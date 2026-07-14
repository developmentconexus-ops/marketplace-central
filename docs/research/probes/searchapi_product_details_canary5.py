#!/usr/bin/env python3
"""SearchAPI product-detail canary. Dry-run by default; never persists tokens or credentials."""
from __future__ import annotations

import argparse
import json
import os
import re
import unicodedata
from pathlib import Path
from typing import Any

import searchapi_canary5 as base


INPUT = Path(__file__).resolve().parents[1] / "evidence" / "2026-07-12-pricefy-canary5-input.json"
OUTPUT = Path(__file__).resolve().parents[1] / "evidence" / "2026-07-13-searchapi-product-details-canary5-result.json"
STOPWORDS = {"a", "e", "o", "as", "os", "de", "da", "do", "das", "dos", "com", "para", "por", "the", "and"}
TEXT_KEYS = {"title", "name", "value", "description", "label", "type"}


def fold(value: Any) -> str:
    text = unicodedata.normalize("NFKD", str(value or ""))
    return "".join(ch.lower() for ch in text if not unicodedata.combining(ch))


def title_tokens(value: Any) -> set[str]:
    return {token for token in re.findall(r"[a-z0-9]+", fold(value)) if len(token) >= 2 and token not in STOPWORDS}


def token_recall(target: Any, observed: Any) -> float:
    target_set = title_tokens(target)
    observed_set = title_tokens(observed)
    if not target_set:
        return 0.0
    matched = 0
    for token in target_set:
        if any(
            token == candidate
            or (len(token) >= 3 and candidate.startswith(token))
            or (len(candidate) >= 3 and token.startswith(candidate))
            for candidate in observed_set
        ):
            matched += 1
    return matched / len(target_set)


def selected_text(value: Any) -> list[str]:
    out: list[str] = []
    if isinstance(value, dict):
        for key, child in value.items():
            if key in TEXT_KEYS and isinstance(child, (str, int, float)) and not isinstance(child, bool):
                out.append(str(child))
            elif isinstance(child, (dict, list)):
                out.extend(selected_text(child))
    elif isinstance(value, list):
        for child in value:
            out.extend(selected_text(child))
    return out


def contradictions(product: dict[str, Any], candidate_text: str) -> list[str]:
    target = fold(product.get("title"))
    candidate = fold(candidate_text)
    found: list[str] = []
    target_voltages = set(re.findall(r"\b(?:110|127|220)\s*v\b", target))
    candidate_voltages = set(re.findall(r"\b(?:110|127|220)\s*v\b", candidate))
    if target_voltages and candidate_voltages and target_voltages.isdisjoint(candidate_voltages):
        found.append("voltage")
    if ("kit" in title_tokens(candidate)) != ("kit" in title_tokens(target)):
        found.append("bundle")
    return found


def detail_identity(product: dict[str, Any], shopping_row: dict[str, Any], payload: dict[str, Any]) -> dict[str, Any]:
    detail = payload.get("product") if isinstance(payload.get("product"), dict) else {}
    core_detail_parts = [
        detail.get("title"), detail.get("brand"), detail.get("description"),
        *selected_text(payload.get("specifications") or []),
        *selected_text(detail.get("variations") or []),
    ]
    detail_text_parts = [shopping_row.get("title"), shopping_row.get("seller"), *core_detail_parts]
    combined = " ".join(str(value or "") for value in detail_text_parts)
    exact = base.identity(product, combined)
    recall = token_recall(product.get("title"), combined)
    authoritative_detail = " ".join(str(value or "") for value in core_detail_parts).strip()
    conflict = contradictions(product, authoritative_detail or combined)
    strong = exact["brand"] and (exact["reference"] or recall >= 0.60) and not conflict
    return {
        **exact,
        "title_token_recall": round(recall, 3),
        "contradictions": conflict,
        "verdict": "ACCEPT" if strong else "REVIEW",
    }


def offer_identity(product: dict[str, Any], offer_title: Any, detail_brand: Any) -> dict[str, Any]:
    observed = f"{detail_brand or ''} {offer_title or ''}"
    exact = base.identity(product, observed)
    recall = token_recall(product.get("title"), offer_title)
    conflict = contradictions(product, str(offer_title or ""))
    strong = exact["brand"] and (exact["reference"] or recall >= 0.75) and not conflict
    return {
        **exact,
        "title_token_recall": round(recall, 3),
        "contradictions": conflict,
        "verdict": "ACCEPT" if strong else "REVIEW",
    }


def safe_detail(payload: dict[str, Any]) -> dict[str, Any]:
    product = payload.get("product") if isinstance(payload.get("product"), dict) else {}
    specs = selected_text(payload.get("specifications") or [])[:30]
    return {
        "title": product.get("title"),
        "brand": product.get("brand"),
        "description": str(product.get("description") or "")[:600] or None,
        "specification_text": specs,
    }


