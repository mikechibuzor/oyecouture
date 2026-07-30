-- PLACEHOLDER FEES — the client has not supplied delivery pricing yet.
-- These exist so checkout is functional end to end; they are stored in the
-- database (not code) precisely so they can be corrected without a deploy.
-- International is zero because that delivery is quoted on WhatsApp.
INSERT INTO delivery_zones (name, fee_kobo, active)
VALUES
    ('Lagos', 500000, TRUE),
    ('Rest of Nigeria', 1000000, TRUE),
    ('International', 0, TRUE)
ON CONFLICT (name) DO NOTHING;
