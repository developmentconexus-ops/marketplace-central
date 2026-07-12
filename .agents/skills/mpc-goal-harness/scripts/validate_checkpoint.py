#!/usr/bin/env python3
"""Validate a harness checkpoint without mutating repository or state files."""

from __future__ import annotations

import argparse
import json
import re
import sys
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


EVENT_TYPES = {
    "packet_accepted", "feature_accepted", "blocker", "fixed_sha_frozen",
    "review_verdict", "qa_verdict", "terminal_handoff", "heartbeat",
}
REQUIRED_FIELDS = (
    "schema_version", "milestone_id", "source_task_id", "parent_task_id",
    "event_id", "sequence", "event_type", "dispatch_id", "work_item_id",
    "base_sha", "commit_or_frozen_sha", "emitted_at", "callback_due_at",
    "status", "checkpoint_id", "commit", "changed_paths", "evidence",
    "review", "blockers", "next",
)
SHA_RE = re.compile(r"^[0-9a-f]{40}$")


def error(code: str, field: str, detail: str) -> dict[str, str]:
    return {"code": code, "field": field, "detail": detail}


def parse_time(value: Any) -> datetime | None:
    if not isinstance(value, str):
        return None
    try:
        parsed = datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError:
        return None
    if parsed.tzinfo is None:
        return None
    return parsed.astimezone(timezone.utc)


def _nonempty_strings(data: dict[str, Any], fields: tuple[str, ...], errors: list[dict[str, str]]) -> None:
    for field in fields:
        if field in data and (not isinstance(data[field], str) or not data[field].strip()):
            errors.append(error("invalid_field", field, "must be a non-empty string"))


def validate_prior_state(prior_state: Any, errors: list[dict[str, str]]) -> tuple[set[str], list[int]]:
    if not isinstance(prior_state, dict):
        errors.append(error("invalid_prior_state", "$", "must be a JSON object"))
        return set(), []
    expected = {"events", "max_sequence"}
    for field in sorted(set(prior_state) - expected):
        errors.append(error("unexpected_prior_state_field", field, "is not supported"))
    for field in sorted(expected - set(prior_state)):
        errors.append(error("missing_prior_state_field", field, "is required when prior state is supplied"))
    events = prior_state.get("events")
    max_sequence = prior_state.get("max_sequence")
    if not isinstance(events, list):
        errors.append(error("invalid_prior_state", "events", "must be an array"))
        return set(), []
    if isinstance(max_sequence, bool) or not isinstance(max_sequence, int) or max_sequence < 0:
        errors.append(error("invalid_prior_state", "max_sequence", "must be an integer >= 0"))
    seen_ids: set[str] = set()
    sequences: list[int] = []
    for index, entry in enumerate(events):
        if not isinstance(entry, dict):
            errors.append(error("invalid_prior_state", f"events[{index}]", "must be an object"))
            continue
        event_id = entry.get("event_id")
        sequence = entry.get("sequence")
        if not isinstance(event_id, str) or not event_id.strip():
            errors.append(error("invalid_prior_state", f"events[{index}].event_id", "must be a non-empty string"))
        else:
            seen_ids.add(event_id)
        if isinstance(sequence, bool) or not isinstance(sequence, int) or sequence < 1:
            errors.append(error("invalid_prior_state", f"events[{index}].sequence", "must be an integer >= 1"))
        else:
            sequences.append(sequence)
    if isinstance(max_sequence, int) and not isinstance(max_sequence, bool) and sequences and max_sequence < max(sequences):
        errors.append(error("invalid_prior_state", "max_sequence", "must be at least the largest event sequence"))
    return seen_ids, sequences + ([max_sequence] if isinstance(max_sequence, int) and not isinstance(max_sequence, bool) else [])


