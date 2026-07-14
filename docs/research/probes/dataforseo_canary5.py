#!/usr/bin/env python3
"""Canário DataForSEO -> Google Shopping -> Mercado Livre.

Default is dry-run. Real execution requires --execute and credentials in
DATAFORSEO_LOGIN / DATAFORSEO_PASSWORD. Secrets are never printed or persisted.
"""

from __future__ import annotations

import argparse
import base64
import json
import os
import re
import statistics
import sys
import time
import unicodedata
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any


API = "https://api.dataforseo.com/v3"
INPUT = Path(__file__).resolve().parents[1] / "evidence" / "2026-07-12-pricefy-canary5-input.json"


def norm(value: Any) -> str:
    text = unicodedata.normalize("NFKD", str(value or ""))
    return re.sub(r"[^a-z0-9]", "", text.encode("ascii", "ignore").decode().lower())


def safe_url(value: Any) -> str | None:
    if not value:
        return None
    parsed = urllib.parse.urlsplit(str(value))
    if parsed.scheme not in {"http", "https"}:
        return None
    return urllib.parse.urlunsplit((parsed.scheme, parsed.netloc.lower(), parsed.path, "", ""))


def is_ml(domain: Any = None, url: Any = None) -> bool:
    host = str(domain or "").lower().split(":", 1)[0].lstrip("www.")
    if not host and url:
        host = (urllib.parse.urlsplit(str(url)).hostname or "").lower().lstrip("www.")
    return host == "mercadolivre.com.br" or host.endswith(".mercadolivre.com.br")


def walk_dicts(value: Any):
    if isinstance(value, dict):
        yield value
        for child in value.values():
            yield from walk_dicts(child)
    elif isinstance(value, list):
        for child in value:
            yield from walk_dicts(child)


def product_candidates(task_payload: dict[str, Any], response: dict[str, Any]) -> list[dict[str, Any]]:
    out: list[dict[str, Any]] = []
    seen: set[tuple[Any, ...]] = set()
    for row in walk_dicts(response.get("result", [])):
        if not any(row.get(k) for k in ("product_id", "gid", "data_docid")):
            continue
        if not row.get("title"):
            continue
        key = (row.get("product_id"), row.get("gid"), row.get("data_docid"), row.get("title"))
        if key in seen:
            continue
        seen.add(key)
        out.append({
            "tag": task_payload["tag"],
            "query_route": task_payload["tag"].rsplit(":", 1)[-1],
            "product_id": row.get("product_id"),
            "gid": row.get("gid"),
            "data_docid": row.get("data_docid"),
            "title": row.get("title"),
            "description": row.get("description"),
            "seller": row.get("seller"),
            "domain": row.get("domain"),
            "price": row.get("price"),
            "currency": row.get("currency"),
            "url": safe_url(row.get("url")),
            "additional_specifications": row.get("additional_specifications")
            if isinstance(row.get("additional_specifications"), dict) else {},
            "pvf": row.get("pvf") or (
                row.get("additional_specifications", {}).get("pvf")
                if isinstance(row.get("additional_specifications"), dict) else None
            ),
        })
    return out


def identity_score(product: dict[str, Any], candidate: dict[str, Any]) -> tuple[int, dict[str, bool]]:
    text = norm(" ".join(str(candidate.get(k) or "") for k in ("title", "description")))
    evidence = {
        "ean": norm(product["ean"]) in text,
        "brand": norm(product["brand"]) in text,
        "reference": norm(product["manufacturer_reference"]) in text,
    }
    return sum(evidence.values()), evidence


def seller_rows(response: dict[str, Any]) -> list[dict[str, Any]]:
    rows: list[dict[str, Any]] = []
    for row in walk_dicts(response.get("result", [])):
        if not row.get("title") or row.get("base_price") is None:
            continue
        if not (row.get("domain") or row.get("url")):
            continue
        rows.append({
            "domain": str(row.get("domain") or "").lower(),
            "url": safe_url(row.get("url")),
            "title": row.get("title"),
            "seller_name": row.get("seller_name"),
            "base_price": row.get("base_price"),
            "total_price": row.get("total_price"),
            "currency": row.get("currency"),
            "condition": row.get("product_condition"),
            "availability": row.get("product_availability"),
        })
    return rows


def request_json(method: str, path: str, login: str, password: str, payload: Any = None) -> dict[str, Any]:
    token = base64.b64encode(f"{login}:{password}".encode()).decode()
    body = None if payload is None else json.dumps(payload, ensure_ascii=False).encode("utf-8")
    req = urllib.request.Request(
        API + path,
        data=body,
        method=method,
        headers={"Authorization": f"Basic {token}", "Content-Type": "application/json"},
    )
    try:
        with urllib.request.urlopen(req, timeout=30) as res:
            return json.load(res)
    except urllib.error.HTTPError as exc:
        message = exc.read(1000).decode("utf-8", "replace")
        raise RuntimeError(f"DataForSEO HTTP {exc.code}: {message}") from exc


