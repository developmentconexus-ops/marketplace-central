"""Deterministic tests for the native harness control artifacts."""

from __future__ import annotations

import importlib.util
import json
import tomllib
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
PROJECT_ROOT = SKILL_ROOT.parents[2]
CODEX_CONFIG = PROJECT_ROOT / ".codex" / "config.toml"
MILESTONE_AGENT = PROJECT_ROOT / ".codex" / "agents" / "mpc-milestone.toml"
FIXTURES = Path(__file__).with_name("fixtures")


def load_module(name: str, path: Path):
    spec = importlib.util.spec_from_file_location(name, path)
    assert spec and spec.loader
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


CHECKPOINT = load_module("checkpoint_validator", SKILL_ROOT / "scripts" / "validate_checkpoint.py")
PREFLIGHT = load_module("dispatch_preflight", SKILL_ROOT / "scripts" / "dispatch_preflight.py")


def fixture(name: str):
    return json.loads((FIXTURES / name).read_text(encoding="utf-8"))


def codes(errors):
    return {item["code"] for item in errors}


ACCEPTED_BASE_SHA = fixture("valid_dispatch.json")["base_sha"]


def preflight(packet, current_writers=None, accepted_base_sha=ACCEPTED_BASE_SHA):
    if current_writers is None:
        current_writers = fixture("current_writers_empty.json")
    return PREFLIGHT.validate_packet(packet, SKILL_ROOT.parents[2], accepted_base_sha, current_writers)


class CheckpointValidationTests(unittest.TestCase):
    def test_valid_checkpoint(self):
        self.assertEqual([], CHECKPOINT.validate_checkpoint(fixture("valid_checkpoint.json")))

    def test_duplicate_event_is_rejected(self):
        errors = CHECKPOINT.validate_checkpoint(
            fixture("valid_checkpoint.json"), fixture("prior_duplicate_event.json")
        )
        self.assertIn("duplicate_event_id", codes(errors))

    def test_stale_sequence_is_rejected(self):
        errors = CHECKPOINT.validate_checkpoint(
            fixture("valid_checkpoint.json"), fixture("prior_stale_sequence.json")
        )
        self.assertIn("stale_sequence", codes(errors))

    def test_missing_required_field_is_rejected(self):
        checkpoint = fixture("valid_checkpoint.json")
        del checkpoint["event_id"]
        errors = CHECKPOINT.validate_checkpoint(checkpoint)
        self.assertIn("missing_required_field", codes(errors))

    def test_malformed_prior_state_is_rejected(self):
        errors = CHECKPOINT.validate_checkpoint(
            fixture("valid_checkpoint.json"), {"events": [None], "max_sequence": 0}
        )
        self.assertIn("invalid_prior_state", codes(errors))

    def test_heartbeat_is_only_an_overdue_fallback(self):
        heartbeat = fixture("early_heartbeat.json")
        self.assertIn("heartbeat_not_overdue", codes(CHECKPOINT.validate_checkpoint(heartbeat)))
        heartbeat["emitted_at"] = "2026-07-12T12:10:00Z"
        self.assertIn("heartbeat_not_overdue", codes(CHECKPOINT.validate_checkpoint(heartbeat)))
        heartbeat["emitted_at"] = "2026-07-12T12:10:01Z"
        self.assertEqual([], CHECKPOINT.validate_checkpoint(heartbeat))


