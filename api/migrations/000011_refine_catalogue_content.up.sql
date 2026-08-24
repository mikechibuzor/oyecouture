UPDATE products
SET name = 'Edo Traditional Attire'
WHERE slug = 'ivory-buba-set';

UPDATE products
SET active = FALSE
WHERE slug = 'wedding-ensemble';

UPDATE products
SET
    name = 'The Blue Linen Set',
    description = 'Deep blue linen with a camp collar and shell buttons, cut generously through the body.'
WHERE slug = 'indigo-linen-set';

UPDATE products
SET
    name = 'The Ambassador Two-Piece',
    category = 'Two-Piece'
WHERE slug = 'ambassador';

UPDATE products
SET
    name = 'The Chairman Two-Piece',
    category = 'Two-Piece'
WHERE slug = 'chairman';