def post_tasks(path: str, payloads: list[dict[str, Any]], login: str, password: str) -> list[str]:
    response = request_json("POST", path, login, password, payloads)
    if response.get("status_code") != 20000 or response.get("tasks_error"):
        raise RuntimeError(f"Task POST failed: {response.get('status_code')} {response.get('status_message')}")
    return [task["id"] for task in response.get("tasks", [])]


def poll_tasks(path: str, ids: list[str], login: str, password: str, timeout_s: int) -> list[dict[str, Any]]:
    pending = set(ids)
    done: dict[str, dict[str, Any]] = {}
    deadline = time.monotonic() + timeout_s
    while pending and time.monotonic() < deadline:
        for task_id in list(pending):
            response = request_json("GET", f"{path}/{task_id}", login, password)
            task = (response.get("tasks") or [{}])[0]
            if task.get("status_code") == 20000 and task.get("result") is not None:
                done[task_id] = task
                pending.remove(task_id)
            elif task.get("status_code", 0) >= 40000 and task.get("status_code") not in {40601, 40602}:
                raise RuntimeError(f"Task {task_id} failed: {task.get('status_code')} {task.get('status_message')}")
        if pending:
            time.sleep(5)
    if pending:
        raise TimeoutError(f"Timed out waiting for {len(pending)} task(s)")
    return [done[task_id] for task_id in ids]


def build_product_tasks(products: list[dict[str, Any]]) -> list[dict[str, Any]]:
    tasks: list[dict[str, Any]] = []
    for index, p in enumerate(products, 1):
        base = {
            "location_name": "Brazil",
            "language_code": "pt",
            "se_domain": "google.com.br",
            "depth": 40,
            "priority": 2,
        }
        tasks.append({**base, "keyword": p["ean"], "tag": f"DFS-CANARY-{index}:ean"})
        tasks.append({
            **base,
            "keyword": f"{p['brand']} {p['manufacturer_reference']} {p['title']}",
            "tag": f"DFS-CANARY-{index}:semantic",
        })
    return tasks


def aggregate(values: list[float]) -> dict[str, float] | None:
    if not values:
        return None
    return {
        "min": min(values),
        "median": statistics.median(values),
        "mean": round(statistics.fmean(values), 2),
    }


def execute(products: list[dict[str, Any]], timeout_s: int) -> dict[str, Any]:
    login, password = os.getenv("DATAFORSEO_LOGIN"), os.getenv("DATAFORSEO_PASSWORD")
    if not login or not password:
        raise RuntimeError("Set DATAFORSEO_LOGIN and DATAFORSEO_PASSWORD before --execute")

    product_tasks = build_product_tasks(products)
    product_ids = post_tasks("/merchant/google/products/task_post", product_tasks, login, password)
    product_responses = poll_tasks(
        "/merchant/google/products/task_get/advanced", product_ids, login, password, timeout_s
    )

    by_product: dict[int, list[dict[str, Any]]] = {i: [] for i in range(len(products))}
    for payload, response in zip(product_tasks, product_responses):
        index = int(payload["tag"].split(":", 1)[0].rsplit("-", 1)[1]) - 1
        by_product[index].extend(product_candidates(payload, response))

    seller_payloads: list[dict[str, Any]] = []
    seller_context: list[tuple[int, dict[str, Any], dict[str, bool], str]] = []
    for index, product in enumerate(products):
        ranked = []
        for candidate in by_product[index]:
            score, evidence = identity_score(product, candidate)
            strong = evidence["brand"] and evidence["reference"]
            # An EAN query may resolve to a Google product cluster without the
            # identifier being echoed. Keep it for REVIEW, but semantic-query
            # candidates must already prove brand + manufacturer reference.
            if strong or candidate["query_route"] == "ean":
                verdict = "ACCEPT" if strong else "REVIEW"
                ranked.append((verdict, score, candidate, evidence))
        ranked.sort(
            key=lambda x: (x[0] == "ACCEPT", x[1], bool(x[2].get("domain"))),
            reverse=True,
        )
        for candidate_verdict, _, candidate, evidence in ranked[:2]:
            payload = {
                "location_name": "Brazil",
                "language_code": "pt",
                "se_domain": "google.com.br",
                "priority": 2,
                "tag": f"DFS-CANARY-{index + 1}:seller",
            }
            for key in ("product_id", "data_docid", "gid"):
                if candidate.get(key):
                    payload[key] = candidate[key]
            seller_payloads.append(payload)
            seller_context.append((index, candidate, evidence, candidate_verdict))

    seller_responses: list[dict[str, Any]] = []
    if seller_payloads:
        ids = post_tasks("/merchant/google/sellers/task_post", seller_payloads, login, password)
        seller_responses = poll_tasks(
            "/merchant/google/sellers/task_get/advanced", ids, login, password, timeout_s
        )

    results = []
    per_product_offers: dict[int, list[dict[str, Any]]] = {i: [] for i in range(len(products))}
    for (index, candidate, evidence, candidate_verdict), response in zip(seller_context, seller_responses):
        existing = {(row["url"], row["base_price"]) for row in per_product_offers[index]}
        for row in seller_rows(response):
            if not is_ml(row["domain"], row["url"]):
                continue
            key = (row["url"], row["base_price"])
            if key in existing:
                continue
            existing.add(key)
            row["identity"] = evidence
            row["identity_verdict"] = candidate_verdict
            row["query_route"] = candidate["query_route"]
            per_product_offers[index].append(row)

    for index, product in enumerate(products):
        offers = per_product_offers[index]
        accepted_offers = [row for row in offers if row["identity_verdict"] == "ACCEPT"]
        verdict = "ACCEPT" if accepted_offers else "REVIEW" if offers or by_product[index] else "REJECT"
        prices = [
            float(o["base_price"])
            for o in accepted_offers
            if o.get("currency") == "BRL"
            and isinstance(o.get("base_price"), (int, float))
            and not isinstance(o.get("base_price"), bool)
            and o["base_price"] > 0
        ]
        results.append({
            "synthetic_id": f"DFS-CANARY-{index + 1}",
            "ean": product["ean"],
            "verdict": verdict,
            "product_candidates": len(by_product[index]),
            "mercado_livre_offers": offers,
            "accepted_mercado_livre_offers": len(accepted_offers),
            "price_brl": aggregate(prices),
        })

    return {
        "provider": "DataForSEO",
        "captured_at_epoch": int(time.time()),
        "products": results,
        "summary": {
            "accept": sum(r["verdict"] == "ACCEPT" for r in results),
            "review": sum(r["verdict"] == "REVIEW" for r in results),
            "reject": sum(r["verdict"] == "REJECT" for r in results),
            "false_accept_gate": "manual audit required; known collision products must remain non-ACCEPT",
        },
    }


