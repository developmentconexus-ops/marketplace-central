ALTER TABLE profitability_manual_adjustments
    ADD CONSTRAINT profitability_manual_adjustments_actor_type_nonempty
        CHECK (btrim(actor_type) <> '') NOT VALID,
    ADD CONSTRAINT profitability_manual_adjustments_actor_id_nonempty
        CHECK (btrim(actor_id) <> '') NOT VALID,
    ADD CONSTRAINT profitability_manual_adjustments_reason_nonempty
        CHECK (btrim(reason) <> '') NOT VALID,
    ADD CONSTRAINT profitability_manual_adjustments_currency_nonempty
        CHECK (btrim(currency) <> '') NOT VALID,
    ADD CONSTRAINT profitability_manual_adjustments_scope_valid
        CHECK (scope IN ('order', 'item')) NOT VALID,
    ADD CONSTRAINT profitability_manual_adjustments_category_valid
        CHECK (category IN ('freight', 'commission', 'cost', 'generic_adjustment')) NOT VALID,
    ADD CONSTRAINT profitability_manual_adjustments_scope_identity_valid
        CHECK (
            (scope = 'order' AND btrim(provider_item_id) = '' AND btrim(provider_variation_id) = '')
            OR (scope = 'item' AND btrim(provider_item_id) <> '')
        ) NOT VALID;
