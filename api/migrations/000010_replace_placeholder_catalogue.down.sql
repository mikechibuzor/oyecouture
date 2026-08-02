-- Restore the placeholder catalogue.
DELETE FROM products WHERE slug IN ('gold-agbada', 'crimson-agbada', 'burgundy-agbada', 'sand-agbada', 'red-asooke-set', 'purple-asooke-set', 'burgundy-asooke-set', 'camel-asooke-set', 'white-buba-set', 'ivory-buba-set', 'wedding-ensemble', 'adire-linen-set', 'sky-linen-set', 'sage-linen-set', 'blush-linen-set', 'striped-linen-set', 'ivory-black-set', 'onyx-linen-set', 'ivory-linen-set', 'indigo-linen-set', 'cream-linen-set', 'botanical-short-set');
UPDATE products SET active = TRUE WHERE section IN ('trad', 'casuals-men', 'casuals-women');