class DispatchPreflightTests(unittest.TestCase):
    def test_valid_packet(self):
        self.assertEqual([], preflight(fixture("valid_dispatch.json")))

    def test_base_sha_mismatch_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        packet["base_sha"] = packet["accepted_base_sha"] = "a" * 40
        self.assertIn("base_sha_mismatch", codes(preflight(packet)))

    def test_nonexistent_milestone_accepted_sha_is_rejected(self):
        errors = preflight(fixture("valid_dispatch.json"), accepted_base_sha="a" * 40)
        self.assertIn("accepted_base_sha_not_found", codes(errors))

    def test_completed_work_item_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        packet["completed_work_item_ids"] = ["F-01"]
        self.assertIn("duplicate_completed_work_item", codes(preflight(packet)))

    def test_missing_feature_identity_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        del packet["feature_id"]
        self.assertIn("missing_required_field", codes(preflight(packet)))

    def test_missing_documented_feature_packet_field_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        del packet["context_files"]
        self.assertIn("missing_required_field", codes(preflight(packet)))

    def test_missing_documented_feature_file_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        packet["feature_file"] = "does/not/exist.md"
        self.assertIn("missing_required_file", codes(preflight(packet)))

    def test_overlapping_writer_path_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        writers = fixture("current_writers_overlapping.json")
        self.assertIn("overlapping_writer_path", codes(preflight(packet, writers)))

    def test_missing_authoritative_writer_state_is_rejected(self):
        errors = PREFLIGHT.validate_packet(fixture("valid_dispatch.json"), SKILL_ROOT.parents[2], ACCEPTED_BASE_SHA)
        self.assertIn("missing_authoritative_writer_state", codes(errors))

    def test_missing_packet_file_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        packet["required_files"] = ["does/not/exist.md"]
        self.assertIn("missing_required_file", codes(preflight(packet)))

    def test_malformed_completed_marker_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        packet["completed_work_item_ids"] = {"F-01": True}
        self.assertIn("invalid_field", codes(preflight(packet)))

    def test_malformed_dispatch_state_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        packet["dispatch_state"] = {"state": "ready"}
        self.assertIn("invalid_marker_state", codes(preflight(packet)))

    def test_active_dispatch_state_is_rejected(self):
        packet = fixture("valid_dispatch.json")
        packet["dispatch_state"] = "dispatched"
        self.assertIn("duplicate_active_dispatch", codes(preflight(packet)))


class MilestoneAgentConfigurationTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config = tomllib.loads(CODEX_CONFIG.read_text(encoding="utf-8"))
        cls.agent = tomllib.loads(MILESTONE_AGENT.read_text(encoding="utf-8"))
        cls.skill = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")

    def test_native_tree_is_bounded_to_portfolio_milestone_and_children(self):
        self.assertTrue(self.config["features"]["multi_agent"])
        self.assertEqual(2, self.config["agents"]["max_depth"])
        self.assertEqual(4, self.config["agents"]["max_threads"])

    def test_only_mpc_milestone_custom_role_is_registered(self):
        registered = {
            key for key, value in self.config["agents"].items() if isinstance(value, dict)
        }
        self.assertEqual({"mpc-milestone"}, registered)
        role = self.config["agents"]["mpc-milestone"]
        self.assertEqual("agents/mpc-milestone.toml", role["config_file"])

    def test_custom_agent_has_required_official_fields(self):
        self.assertEqual("mpc-milestone", self.agent["name"])
        self.assertTrue(self.agent["description"].strip())
        self.assertTrue(self.agent["developer_instructions"].strip())
        self.assertEqual("gpt-5.6-terra", self.agent["model"])
        self.assertEqual("high", self.agent["model_reasoning_effort"])

    def test_parent_channel_is_compact_and_outcome_only(self):
        prompt = self.agent["developer_instructions"]
        self.assertIn("`needs_input`", prompt)
        self.assertIn("`terminal`", prompt)
        self.assertIn("under 2,000 characters", prompt)
        self.assertIn("Never include raw logs or transcripts", prompt)
        self.assertIn("needs_input", self.skill)
        self.assertIn("terminal", self.skill)

    def test_config_adds_no_hooks_or_external_capabilities(self):
        forbidden = {
            "hooks", "mcp_servers", "model_provider", "model_providers",
            "approval_policy", "sandbox_mode", "sandbox_workspace_write",
            "network_access", "otel", "notify",
        }
        self.assertTrue(forbidden.isdisjoint(self.config))


if __name__ == "__main__":
    unittest.main()
