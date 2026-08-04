-- 0093: ICMS matrix schema (P2.b, Task 1) — imposto ex-ante no motor de preco.
--
-- Adds the fiscal fields the pricing engine needs to read ICMS ex-ante instead
-- of fabricating it. These five columns land on products_mirror because ICMS
-- treatment is a per-product classification (grupo_icms/origprod) plus a
-- per-product fiscal snapshot (st_retido_entrada/restituicao_unit/fiscal_dt_ref),
-- not a per-order fact — the sync adapters populate them in a later task
-- (Task 2). Honest-unknown (ADR-17): every column here is bare nullable, no
-- DEFAULT — a product whose fiscal classification was never read stays NULL,
-- never 0 or a plausible guess.
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS grupo_icms          integer;
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS origprod            smallint;
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS st_retido_entrada   numeric(14,4);
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS restituicao_unit    numeric(14,4);
ALTER TABLE products_mirror ADD COLUMN IF NOT EXISTS fiscal_dt_ref       date;

-- icms_aliquota_interna is a LEGAL table: it is OUR data, seeded from the
-- states' own legislation (D-42) — it is never populated from TGFICM, which
-- carries stale/wrong internal rates for at least AL/DF/PR/RJ (see the seed
-- comment below the table). fcp_embutido is already folded into the state's
-- headline rate wherever the FCP is a general, non-segment-specific surcharge
-- (D-43); a segment-specific FCP is out of scope here. One row per UF per
-- vigencia: the partial unique index below is what actually guarantees at
-- most one open (vigente_ate IS NULL) row per UF — never rely on writer
-- discipline alone for that guarantee.
CREATE TABLE IF NOT EXISTS icms_aliquota_interna (
    uf             text        NOT NULL,
    aliquota       numeric(6,3) NOT NULL,
    fcp_embutido   numeric(6,3) NOT NULL DEFAULT 0,
    fonte          text        NOT NULL,
    lei            text        NOT NULL,
    vigente_desde  date        NOT NULL,
    vigente_ate    date,
    PRIMARY KEY (uf, vigente_desde)
);

CREATE UNIQUE INDEX IF NOT EXISTS icms_aliquota_interna_vigente
    ON icms_aliquota_interna (uf)
    WHERE vigente_ate IS NULL;

-- icms_matrix_mirror is the synced (uf_origem, uf_destino, grupo_icms) ICMS
-- cell, mirrored from the ERP by a later task. codtrib is nullable (a cell can
-- be read before its tributary classification is known). This slice's
-- calculation reads ONLY codtrib (is this destination ST?) — aliquota,
-- red_base and perc_fcp are still captured verbatim here so P2.c can
-- reconcile against the emitted nota, but no calculation in this task reads
-- them. red_base is a base-composing fact (it is nonzero in 63% of rows), not
-- part of the rate, which is exactly why it is its own column instead of
-- being folded into aliquota. The partial unique index is the actual
-- guarantee against two open versions of the same cell coexisting; the
-- writer's discipline is not.
CREATE TABLE IF NOT EXISTS icms_matrix_mirror (
    tenant_id          text        NOT NULL,
    uf_origem          text        NOT NULL,
    uf_destino         text        NOT NULL,
    grupo_icms         integer     NOT NULL,
    codtrib            smallint,
    aliquota           numeric(6,3),
    red_base           numeric(6,3),
    perc_fcp           numeric(6,3),
    linhas_candidatas  smallint    NOT NULL,
    ambiguo            boolean     NOT NULL,
    vigente_desde      timestamptz NOT NULL,
    vigente_ate        timestamptz,
    synced_at          timestamptz NOT NULL,
    PRIMARY KEY (tenant_id, uf_origem, uf_destino, grupo_icms, vigente_desde)
);

CREATE UNIQUE INDEX IF NOT EXISTS icms_matrix_mirror_vigente
    ON icms_matrix_mirror (tenant_id, uf_origem, uf_destino, grupo_icms)
    WHERE vigente_ate IS NULL;

-- Seed: 27 UFs, verified against the legislation in effect on 2026-08-03.
-- fcp_embutido is nonzero only where the FCP is folded into the headline
-- rate as a general surcharge (AL, RJ, SE). AL/DF/PR/RJ carry a specific
-- vigente_desde because these are exactly the four UFs where the ERP's own
-- TGFICM rate is stale or wrong (D-42/D-43) — fonte/lei are columns, not a
-- comment, so a reconciliation can cite the law, not prose.
-- ON CONFLICT DO NOTHING porque schema_migrations é chaveado por NOME DE
-- ARQUIVO (platform/migrate/runner.go): esta migração já rodou como
-- 0093_icms_matrix.sql em ambientes vivos, e o renome para 0094 (colisão de
-- prefixo entre três fatias paralelas) faz o runner reaplicá-la. DO NOTHING
-- mantém o seed já semeado intacto — nunca sobrescreve uma alíquota que uma
-- migração posterior tenha corrigido.
INSERT INTO icms_aliquota_interna (uf, aliquota, fcp_embutido, fonte, lei, vigente_desde) VALUES
    ('AC', 19.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('AL', 21.5, 1.0, 'Lei 9.776/2025 (20,5% ICMS + 1% FECOEP = 21,5%)', 'Lei 9.776/2025', '2026-04-01'),
    ('AM', 20.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('AP', 18.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('BA', 20.5, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('CE', 20.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('DF', 20.0, 0,   'Lei 7.326/2023', 'Lei 7.326/2023', '2024-01-21'),
    ('ES', 17.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('GO', 19.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('MA', 23.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('MT', 17.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('MS', 17.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('MG', 18.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('PA', 19.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('PB', 20.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('PR', 19.5, 0,   'Lei 21.850/2023', 'Lei 21.850/2023', '2024-03-18'),
    ('PE', 20.5, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('PI', 22.5, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('RJ', 22.0, 2.0, 'Lei 10.253/2023 + LC 210/2023 (20% ICMS + 2% FECP = 22%)', 'Lei 10.253/2023 + LC 210/2023', '2024-03-20'),
    ('RN', 20.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('RO', 19.5, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('RR', 20.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('RS', 17.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('SC', 17.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('SP', 18.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('SE', 20.0, 1.0, 'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01'),
    ('TO', 20.0, 0,   'legislação estadual vigente, sem alteração recente conhecida', 'legislação estadual vigente', '2000-01-01')
ON CONFLICT (uf, vigente_desde) DO NOTHING;
