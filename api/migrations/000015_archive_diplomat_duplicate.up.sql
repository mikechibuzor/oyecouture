-- The Diplomat and Chairman entries use the same product photography.
-- Keep Chairman as the client-approved catalogue entry and archive Diplomat.
-- Archiving preserves existing order references and keeps this reversible.
UPDATE products
SET active = FALSE
WHERE slug = 'diplomat';
