#!/usr/bin/env python3
"""Read-only fail-closed preflight for a native child dispatch packet."""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from pathlib import Path, PurePosixPath
from typing import Any


SHA_RE = re.compile(r"^[0-9a-f]{40}$")
REQUIRED_FIELDS = (
    "role", "base_sha", "accepted_base_sha", "dispatch_id", "feature_id", "work_item_id",
    "source_task_id", "parent_task_id", "callback_task_id", "required_files",
    "allowed_paths", "dirty_paths", "one_writer", "milestone_task_id", "feature_file",
    "context_files", "knowledge_routes", "forbidden_paths", "shared_seams",
    "side_effects", "proof", "stop_conditions",
)


def error(code: str, field: str, detail: str) -> dict[str, str]:
    return {"code": code, "field": field, "detail": detail}


def clean_relative_path(value: Any) -> str | None:
    if not isinstance(value, str) or not value.strip():
        return None
    normalized = value.replace("\\", "/")
    path = PurePosixPath(normalized)
    if path.is_absolute() or ".." in path.parts:
        return None
    return path.as_posix().rstrip("/")


def overlaps(left: str, right: str) -> bool:
    return left == right or left.startswith(right + "/") or right.startswith(left + "/")


def path_is_owned(path: str, owned_paths: list[str]) -> bool:
    return any(path == owned or path.startswith(owned + "/") for owned in owned_paths)


def as_paths(value: Any, field: str, errors: list[dict[str, str]]) -> list[str]:
    if not isinstance(value, list):
        errors.append(error("invalid_field", field, "must be an array of relative paths"))
        return []
    parsed: list[str] = []
    for index, item in enumerate(value):
        path = clean_relative_path(item)
        if path is None:
            errors.append(error("invalid_path", f"{field}[{index}]", "must be a non-empty repository-relative path"))
        else:
            parsed.append(path)
    return parsed


def as_strings(value: Any, field: str, errors: list[dict[str, str]], require_nonempty: bool = False) -> list[str]:
    if not isinstance(value, list) or not all(isinstance(item, str) and item.strip() for item in value):
        errors.append(error("invalid_field", field, "must be an array of non-empty strings"))
        return []
    if require_nonempty and not value:
        errors.append(error("missing_required_value", field, "must not be empty"))
    return value


def checked_files(value: Any, field: str, repo_root: Path, errors: list[dict[str, str]], require_nonempty: bool = False) -> list[str]:
    paths = as_paths(value, field, errors)
    if require_nonempty and not paths:
        errors.append(error("missing_required_value", field, "must name at least one file"))
    root = repo_root.resolve()
    for relative in paths:
        candidate = (root / relative).resolve()
        try:
            candidate.relative_to(root)
        except ValueError:
            errors.append(error("invalid_path", field, "must remain inside repo_root"))
            continue
        if not candidate.is_file():
            errors.append(error("missing_required_file", relative, "does not exist"))
    return paths


def exact_string_lists(value: Any, field: str, keys: tuple[str, ...], errors: list[dict[str, str]], require_nonempty: bool = False) -> None:
    if not isinstance(value, dict):
        errors.append(error("invalid_field", field, "must be an object"))
        return
    for key in keys:
        as_strings(value.get(key), f"{field}.{key}", errors, require_nonempty)


def git_has_commit(repo_root: Path, sha: str) -> bool:
    """Read Git object existence only; never change refs, index, or worktree."""
    try:
        result = subprocess.run(
            ["git", "-C", str(repo_root), "cat-file", "-e", f"{sha}^{{commit}}"],
            check=False,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            timeout=5,
        )
    except (OSError, subprocess.TimeoutExpired):
        return False
    return result.returncode == 0


