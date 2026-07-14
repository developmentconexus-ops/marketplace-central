#!/usr/bin/env python3
"""DataForSEO Products -> Product Info -> Sellers canary. Dry-run by default."""
from __future__ import annotations

import argparse
import json
import os
import time
from pathlib import Path
from typing import Any

import dataforseo_canary5 as dfs
from searchapi_product_details_canary5 import contradictions, offer_identity, selected_text, token_recall


INPUT = Path(__file__).resolve().parents[1] / "evidence" / "2026-07-12-pricefy-canary5-input.json"
OUTPUT = Path(__file__).resolve().parents[1] / "evidence" / "2026-07-13-dataforseo-product-details-canary5-result.json"


def identifiers(candidate: dict[str, Any]) -> dict[str, Any]:
    result = {
        key: candidate[key]
        for key in ("product_id", "data_docid", "gid")
        if candidate.get(key)
    }
    pvf = candidate.get("pvf")
    if pvf:
        result["pvf"] = pvf
    return result


def extract_product_info(task: dict[str, Any]) -> dict[str, Any]:
    for row in dfs.walk_dicts(task.get("result", [])):
        if isinstance(row.get("product_info"), dict):
            return row["product_info"]
    for row in dfs.walk_dicts(task.get("result", [])):
        if row.get("title") and any(key in row for key in ("specifications", "variations", "description")):
            return row
    return {}


def info_identity(product: dict[str, Any], candidate: dict[str, Any], info: dict[str, Any]) -> dict[str, Any]:
    parts = [
        candidate.get("title"), candidate.get("description"),
        info.get("title"), info.get("brand"), info.get("description"),
        *selected_text(info.get("specifications") or []),
        *selected_text(info.get("variations") or []),
    ]
    text = " ".join(str(value or "") for value in parts)
    _, exact = dfs.identity_score(product, {"title": text, "description": ""})
    recall = token_recall(product.get("title"), text)
    conflict = contradictions(product, text)
    strong = exact["brand"] and (exact["reference"] or recall >= 0.60) and not conflict
    return {
        **exact,
        "title_token_recall": round(recall, 3),
        "contradictions": conflict,
        "verdict": "ACCEPT" if strong else "REVIEW",
    }


def safe_info(info: dict[str, Any]) -> dict[str, Any]:
    return {
        "title": info.get("title"),
        "brand": info.get("brand"),
        "description": str(info.get("description") or "")[:600] or None,
        "specification_text": selected_text(info.get("specifications") or [])[:30],
        "variation_count": len(info.get("variations") or []) if isinstance(info.get("variations"), list) else 0,
    }


def choose_candidates(product: dict[str, Any], candidates: list[dict[str, Any]]) -> list[dict[str, Any]]:
    unique: dict[tuple[Any, ...], dict[str, Any]] = {}
    for candidate in candidates:
        key = (candidate.get("product_id"), candidate.get("data_docid"), candidate.get("gid"))
        if not any(key) or key in unique:
            continue
        unique[key] = candidate
    ranked = []
    for candidate in unique.values():
        score, evidence = dfs.identity_score(product, candidate)
        semantic_route = str(candidate.get("tag") or "").endswith(":semantic")
        ranked.append((semantic_route, score, evidence["brand"], evidence["reference"], candidate))
    ranked.sort(key=lambda row: (row[0], row[1], row[2], row[3]), reverse=True)
    return [row[-1] for row in ranked[:2]]


def common_task(candidate: dict[str, Any], *, include_pvf: bool) -> dict[str, Any]:
    ids = identifiers(candidate)
    if not include_pvf:
        ids.pop("pvf", None)
    task = {
        "location_name": "Brazil",
        "language_code": "pt",
        "priority": 2,
        **ids,
    }
    return task


