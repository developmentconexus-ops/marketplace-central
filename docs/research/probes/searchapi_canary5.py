#!/usr/bin/env python3
"""SearchAPI.io Google Shopping canary. Dry-run by default; never prints credentials."""
from __future__ import annotations

import argparse
import json
import os
import statistics
import time
import unicodedata
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any

API = "https://www.searchapi.io/api/v1/search"
INPUT = Path(__file__).resolve().parents[1] / "evidence" / "2026-07-12-pricefy-canary5-input.json"
OUTPUT = Path(__file__).resolve().parents[1] / "evidence" / "2026-07-13-searchapi-canary5-result.json"


def norm(value: Any) -> str:
    text = unicodedata.normalize("NFKD", str(value or ""))
    return "".join(c.lower() for c in text if c.isalnum())


def safe_url(value: Any) -> str | None:
    if not isinstance(value, str):
        return None
    parsed = urllib.parse.urlsplit(value)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        return None
    return urllib.parse.urlunsplit((parsed.scheme, parsed.netloc.lower(), parsed.path, "", ""))


def is_ml(url: str | None) -> bool:
    if not url:
        return False
    host = urllib.parse.urlsplit(url).hostname or ""
    return host == "mercadolivre.com.br" or host.endswith(".mercadolivre.com.br")


def identity(product: dict[str, Any], *texts: Any) -> dict[str, bool]:
    haystack = norm(" ".join(str(x or "") for x in texts))
    return {
        "ean": norm(product["ean"]) in haystack,
        "brand": norm(product["brand"]) in haystack,
        "reference": norm(product["manufacturer_reference"]) in haystack,
    }


def strong_identity(evidence: dict[str, bool]) -> bool:
    return evidence["ean"] or (evidence["brand"] and evidence["reference"])


def request_json(params: dict[str, Any], key: str, timeout: int) -> tuple[int, dict[str, Any], float]:
    url = API + "?" + urllib.parse.urlencode(params)
    request = urllib.request.Request(
        url,
        headers={"Authorization": f"Bearer {key}", "Accept": "application/json"},
        method="GET",
    )
    started = time.monotonic()
    try:
        with urllib.request.urlopen(request, timeout=timeout) as response:
            payload = json.loads(response.read().decode("utf-8"))
            return response.status, payload, time.monotonic() - started
    except urllib.error.HTTPError as exc:
        return exc.code, {}, time.monotonic() - started
    except (urllib.error.URLError, TimeoutError, UnicodeDecodeError, json.JSONDecodeError):
        return 0, {}, time.monotonic() - started


def query_tasks(products: list[dict[str, Any]]) -> list[dict[str, Any]]:
    tasks = []
    for index, product in enumerate(products):
        tasks.append({"product_index": index, "route": "ean", "q": product["ean"]})
        tasks.append({
            "product_index": index,
            "route": "semantic",
            "q": f'{product["brand"]} {product["manufacturer_reference"]} {product["title"]}',
        })
    return tasks


def shopping_candidates(product: dict[str, Any], payload: dict[str, Any], route: str) -> list[dict[str, Any]]:
    candidates = []
    seen = set()
    for row in payload.get("shopping_results") or []:
        if not isinstance(row, dict) or not row.get("product_token"):
            continue
        token = row["product_token"]
        if token in seen:
            continue
        seen.add(token)
        evidence = identity(product, row.get("title"), row.get("seller"), row.get("snippet"))
        # Google may consume the EAN query without echoing it in the normalized
        # product title. Preserve those tokens only as REVIEW candidates, then
        # require identity again on the concrete merchant offer. Semantic-query
        # candidates must already carry brand + manufacturer reference.
        if not strong_identity(evidence) and route != "ean":
            continue
        candidates.append({
            "token": token,
            "product_id": row.get("product_id"),
            "title": row.get("title"),
            "seller": row.get("seller"),
            "shopping_price": row.get("extracted_price"),
            "identity": evidence,
            "query_route": route,
            "candidate_verdict": "ACCEPT" if strong_identity(evidence) else "REVIEW",
        })
    return candidates[:2]