def self_test() -> None:
    product = {"ean": "7892162167027", "brand": "FRANKE", "manufacturer_reference": "16702"}
    good = {"title": "Coifa Franke Linea Touch 16702", "description": "EAN 7892162167027"}
    bad = {"title": "Furadeira Menegotti", "description": "7890013123284"}
    score, ev = identity_score(product, good)
    assert score == 3 and all(ev.values())
    assert identity_score(product, bad)[0] == 0
    assert is_ml("www.mercadolivre.com.br", None)
    assert not is_ml("example.com", "https://example.com/x")
    assert safe_url("https://produto.mercadolivre.com.br/MLB-1?x=secret#frag") == "https://produto.mercadolivre.com.br/MLB-1"
    assert aggregate([10.0, 20.0, 30.0]) == {"min": 10.0, "median": 20.0, "mean": 20.0}
    assert build_product_tasks([{
        "ean": "7892162167027", "brand": "FRANKE", "manufacturer_reference": "16702",
        "title": "Coifa Linea",
    }])[0]["tag"].endswith(":ean")
    print("PASS dataforseo_canary5 self-test")


def main() -> int:
    parser = argparse.ArgumentParser()
    mode = parser.add_mutually_exclusive_group()
    mode.add_argument("--execute", action="store_true")
    mode.add_argument("--self-test", action="store_true")
    parser.add_argument("--timeout", type=int, default=180)
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()

    if args.self_test:
        self_test()
        return 0
    products = json.loads(INPUT.read_text(encoding="utf-8"))["products"]
    if not args.execute:
        print(json.dumps({
            "mode": "dry-run",
            "external_requests": 0,
            "products": len(products),
            "product_tasks": len(build_product_tasks(products)),
            "routes": ["ean", "brand + manufacturer reference + title"],
            "maximum_seller_tasks": len(products) * 2,
            "required_env": [
                "DATAFORSEO_LOGIN",
                "DATAFORSEO_PASSWORD",
                "MPC_DATAFORSEO_TRANSFER_AUTHORIZED=YES"
            ],
        }, ensure_ascii=False, indent=2))
        return 0

    if os.getenv("MPC_DATAFORSEO_TRANSFER_AUTHORIZED") != "YES":
        raise RuntimeError("External transfer authorization gate is not set")
    result = execute(products, args.timeout)
    rendered = json.dumps(result, ensure_ascii=False, indent=2)
    if args.output:
        args.output.write_text(rendered + "\n", encoding="utf-8")
    print(rendered)
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        print(f"ERROR: {exc}", file=sys.stderr)
        raise SystemExit(1)