def execute(products: list[dict[str, Any]], timeout_s: int) -> dict[str, Any]:
    login, password = os.getenv("DATAFORSEO_LOGIN"), os.getenv("DATAFORSEO_PASSWORD")
    if not login or not password:
        raise RuntimeError("Set DATAFORSEO_LOGIN and DATAFORSEO_PASSWORD before --execute")

    product_tasks = dfs.build_product_tasks(products)
    product_ids = dfs.post_tasks("/merchant/google/products/task_post", product_tasks, login, password)
    product_responses = dfs.poll_tasks(
        "/merchant/google/products/task_get/advanced", product_ids, login, password, timeout_s
    )
    by_product: dict[int, list[dict[str, Any]]] = {index: [] for index in range(len(products))}
    for payload, response in zip(product_tasks, product_responses):
        index = int(payload["tag"].split(":", 1)[0].rsplit("-", 1)[1]) - 1
        by_product[index].extend(dfs.product_candidates(payload, response))

    info_payloads: list[dict[str, Any]] = []
    info_context: list[tuple[int, dict[str, Any]]] = []
    for index, product in enumerate(products):
        for candidate in choose_candidates(product, by_product[index]):
            info_payloads.append(common_task(candidate, include_pvf=False))
            info_context.append((index, candidate))
    info_responses: list[dict[str, Any]] = []
    if info_payloads:
        ids = dfs.post_tasks("/merchant/google/product_info/task_post", info_payloads, login, password)
        info_responses = dfs.poll_tasks(
            "/merchant/google/product_info/task_get/advanced", ids, login, password, timeout_s
        )

    seller_payloads: list[dict[str, Any]] = []
    seller_context: list[tuple[int, dict[str, Any], dict[str, Any], dict[str, Any]]] = []
    candidate_records: dict[int, list[dict[str, Any]]] = {index: [] for index in range(len(products))}
    for (index, candidate), response in zip(info_context, info_responses):
        info = extract_product_info(response)
        evidence = info_identity(products[index], candidate, info)
        record = {
            "shopping_title": candidate.get("title"),
            "shopping_seller": candidate.get("seller"),
            "shopping_price": candidate.get("price"),
            "identifiers_present": sorted(identifiers(candidate)),
            "pvf_preserved": bool(candidate.get("pvf")),
            "product_info": safe_info(info),
            "identity": evidence,
        }
        candidate_records[index].append(record)
        if evidence["verdict"] == "ACCEPT":
            seller_payloads.append(common_task(candidate, include_pvf=True))
            seller_context.append((index, candidate, info, record))

    seller_responses: list[dict[str, Any]] = []
    if seller_payloads:
        ids = dfs.post_tasks("/merchant/google/sellers/task_post", seller_payloads, login, password)
        seller_responses = dfs.poll_tasks(
            "/merchant/google/sellers/task_get/advanced", ids, login, password, timeout_s
        )

    accepted_by_product: dict[int, list[dict[str, Any]]] = {index: [] for index in range(len(products))}
    for (index, _candidate, info, record), response in zip(seller_context, seller_responses):
        rows = [row for row in dfs.seller_rows(response) if dfs.is_ml(row.get("domain"), row.get("url"))]
        record["mercado_livre_offers"] = len(rows)
        for row in rows:
            evidence = offer_identity(products[index], row.get("title"), info.get("brand"))
            row["identity"] = evidence
            row["identity_verdict"] = evidence["verdict"]
            if evidence["verdict"] == "ACCEPT":
                accepted_by_product[index].append(row)

    results = []
    for index, product in enumerate(products):
        accepted = list({(row.get("url"), row.get("base_price")): row for row in accepted_by_product[index]}.values())
        prices = [
            float(row["base_price"])
            for row in accepted
            if row.get("currency") == "BRL"
            and isinstance(row.get("base_price"), (int, float))
            and not isinstance(row.get("base_price"), bool)
            and row["base_price"] > 0
        ]
        candidates = candidate_records[index]
        results.append({
            "synthetic_id": f"DFS-DETAIL-CANARY-{index + 1}",
            "ean": product["ean"],
            "stratum": product["stratum"],
            "candidate_count": len(candidates),
            "candidates": candidates,
            "accepted_mercado_livre_offers": accepted,
            "price_brl": dfs.aggregate(prices),
            "verdict": "ACCEPT" if prices else "REVIEW" if candidates else "REJECT",
        })
    return {
        "provider": "DataForSEO",
        "method": "Products -> Product Info -> identity gate -> Sellers with all identifiers and pvf -> offer gate",
        "captured_at_epoch": int(time.time()),
        "task_counts": {
            "products": len(product_tasks),
            "product_info": len(info_payloads),
            "sellers": len(seller_payloads),
        },
        "products": results,
        "summary": {
            "accept": sum(row["verdict"] == "ACCEPT" for row in results),
            "review": sum(row["verdict"] == "REVIEW" for row in results),
            "reject": sum(row["verdict"] == "REJECT" for row in results),
            "products_with_accepted_ml_price": sum(bool(row["price_brl"]) for row in results),
        },
    }


def self_test() -> None:
    candidate = {
        "product_id": "p", "data_docid": "d", "gid": "g",
        "additional_specifications": {"pvf": "color:red"}, "pvf": "color:red",
    }
    assert identifiers(candidate) == {
        "product_id": "p", "data_docid": "d", "gid": "g", "pvf": "color:red"
    }
    assert "pvf" not in common_task(candidate, include_pvf=False)
    assert common_task(candidate, include_pvf=True)["pvf"] == "color:red"
    task = {"result": [{"product_info": {"title": "Coifa Franke Linea", "brand": "Franke"}}]}
    assert extract_product_info(task)["brand"] == "Franke"
    product = {
        "ean": "7892162167027", "brand": "FRANKE", "manufacturer_reference": "16702",
        "title": "COIFA ILHA LINEA TOUCH 90CM 220V",
    }
    info = {"title": "Coifa de Ilha Franke Linea Touch 90cm 220V", "brand": "Franke"}
    evidence = info_identity(product, {"title": info["title"]}, info)
    assert evidence["verdict"] == "ACCEPT"
    assert offer_identity(product, "Coifa Franke Crystal FCR925 90cm 220V", "Franke")["verdict"] == "REVIEW"


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--execute", action="store_true")
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--timeout", type=int, default=240)
    args = parser.parse_args()
    self_test()
    if args.self_test:
        print("PASS dataforseo_product_details_canary5 self-test")
        return 0
    products = json.loads(INPUT.read_text(encoding="utf-8"))["products"]
    if not args.execute:
        print(json.dumps({
            "mode": "dry-run",
            "external_requests": 0,
            "products": len(products),
            "product_tasks": len(products) * 2,
            "maximum_product_info_tasks": len(products) * 2,
            "maximum_seller_tasks": len(products) * 2,
            "required_env": [
                "DATAFORSEO_LOGIN", "DATAFORSEO_PASSWORD", "MPC_DATAFORSEO_TRANSFER_AUTHORIZED=YES"
            ],
        }))
        return 0
    if os.getenv("MPC_DATAFORSEO_TRANSFER_AUTHORIZED") != "YES":
        raise RuntimeError("External transfer authorization gate is not set")
    document = execute(products, args.timeout)
    OUTPUT.write_text(json.dumps(document, ensure_ascii=False, indent=2), encoding="utf-8")
    print(json.dumps(document["summary"]))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