def offers(product: dict[str, Any], candidate: dict[str, Any], payload: dict[str, Any]) -> list[dict[str, Any]]:
    result = []
    for row in payload.get("offers") or []:
        if not isinstance(row, dict):
            continue
        url = safe_url(row.get("link"))
        price = row.get("extracted_price")
        if not is_ml(url) or isinstance(price, bool) or not isinstance(price, (int, float)) or price <= 0:
            continue
        merchant = row.get("merchant") if isinstance(row.get("merchant"), dict) else {}
        evidence = identity(product, candidate.get("title"), row.get("title"), merchant.get("name"))
        offer_verdict = "ACCEPT" if strong_identity(evidence) else "REVIEW"
        result.append({
            "url": url,
            "title": row.get("title"),
            "merchant": merchant.get("name"),
            "price_brl": float(price),
            "delivery_price_brl": row.get("extracted_delivery_price"),
            "total_price_brl": row.get("extracted_total_price"),
            "condition": row.get("condition"),
            "availability": row.get("availability"),
            "identity": evidence,
            "identity_verdict": offer_verdict,
        })
    return result


def summarize_prices(rows: list[dict[str, Any]]) -> dict[str, Any] | None:
    values = sorted({row["price_brl"] for row in rows})
    if not values:
        return None
    return {
        "n_unique_prices": len(values),
        "min_brl": min(values),
        "median_brl": statistics.median(values),
        "mean_brl": statistics.fmean(values),
        "max_brl": max(values),
    }


def execute(products: list[dict[str, Any]], key: str, timeout: int) -> dict[str, Any]:
    task_results = []
    candidates_by_product: dict[int, list[dict[str, Any]]] = {i: [] for i in range(len(products))}
    tokens_seen: dict[int, set[str]] = {i: set() for i in range(len(products))}
    for task in query_tasks(products):
        status, payload, latency = request_json(
            {"engine": "google_shopping", "q": task["q"], "gl": "br", "hl": "pt-br"},
            key,
            timeout,
        )
        index = task["product_index"]
        found = shopping_candidates(products[index], payload, task["route"])
        for candidate in found:
            if candidate["token"] not in tokens_seen[index]:
                tokens_seen[index].add(candidate["token"])
                candidates_by_product[index].append(candidate)
        task_results.append({"product_index": index, "route": task["route"], "http_status": status, "latency_s": round(latency, 3), "accepted_candidates": len(found)})

    offer_requests = []
    product_offers: dict[int, list[dict[str, Any]]] = {i: [] for i in range(len(products))}
    for index, product in enumerate(products):
        ranked_candidates = sorted(
            candidates_by_product[index],
            key=lambda row: (row["candidate_verdict"] != "ACCEPT", row["query_route"] != "semantic"),
        )
        for candidate in ranked_candidates[:2]:
            status, payload, latency = request_json(
                {"engine": "google_product_offers", "product_token": candidate["token"], "gl": "br", "hl": "pt-br"},
                key,
                timeout,
            )
            found = offers(product, candidate, payload)
            existing = {(row["url"], row["price_brl"]) for row in product_offers[index]}
            product_offers[index].extend(
                row for row in found if (row["url"], row["price_brl"]) not in existing
            )
            offer_requests.append({
                "product_index": index,
                "query_route": candidate["query_route"],
                "candidate_verdict": candidate["candidate_verdict"],
                "http_status": status,
                "latency_s": round(latency, 3),
                "ml_offers": len(found),
            })

    product_results = []
    for index, product in enumerate(products):
        rows = product_offers[index]
        accepted_rows = [row for row in rows if row["identity_verdict"] == "ACCEPT"]
        verdict = "ACCEPT" if accepted_rows else "REVIEW" if candidates_by_product[index] or rows else "REJECT"
        product_results.append({
            "synthetic_id": f"SA-CANARY-{index + 1}",
            "ean": product["ean"],
            "stratum": product["stratum"],
            "candidate_count": len(candidates_by_product[index]),
            "mercado_livre_offer_count": len(rows),
            "accepted_mercado_livre_offer_count": len(accepted_rows),
            "verdict": verdict,
            "price_stats": summarize_prices(accepted_rows),
            "offers": rows,
        })
    return {
        "provider": "SearchAPI.io",
        "method": "Google Shopping then Google Product Offers; gl=br; hl=pt-br",
        "shopping_requests": task_results,
        "offer_requests": offer_requests,
        "products": product_results,
        "summary": {
            "accept": sum(x["verdict"] == "ACCEPT" for x in product_results),
            "review": sum(x["verdict"] == "REVIEW" for x in product_results),
            "reject": sum(x["verdict"] == "REJECT" for x in product_results),
            "products_with_ml_price": sum(bool(x["price_stats"]) for x in product_results),
            "mercado_livre_offers_total": sum(x["mercado_livre_offer_count"] for x in product_results),
            "unverified_mercado_livre_offers": sum(
                row["identity_verdict"] == "REVIEW"
                for product in product_results
                for row in product["offers"]
            ),
        },
    }


