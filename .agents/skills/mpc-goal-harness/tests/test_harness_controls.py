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
FIXTURES = Path(__file__).with_name("fixtures")


def load_module(name: str, path: Path):
    spec = importlib.util.spec_from_file_location(name, path)
    assert spec and spec.loader
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


CHECKPOINT = load_module(
    "checkpoint_validator", SKILL_ROOT / "scripts" / "validate_checkpoint.py"
)
PREFLIGHT = load_module(
    "dispatch_preflight", SKILL_ROOT / "scripts" / "dispatch_preflight.py"
)


def fixture(name: str):
    return json.loads((FIXTURES / name).read_text(encoding="utf-8"))


def codes(errors):
    return {item["code"] for item in errors}


ACCEPTED_BASE_SHA = fixture("valid_dispatch.json")["base_sha"]


def preflight(packet, current_writers=None, accepted_base_sha=ACCEPTED_BASE_SHA):
    if current_writers is None:
        current_writers = fixture("current_writers_empty.json")
    return PREFLIGHT.validate_packet(
        packet, SKILL_ROOT.parents[2], accepted_base_sha, current_writers
    )


class CheckpointValidationTests(unittest.TestCase):
    def test_valid_checkpoint(self):
        self.assertEqual(
            [], CHECKPOINT.validate_checkpoint(fixture("valid_checkpoint.json"))
        )

    def test_duplicate_event_is_rejected(self):
        errors = CHECKPOINT.validate_checkpoint(
            fixture("valid_checkpoint.json"),
            fixture("prior_duplicate_event.json"),
        )
        self.assertIn("duplicate_event_id", codes(errors))

    def test_stale_sequence_is_rejected(self):
        errors = CHECKPOINT.validate_checkpoint(
            fixture("valid_checkpoint.json"),
            fixture("prior_stale_sequence.json"),
        )
        self.assertIn("stale_sequence", codes(errors))

    def test_missing_required_field_is_rejected(self):
        checkpoint = fixture("valid_checkpoint.json")
        del checkpoint["event_id"]
        errors = CHECKPOINT.validate_checkpoint(checkpoint)
        self.assertIn("missing_required_field", codes(errors))

    def test_malformed_prior_state_is_rejected(self):
        errors = CHECKPOINT.validate_checkpoint(
            fixture("valid_checkpoint.json"),
            {"events": [None], "max_sequence": 0},
        )
        self.assertIn("invalid_prior_state", codes(errors))

    def test_heartbeat_is_only_an_overdue_fallback(self):
        heartbeat = fixture("early_heartbeat.json")
        self.assertIn(
            "heartbeat_not_overdue",
            codes(CHECKPOINT.validate_checkpoint(heartbeat)),
        )
        heartbeat["emitted_at"] = "2026-07-12T12:10:00Z"
        self.assertIn(
            "heartbeat_not_overdue",
            codes(CHECKPOINT.validate_checkpoint(heartbeat)),
        )
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
        errors = preflight(
            fixture("valid_dispatch.json"),
            accepted_base_sha="a" * 40,
        )
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
        self.assertIn(
            "overlapping_writer_path",
            codes(preflight(packet, writers)),
        )

    def test_missing_authoritative_writer_state_is_rejected(self):
        errors = PREFLIGHT.validate_packet(
            fixture("valid_dispatch.json"),
            SKILL_ROOT.parents[2],
            ACCEPTED_BASE_SHA,
        )
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


class MilestoneRuntimeConfigurationTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.config = tomllib.loads(CODEX_CONFIG.read_text(encoding="utf-8"))
        cls.skill = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")

    def test_standalone_milestone_tree_is_bounded_to_direct_children(self):
        self.assertTrue(self.config["features"]["multi_agent"])
        self.assertEqual(1, self.config["agents"]["max_depth"])
        self.assertEqual(3, self.config["agents"]["max_threads"])

    def test_no_unselectable_custom_roles_are_registered(self):
        registered = {
            key
            for key, value in self.config["agents"].items()
            if isinstance(value, dict)
        }
        self.assertEqual(set(), registered)
        self.assertIn(
            "no agent_type, model, or reasoning selector",
            self.skill,
        )
        self.assertIn(
            "never use task_name as evidence",
            self.skill,
        )

    def test_manual_root_and_generic_children_are_explicit(self):
        self.assertIn(
            "session_type: manually_created_standalone_visible_root",
            self.skill,
        )
        self.assertIn("role: Milestone Orchestrator", self.skill)
        self.assertIn("model_to_select: gpt-5.6-terra", self.skill)
        self.assertIn("reasoning_effort_to_select: medium", self.skill)
        self.assertIn(
            "harness_skill: .agents/skills/mpc-goal-harness/SKILL.md",
            self.skill,
        )
        self.assertIn("runtime_agent_type: generic", self.skill)
        self.assertIn("runtime_model_policy: runtime_managed", self.skill)
        self.assertIn(
            "task_name is a correlation label, not an agent",
            self.skill,
        )
        self.assertNotIn("agent_type: mpc-milestone", self.skill)
        self.assertNotIn("agent_type: mpc-implementer", self.skill)
        self.assertNotIn("agent_type: mpc-verifier", self.skill)

    def test_parent_channel_is_compact_and_outcome_only(self):
        self.assertIn("needs_input", self.skill)
        self.assertIn("terminal", self.skill)
        self.assertIn(
            "For terminal, persist and validate it before messaging",
            self.skill,
        )
        self.assertIn("under 2,000 characters", self.skill)
        self.assertIn("blockers:", self.skill)
        self.assertIn("next:", self.skill)
        self.assertIn("No transcript replay", self.skill)

    def test_default_child_budget_prevents_microfeature_fanout(self):
        self.assertIn(
            "one coherent Feature Implementer run, one final fixed-SHA",
            self.skill,
        )
        self.assertIn(
            "review, and one proportional QA run",
            self.skill,
        )
        self.assertIn("compact human cost decision", self.skill)

    def test_normal_control_plane_forbids_heartbeat_and_polling(self):
        self.assertIn(
            "No heartbeat, callback guard, cron, hook, polling loop",
            self.skill,
        )
        self.assertIn(
            "Legacy heartbeat remains schema-valid",
            self.skill,
        )
        self.assertIn("new dispatches must not create it", self.skill)
        self.assertNotIn("must create one native Codex heartbeat", self.skill)

    def test_terminal_callback_is_explicit_and_durable_first(self):
        self.assertIn(
            "persist and validate one compact checkpoint, then explicitly call",
            self.skill,
        )
        self.assertIn(
            "send_message_to_thread with portfolio_task_id",
            self.skill,
        )
        self.assertIn("Native completion is not a callback", self.skill)
        self.assertIn(
            "checkpoint: <repository-relative checkpoint path>",
            self.skill,
        )
        self.assertIn("return a pasteable callback", self.skill)

    def test_portfolio_handoff_is_manual(self):
        self.assertIn("Portfolio never creates that task", self.skill)
        self.assertIn(
            "user, not Portfolio, creates the new task",
            self.skill,
        )
        self.assertIn("After emitting the prompt, Portfolio stops", self.skill)
        self.assertIn(
            "It does not poll or create the\nMilestone",
            self.skill,
        )

    def test_config_adds_no_hooks_or_external_capabilities(self):
        forbidden = {
            "hooks",
            "mcp_servers",
            "model_provider",
            "model_providers",
            "approval_policy",
            "sandbox_mode",
            "sandbox_workspace_write",
            "network_access",
            "otel",
            "notify",
        }
        self.assertTrue(forbidden.isdisjoint(self.config))


if __name__ == "__main__":
    unittest.main()