def authoritative_writers(writer_state: Any, errors: list[dict[str, str]]) -> list[tuple[str, list[str]]]:
    if not isinstance(writer_state, dict) or writer_state.get("authoritative") is not True:
        errors.append(error("missing_authoritative_writer_state", "current_writers", "must be an authoritative writer snapshot"))
        return []
    entries = writer_state.get("writers")
    if not isinstance(entries, list):
        errors.append(error("invalid_field", "current_writers.writers", "must be an array"))
        return []
    writers: list[tuple[str, list[str]]] = []
    for index, entry in enumerate(entries):
        if not isinstance(entry, dict):
            errors.append(error("invalid_field", f"current_writers.writers[{index}]", "must be an object"))
            continue
        owner = entry.get("owner_task_id")
        if not isinstance(owner, str) or not owner.strip():
            errors.append(error("invalid_field", f"current_writers.writers[{index}].owner_task_id", "must be a non-empty string"))
            continue
        writers.append((owner, as_paths(entry.get("paths"), f"current_writers.writers[{index}].paths", errors)))
    return writers


def validate_packet(
    packet: Any,
    repo_root: Path,
    accepted_base_sha: Any = None,
    current_writers: Any = None,
) -> list[dict[str, str]]:
    """Validate packet-only declarations and read file existence without writes."""
    errors: list[dict[str, str]] = []
    if not isinstance(packet, dict):
        return [error("invalid_packet", "$", "must be a JSON object")]
    for field in REQUIRED_FIELDS:
        if field not in packet:
            errors.append(error("missing_required_field", field, "is required"))
    if errors:
        return errors
    for field in ("role", "dispatch_id", "feature_id", "work_item_id", "source_task_id", "parent_task_id", "callback_task_id", "milestone_task_id"):
        if not isinstance(packet[field], str) or not packet[field].strip():
            errors.append(error("invalid_field", field, "must be a non-empty string"))
    if packet["role"] != "Feature Implementer":
        errors.append(error("invalid_role", "role", "must be Feature Implementer"))
    for field in ("base_sha", "accepted_base_sha"):
        if not isinstance(packet[field], str) or not SHA_RE.fullmatch(packet[field]):
            errors.append(error("invalid_sha", field, "must be a lowercase 40-character SHA"))
    if not isinstance(accepted_base_sha, str) or not SHA_RE.fullmatch(accepted_base_sha):
        errors.append(error("invalid_accepted_base_sha", "accepted_base_sha", "must be supplied by the Milestone as a lowercase 40-character SHA"))
    elif not git_has_commit(repo_root, accepted_base_sha):
        errors.append(error("accepted_base_sha_not_found", "accepted_base_sha", "is not a commit in repo_root"))
    if packet["base_sha"] != accepted_base_sha or packet["accepted_base_sha"] != accepted_base_sha:
        errors.append(error("base_sha_mismatch", "base_sha", "does not match the Milestone-supplied accepted base SHA"))
    if packet["callback_task_id"] != packet["parent_task_id"]:
        errors.append(error("callback_identity_mismatch", "callback_task_id", "must return to parent_task_id"))
    if packet["milestone_task_id"] != packet["parent_task_id"]:
        errors.append(error("milestone_identity_mismatch", "milestone_task_id", "must equal parent_task_id"))

    checked_files(packet["required_files"], "required_files", repo_root, errors, require_nonempty=True)
    feature_files = checked_files([packet["feature_file"]], "feature_file", repo_root, errors, require_nonempty=True)
    if not feature_files:
        errors.append(error("invalid_field", "feature_file", "must be a repository-relative file path"))
    checked_files(packet["context_files"], "context_files", repo_root, errors, require_nonempty=True)
    as_strings(packet["knowledge_routes"], "knowledge_routes", errors, require_nonempty=True)
    allowed_paths = as_paths(packet["allowed_paths"], "allowed_paths", errors)
    dirty_paths = as_paths(packet["dirty_paths"], "dirty_paths", errors)
    as_paths(packet["forbidden_paths"], "forbidden_paths", errors)
    as_paths(packet["shared_seams"], "shared_seams", errors)
    exact_string_lists(packet["side_effects"], "side_effects", ("allowed", "forbidden"), errors)
    exact_string_lists(packet["proof"], "proof", ("command_ids", "evidence_targets"), errors, require_nonempty=True)
    as_strings(packet["stop_conditions"], "stop_conditions", errors, require_nonempty=True)

    one_writer = packet["one_writer"]
    if not isinstance(one_writer, dict):
        errors.append(error("missing_one_writer_declaration", "one_writer", "must name one owner and owned paths"))
        owned_paths: list[str] = []
    else:
        if one_writer.get("owner_task_id") != packet["source_task_id"]:
            errors.append(error("one_writer_owner_mismatch", "one_writer.owner_task_id", "must equal source_task_id"))
        owned_paths = as_paths(one_writer.get("owned_paths"), "one_writer.owned_paths", errors)
        if not owned_paths:
            errors.append(error("missing_one_writer_declaration", "one_writer.owned_paths", "must declare at least one path"))
    for path in allowed_paths:
        if not path_is_owned(path, owned_paths):
            errors.append(error("unowned_allowed_path", path, "is not covered by one_writer.owned_paths"))
    for path in dirty_paths:
        if not path_is_owned(path, owned_paths):
            errors.append(error("unowned_dirty_path", path, "is not covered by one_writer.owned_paths"))
    if "dirty_owner_task_id" in packet and packet["dirty_owner_task_id"] != packet["source_task_id"]:
        errors.append(error("dirty_owner_mismatch", "dirty_owner_task_id", "must equal source_task_id"))

    marker_fields = (
        "completed_work_item_ids", "completed_dispatch_ids", "stale_work_item_ids", "stale_dispatch_ids",
    )
    marker_values: dict[str, list[str]] = {}
    for field in marker_fields:
        if field in packet:
            marker_values[field] = as_strings(packet[field], field, errors)
        else:
            marker_values[field] = []
    if "dispatch_state" in packet and (
        not isinstance(packet["dispatch_state"], str)
        or packet["dispatch_state"] not in {"ready", "pending", "dispatched", "completed", "stale", "cancelled"}
    ):
        errors.append(error("invalid_marker_state", "dispatch_state", "must be a supported dispatch state"))
    completed_work = marker_values["completed_work_item_ids"]
    completed_dispatch = marker_values["completed_dispatch_ids"]
    stale_work = marker_values["stale_work_item_ids"]
    stale_dispatch = marker_values["stale_dispatch_ids"]
    if packet["work_item_id"] in completed_work:
        errors.append(error("duplicate_completed_work_item", "work_item_id", "is already completed"))
    if isinstance(completed_dispatch, list) and packet["dispatch_id"] in completed_dispatch:
        errors.append(error("duplicate_completed_dispatch", "dispatch_id", "is already completed"))
    if isinstance(stale_work, list) and packet["work_item_id"] in stale_work:
        errors.append(error("stale_work_item", "work_item_id", "is marked stale"))
    if isinstance(stale_dispatch, list) and packet["dispatch_id"] in stale_dispatch:
        errors.append(error("stale_dispatch", "dispatch_id", "is marked stale"))
    if packet.get("dispatch_state") == "dispatched":
        errors.append(error("duplicate_active_dispatch", "dispatch_state", "is already active"))
    if isinstance(packet.get("dispatch_state"), str) and packet["dispatch_state"] in {"completed", "stale", "cancelled"}:
        errors.append(error("stale_dispatch", "dispatch_state", "cannot be dispatched again"))

    one_writer_owner = one_writer.get("owner_task_id") if isinstance(one_writer, dict) else None
    for owner, other_paths in authoritative_writers(current_writers, errors):
        if owner == one_writer_owner:
            continue
        for own in owned_paths:
            for other in other_paths:
                if overlaps(own, other):
                    errors.append(error("overlapping_writer_path", other, "overlaps current_writers"))
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--packet", required=True, type=Path)
    parser.add_argument("--repo-root", type=Path, default=Path.cwd())
    parser.add_argument("--accepted-base-sha", required=True)
    parser.add_argument("--current-writers", required=True, type=Path)
    args = parser.parse_args()
    try:
        packet = json.loads(args.packet.read_text(encoding="utf-8"))
        current_writers = json.loads(args.current_writers.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        print(json.dumps({"ok": False, "errors": [error("input_error", "$", str(exc))]}, sort_keys=True))
        return 2
    errors = validate_packet(packet, args.repo_root, args.accepted_base_sha, current_writers)
    print(json.dumps({"ok": not errors, "errors": errors}, sort_keys=True))
    return 0 if not errors else 1


if __name__ == "__main__":
    sys.exit(main())
