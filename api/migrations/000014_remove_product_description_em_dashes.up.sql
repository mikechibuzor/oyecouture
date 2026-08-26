UPDATE products
SET description = REPLACE(description, ' ' || CHR(8212) || ' ', ', ')
WHERE STRPOS(description, CHR(8212)) > 0;
