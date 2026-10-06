-- The addresses and the sender name become editable from the admin, so changing who hears
-- from a customer no longer needs a deploy.
--
-- contact_email is published on the contact page. It is seeded with the real address because
-- the page previously hardcoded hello@oyecouture.com, a domain this business does not run.
--
-- notify_email is seeded empty on purpose: empty means "use OWNER_NOTIFY_EMAIL", so enquiries
-- keep arriving exactly where they do now until someone chooses otherwise.
--
-- from_name is seeded with the name already configured in EMAIL_FROM, so the admin shows what
-- recipients actually see rather than an empty box.
INSERT INTO settings (key, value) VALUES
    ('contact_email', 'hello@oyecouturebespoke.com'),
    ('notify_email', ''),
    ('from_name', 'Oyecouture Bespoke')
ON CONFLICT (key) DO NOTHING;
