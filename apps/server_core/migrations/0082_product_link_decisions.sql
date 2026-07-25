-- 0082: product_link_decisions — E10 auto-link audit trail (MIS-006 M-05).
-- One row per approval decision on a product link, written in the SAME
-- transaction as the link's state transition: an approved link without its
-- decision row, or the reverse, would make the trail a guess rather than a
-- record.
--
-- Bloco B+: runs after M-02's products_mirror/active_source (bloco B). Additive
-- only — product_links is NOT altered. A8 (idempotency) is already satisfied by
-- that table's existing PRIMARY KEY
-- (tenant_id, installation_id, provider_item_id, provider_variation_id); a new
-- key on (internal_product_id, ...) would have collapsed two variations of the
-- same listing and capped a product at one listing, both of which are wrong
-- (D-121: one product may link to N distinct listings, unflagged).
--
-- actor distinguishes the only automatic path from every human one. D-121-2:
-- the automatic path is corroboration alone — CODPROD and EAN naming the same
-- product (concordant_codprod_ean). A single anchor is an operator decision, so
-- exact_codprod_unique/exact_ean_unique only ever carry actor='operator'; that
-- pairing is enforced below, not left to application discipline.
--
-- ADR-17: collisions_at_decision is NULLABLE. It records how many products the
-- deciding anchor matched, and it is written only when that count was actually
-- read. A decision taken without reading it stores NULL — never 0, which would
-- read as "matched nothing" and is a fact nobody established.
CREATE TABLE IF NOT EXISTS product_link_decisions (
    tenant_id              text        NOT NULL,
    decision_id            text        NOT NULL,
    installation_id        text        NOT NULL,
    provider_item_id       text        NOT NULL,
    provider_variation_id  text        NOT NULL DEFAULT '',
    -- The contract (§E10) names the subject of a decision `link_id`.
    -- product_links has no surrogate id, so the link IS its natural key. The
    -- column is generated from that key instead of written alongside it, so the
    -- audit trail cannot drift from the link it claims to describe.
    link_id                text        GENERATED ALWAYS AS (
                               installation_id || '|' || provider_item_id || '|' || provider_variation_id
                           ) STORED,
    rule_matched           text        NOT NULL,
    actor                  text        NOT NULL,
    collisions_at_decision integer,
    created_at             timestamptz NOT NULL,
    -- Set on the OLD row when a later decision replaces it. History is
    -- append-only: an operator override writes a new row and stamps the
    -- previous one, never rewrites it in place.
    superseded_by          text,
    PRIMARY KEY (tenant_id, decision_id),
    CONSTRAINT product_link_decisions_actor_check
        CHECK (actor IN ('system', 'operator')),
    CONSTRAINT product_link_decisions_rule_check
        CHECK (rule_matched IN ('exact_codprod_unique', 'exact_ean_unique', 'concordant_codprod_ean', 'manual')),
    -- Corroboration is the only thing the system may decide by itself; every
    -- single-anchor rule is a human confirmation (D-121-2).
    CONSTRAINT product_link_decisions_system_is_corroborated_check
        CHECK (actor <> 'system' OR rule_matched = 'concordant_codprod_ean'),
    -- An anchor that matched nothing could not have decided anything, so 0 is
    -- not a possible reading — only a real count or an honest unknown.
    CONSTRAINT product_link_decisions_collisions_check
        CHECK (collisions_at_decision IS NULL OR collisions_at_decision >= 1),
    CONSTRAINT product_link_decisions_not_self_superseded_check
        CHECK (superseded_by IS NULL OR superseded_by <> decision_id),
    CONSTRAINT product_link_decisions_superseded_by_fkey
        FOREIGN KEY (tenant_id, superseded_by) REFERENCES product_link_decisions (tenant_id, decision_id)
);

-- The read this table exists for: the decision history of one link, newest
-- first. Partial index on the live row answers "who decided this link, and was
-- it a human?" — the precedence check the automatic path must pass before it
-- may touch a link at all.
CREATE INDEX IF NOT EXISTS product_link_decisions_link_idx
    ON product_link_decisions (tenant_id, link_id, created_at DESC);
CREATE INDEX IF NOT EXISTS product_link_decisions_current_idx
    ON product_link_decisions (tenant_id, link_id)
    WHERE superseded_by IS NULL;