def validate_checkpoint(checkpoint: Any, prior_state: Any | None = None) -> list[dict[str, str]]:
    """Return stable, sorted validation errors; never update prior_state."""
    errors: list[dict[str, str]] = []
    if not isinstance(checkpoint, dict):
        return [error("invalid_checkpoint", "$", "must be a JSON object")]

    unknown = sorted(set(checkpoint) - set(REQUIRED_FIELDS))
    for field in unknown:
        errors.append(error("unexpected_field", field, "is not in schema version 1.0"))
    for field in REQUIRED_FIELDS:
        if field not in checkpoint:
            errors.append(error("missing_required_field", field, "is required"))
    if errors:
        return errors

    _nonempty_strings(checkpoint, (
        "milestone_id", "source_task_id", "parent_task_id", "event_id",
        "dispatch_id", "work_item_id", "status", "checkpoint_id", "next",
    ), errors)
    if checkpoint["schema_version"] != "1.0":
        errors.append(error("invalid_schema_version", "schema_version", "must equal 1.0"))
    if isinstance(checkpoint["sequence"], bool) or not isinstance(checkpoint["sequence"], int) or checkpoint["sequence"] < 1:
        errors.append(error("invalid_field", "sequence", "must be an integer >= 1"))
    if checkpoint["event_type"] not in EVENT_TYPES:
        errors.append(error("invalid_event_type", "event_type", "is not a supported checkpoint event"))
    for field in ("base_sha", "commit_or_frozen_sha", "commit"):
        value = checkpoint[field]
        if value is not None and (not isinstance(value, str) or not SHA_RE.fullmatch(value)):
            errors.append(error("invalid_sha", field, "must be a lowercase 40-character SHA or null"))
    if not isinstance(checkpoint["base_sha"], str) or not SHA_RE.fullmatch(checkpoint["base_sha"]):
        errors.append(error("invalid_sha", "base_sha", "must be a lowercase 40-character SHA"))
    emitted = parse_time(checkpoint["emitted_at"])
    due = parse_time(checkpoint["callback_due_at"])
    if emitted is None:
        errors.append(error("invalid_timestamp", "emitted_at", "must be an offset-aware ISO-8601 timestamp"))
    if due is None:
        errors.append(error("invalid_timestamp", "callback_due_at", "must be an offset-aware ISO-8601 timestamp"))
    if checkpoint["event_type"] == "heartbeat" and emitted is not None and due is not None and emitted <= due:
        errors.append(error("heartbeat_not_overdue", "emitted_at", "heartbeat is allowed only strictly after callback_due_at"))
    for field in ("changed_paths", "evidence", "blockers"):
        value = checkpoint[field]
        if not isinstance(value, list) or not all(isinstance(item, str) for item in value):
            errors.append(error("invalid_field", field, "must be an array of strings"))
    if checkpoint["review"] is not None and not isinstance(checkpoint["review"], dict):
        errors.append(error("invalid_field", "review", "must be an object or null"))

    if prior_state is not None:
        seen_ids, previous_sequences = validate_prior_state(prior_state, errors)
        if checkpoint["event_id"] in seen_ids:
            errors.append(error("duplicate_event_id", "event_id", "was already accepted in prior state"))
        if previous_sequences and isinstance(checkpoint["sequence"], int) and checkpoint["sequence"] <= max(previous_sequences):
            errors.append(error("stale_sequence", "sequence", "must be greater than prior maximum"))
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--checkpoint", required=True, type=Path)
    parser.add_argument("--prior-state", type=Path)
    args = parser.parse_args()
    try:
        checkpoint = json.loads(args.checkpoint.read_text(encoding="utf-8"))
        prior_state = json.loads(args.prior_state.read_text(encoding="utf-8")) if args.prior_state else None
    except (OSError, json.JSONDecodeError) as exc:
        print(json.dumps({"ok": False, "errors": [error("input_error", "$", str(exc))]}, sort_keys=True))
        return 2
    errors = validate_checkpoint(checkpoint, prior_state)
    print(json.dumps({"ok": not errors, "errors": errors}, sort_keys=True))
    return 0 if not errors else 1


if __name__ == "__main__":
    sys.exit(main())
