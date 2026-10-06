-- The shop can be closed to visitors from the admin, without a deploy. Seeded off, with
-- wording ready so turning it on never draws an empty page.
INSERT INTO settings (key, value) VALUES
    ('site_held', 'false'),
    ('hold_headline', 'We are preparing something new.'),
    ('hold_message', 'Oyecouture is closed to visitors for a short while. We are a message away in the meantime.')
ON CONFLICT (key) DO NOTHING;
