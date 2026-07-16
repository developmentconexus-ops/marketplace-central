ALTER TABLE listings DROP CONSTRAINT listings_status_check;

ALTER TABLE listings ADD CONSTRAINT listings_status_check
    CHECK (status IN ('active','paused','closed','unknown','under_review','inactive','payment_required','not_yet_active'));