def self_test() -> None:
    product = {"ean": "7892162167027", "brand": "FRANKE", "manufacturer_reference": "16702"}
    good = {"shopping_results": [{"product_token": "token", "product_id": "p", "title": "Coifa Franke Linea 16702"}]}
    bad = {"shopping_results": [{"product_token": "bad", "title": "Furadeira Menegotti"}]}
    candidate = shopping_candidates(product, good, "semantic")[0]
    assert shopping_candidates(product, bad, "semantic") == []
    assert shopping_candidates(product, bad, "ean")[0]["candidate_verdict"] == "REVIEW"
    result = offers(product, candidate, {"offers": [{
        "title": "Coifa Franke 16702", "link": "https://produto.mercadolivre.com.br/MLB-1?tracking=secret",
        "extracted_price": 2289.0, "merchant": {"name": "Mercado Livre"},
    }]})
    assert len(result) == 1 and result[0]["price_brl"] == 2289.0
    assert "?" not in result[0]["url"] and strong_identity(result[0]["identity"])
    assert result[0]["identity_verdict"] == "ACCEPT"
    assert summarize_prices(result)["median_brl"] == 2289.0
    weak_candidate = shopping_candidates(product, bad, "ean")[0]
    weak_result = offers(product, weak_candidate, {"offers": [{
        "title": "Furadeira Menegotti", "link": "https://produto.mercadolivre.com.br/MLB-2",
        "extracted_price": 299.0, "merchant": {"name": "Loja de Ferramentas"},
    }]})
    assert weak_result[0]["identity_verdict"] == "REVIEW"


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--execute", action="store_true")
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--timeout", type=int, default=30)
    args = parser.parse_args()
    self_test()
    if args.self_test:
        print("PASS searchapi_canary5 self-test")
        return
    products = json.loads(INPUT.read_text(encoding="utf-8"))["products"]
    if not args.execute:
        print(json.dumps({
            "mode": "dry-run",
            "external_requests": 0,
            "products": len(products),
            "shopping_requests": len(query_tasks(products)),
            "maximum_offer_requests": len(products) * 2,
            "required_env": ["SEARCHAPI_API_KEY", "MPC_SEARCHAPI_TRANSFER_AUTHORIZED=YES"],
        }))
        return
    if os.environ.get("MPC_SEARCHAPI_TRANSFER_AUTHORIZED") != "YES":
        raise SystemExit("External transfer authorization gate is not set")
    key = os.environ.get("SEARCHAPI_API_KEY", "")
    if not key:
        raise SystemExit("SEARCHAPI_API_KEY is not set")
    document = execute(products, key, args.timeout)
    OUTPUT.write_text(json.dumps(document, ensure_ascii=False, indent=2), encoding="utf-8")
    print(json.dumps(document["summary"]))


if __name__ == "__main__":
    main()
