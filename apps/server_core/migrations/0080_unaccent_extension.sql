-- 0080: unaccent extension for accent-insensitive catalog search.
--
-- The ERP descriptions are stored exactly as the source writes them, accents and
-- all ("VÁLVULA DE MICTÓRIO"). The operator types what is on the keyboard
-- ("valvula"), so a plain ILIKE answers "Nenhum produto encontrado." for a
-- product that exists. unaccent() folds both sides of the comparison at query
-- time; the stored text is never rewritten, so no data is lost or normalized.
CREATE EXTENSION IF NOT EXISTS unaccent;