def execute(products: list[dict[str, Any]], key: str, timeout: int) -> dict[str, Any]:
    product_results: list[dict[str, Any]] = []
    request_counts = {"shopping": 0, "details": 0, "offers": 0}
    for index, product in enumerate(products):
        query = f'{product["brand"]} {product["manufacturer_reference"]} {product["title"]}'
        status, shopping_payload, latency = base.request_json(
            {"engine": "google_shopping", "q": query, "gl": "br", "hl": "pt-br"}, key, timeout
        )
        request_counts["shopping"] += 1
        candidate_results: list[dict[str, Any]] = []
        accepted_offers: list[dict[str, Any]] = []
        rows = [
            row for row in (shopping_payload.get("shopping_results") or [])
            if isinstance(row, dict) and row.get("product_token")
        ][:2]
        for row in rows:
            detail_status, detail_payload, detail_latency = base.request_json(
                {"engine": "google_product", "product_token": row["product_token"], "gl": "br", "hl": "pt-br"},
                key,
                timeout,
            )
            request_counts["details"] += 1
            evidence = detail_identity(product, row, detail_payload)
            candidate_record = {
                "shopping_title": row.get("title"),
                "shopping_seller": row.get("seller"),
                "shopping_price": row.get("extracted_price"),
                "detail_http_status": detail_status,
                "detail_latency_s": round(detail_latency, 3),
                "detail": safe_detail(detail_payload),
                "identity": evidence,
            }
            candidate_results.append(candidate_record)
            if evidence["verdict"] != "ACCEPT":
                continue
            offers_status, offers_payload, offers_latency = base.request_json(
                {"engine": "google_product_offers", "product_token": row["product_token"], "gl": "br", "hl": "pt-br"},
                key,
                timeout,
            )
            request_counts["offers"] += 1
            found = base.offers(product, {"title": row.get("title")}, offers_payload)
            for offer in found:
                direct_evidence = offer_identity(product, offer.get("title"), candidate_record["detail"].get("brand"))
                offer["identity"] = direct_evidence
                offer["identity_verdict"] = direct_evidence["verdict"]
                offer["offer_contradictions"] = direct_evidence["contradictions"]
                if offer["identity_verdict"] == "ACCEPT":
                    accepted_offers.append(offer)
            candidate_record["offers_http_status"] = offers_status
            candidate_record["offers_latency_s"] = round(offers_latency, 3)
            candidate_record["mercado_livre_offers"] = len(found)
        deduped = list({(row["url"], row["price_brl"]): row for row in accepted_offers}.values())
        product_results.append({
            "synthetic_id": f"SA-DETAIL-CANARY-{index + 1}",
            "ean": product["ean"],
            "stratum": product["stratum"],
            "shopping_http_status": status,
            "shopping_latency_s": round(latency, 3),
            "candidate_count": len(candidate_results),
            "candidates": candidate_results,
            "accepted_mercado_livre_offers": deduped,
            "price_stats": base.summarize_prices(deduped),
            "verdict": "ACCEPT" if deduped else "REVIEW" if candidate_results else "REJECT",
        })
    return {
        "provider": "SearchAPI.io",
        "method": "semantic Google Shopping -> Google Product details -> identity gate -> Google Product Offers",
        "request_counts": request_counts,
        "products": product_results,
        "summary": {
            "accept": sum(row["verdict"] == "ACCEPT" for row in product_results),
            "review": sum(row["verdict"] == "REVIEW" for row in product_results),
            "reject": sum(row["verdict"] == "REJECT" for row in product_results),
            "products_with_accepted_ml_price": sum(bool(row["price_stats"]) for row in product_results),
        },
    }


def self_test() -> None:
    product = {
        "ean": "7892162167027", "brand": "FRANKE", "manufacturer_reference": "16702",
        "title": "COIFA ILHA LINEA TOUCH 90CM 220V",
    }
    good = {"title": "Coifa Ilha Franke Linea Touch 90 cm 220V", "seller": "Loja"}
    good_detail = {"product": {"title": good["title"], "brand": "Franke", "description": "Coifa de ilha"}}
    assert detail_identity(product, good, good_detail)["verdict"] == "ACCEPT"
    wrong_brand = {"product": {"title": "Furadeira Menegotti 220V", "brand": "Menegotti"}}
    assert detail_identity(product, {"title": "Furadeira"}, wrong_brand)["verdict"] == "REVIEW"
    wrong_voltage = {"product": {"title": "Coifa Ilha Franke Linea Touch 90 cm 127V", "brand": "Franke"}}
    evidence = detail_identity(product, good, wrong_voltage)
    assert evidence["verdict"] == "REVIEW" and "voltage" in evidence["contradictions"]
    assert offer_identity(product, "Coifa De Ilha Franke Linea Touch 90cm 220V", "Franke")["verdict"] == "ACCEPT"
    wrong_model = offer_identity(product, "Coifa De Ilha Franke Crystal FCR925 90cm 220V", "Franke")
    assert wrong_model["verdict"] == "REVIEW"
    docol = {
        "ean": "7891461390075", "brand": "DOCOL", "manufacturer_reference": "90010032072",
        "title": "ACAB.MON.CHUV E DUCHA NEW EDGE OURO ESC",
    }
    expanded = "Acabamento Monocomando para Chuveiro e Ducha Docol New Edge Ouro Escovado"
    assert offer_identity(docol, expanded, "Docol")["verdict"] == "ACCEPT"


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--execute", action="store_true")
    parser.add_argument("--self-test", action="store_true")
    parser.add_argument("--timeout", type=int, default=30)
    args = parser.parse_args()
    self_test()
    if args.self_test:
        print("PASS searchapi_product_details_canary5 self-test")
        return
    products = json.loads(INPUT.read_text(encoding="utf-8"))["products"]
    if not args.execute:
        print(json.dumps({
            "mode": "dry-run", "external_requests": 0, "products": len(products),
            "maximum_requests": len(products) * 5,
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
