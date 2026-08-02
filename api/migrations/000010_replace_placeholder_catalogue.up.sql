-- The client's real photographed inventory replaces the design-mockup
-- placeholders. Placeholders are deactivated rather than deleted so order
-- history keeps its product links and the change stays reversible.
UPDATE products SET active = FALSE
WHERE section IN ('trad', 'casuals-men', 'casuals-women');

INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('gold-agbada', 'Gold Aso-Oke Agbada', 'trad', 'Agbada', 78000000, 'Hand-woven aso-oke in gold and amber, cut as a full agbada over a contrast tunic and sokoto. Finished with a matching fila.', '', FALSE, TRUE, 1);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('crimson-agbada', 'Crimson Ceremonial Agbada', 'trad', 'Agbada', 78000000, 'A deep crimson agbada with an embroidered placket, cut wide through the sleeve for presence. Worn for the occasions that are remembered.', 'New', FALSE, TRUE, 2);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('burgundy-agbada', 'Burgundy Aso-Oke Agbada', 'trad', 'Agbada', 78000000, 'Striped aso-oke in burgundy and slate, tailored as a flowing agbada with hand-finished neckline embroidery.', '', FALSE, TRUE, 3);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('sand-agbada', 'Sand Aso-Oke Agbada', 'trad', 'Agbada', 78000000, 'Woven in a soft sand tone with tonal embroidery — an agbada for daylight ceremonies and quieter formality.', '', FALSE, TRUE, 4);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('red-asooke-set', 'Red Aso-Oke Two-Piece', 'trad', 'Aso-Oke', 52000000, 'Buba and sokoto in textured red aso-oke with a contrast indigo yoke. The agbada left off, the craft left visible.', '', FALSE, TRUE, 5);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('purple-asooke-set', 'Purple Aso-Oke Two-Piece', 'trad', 'Aso-Oke', 52000000, 'A royal purple buba and sokoto with tonal embroidery down the placket, cut close and worn with a matching fila.', '', FALSE, TRUE, 6);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('burgundy-asooke-set', 'Burgundy Aso-Oke Two-Piece', 'trad', 'Aso-Oke', 52000000, 'Burgundy aso-oke with a fine woven stripe, tailored as a modern two-piece for traditional occasions.', 'New', FALSE, TRUE, 7);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('camel-asooke-set', 'Camel Aso-Oke Two-Piece', 'trad', 'Aso-Oke', 52000000, 'Camel-toned aso-oke with hand-worked chest embroidery — the most versatile piece in the collection.', '', FALSE, TRUE, 8);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('white-buba-set', 'White Embroidered Buba & Sokoto', 'trad', 'Kaftans', 29500000, 'A crisp white buba with dark medallion embroidery, worn over a contrast sokoto. Light enough for the whole day.', '', FALSE, TRUE, 9);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('ivory-buba-set', 'Ivory Embroidered Buba & Sokoto', 'trad', 'Kaftans', 29500000, 'Ivory buba with tonal hand embroidery and a relaxed sleeve, finished with a woven sokoto.', '', FALSE, TRUE, 10);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('wedding-ensemble', 'Aso-Oke Wedding Ensemble', 'trad', 'Aso-Oke', 86000000, 'A coordinated wedding ensemble made to a single measure across both outfits, woven and finished together so the pair reads as one.', 'Bespoke', FALSE, TRUE, 11);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('adire-linen-set', 'The Adire Linen Set', 'casuals-men', 'Shirts', 11500000, 'A mandarin-collar linen shirt worn over hand-dyed adire trousers — the collection''s most personal pairing.', 'New', TRUE, TRUE, 12);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('sky-linen-set', 'The Sky Linen Set', 'casuals-men', 'Shirts', 11500000, 'A pale blue linen shirt with a clean mandarin collar, cut long and worn over tailored navy trousers.', '', TRUE, TRUE, 13);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('sage-linen-set', 'The Sage Linen Set', 'casuals-men', 'Shirts', 11500000, 'Sage linen with a soft placket and an unstructured shoulder, matched to olive trousers that move with it.', '', TRUE, TRUE, 14);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('blush-linen-set', 'The Blush Linen Set', 'casuals-men', 'Shirts', 11500000, 'A blush two-piece in washed linen with a covered button placket — restrained, and quietly the boldest thing in the room.', 'New', TRUE, TRUE, 15);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('striped-linen-set', 'The Striped Linen Set', 'casuals-men', 'Shirts', 9800000, 'A short-sleeve shirt in fine woven stripe, cut square and worn with a straight black trouser.', '', TRUE, TRUE, 16);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('ivory-black-set', 'The Ivory & Black Set', 'casuals-men', 'Shirts', 9800000, 'An ivory short-sleeve shirt with a relaxed body, paired with a wide black linen trouser.', '', TRUE, TRUE, 17);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('onyx-linen-set', 'The Onyx Linen Set', 'casuals-men', 'Shirts', 9800000, 'All black, all linen — a short-sleeve shirt and trouser cut for evenings that start warm and end late.', '', TRUE, TRUE, 18);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('ivory-linen-set', 'The Ivory Linen Set', 'casuals-men', 'Shirts', 9800000, 'A full ivory two-piece in breathable linen, the short sleeve finished with a clean turned hem.', '', TRUE, TRUE, 19);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('indigo-linen-set', 'The Indigo Linen Set', 'casuals-men', 'Shirts', 9800000, 'Deep indigo linen with a camp collar and shell buttons, cut generously through the body.', 'New', TRUE, TRUE, 20);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('cream-linen-set', 'The Cream Linen Set', 'casuals-women', 'Sets', 10500000, 'A cream linen shirt and trouser cut for ease — the shirt long enough to wear open, the trouser tailored straight.', '', TRUE, TRUE, 21);
INSERT INTO products (slug, name, section, category, price_kobo, description, tag, has_sizes, active, sort_order)
VALUES ('botanical-short-set', 'The Botanical Short Set', 'casuals-women', 'Sets', 10500000, 'A botanical-print short-sleeve shirt over tailored red linen shorts. Made for weekends that photograph well.', 'New', TRUE, TRUE, 22);

INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-01-1', 1 FROM products WHERE slug = 'gold-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-01-2', 2 FROM products WHERE slug = 'gold-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-01-3', 3 FROM products WHERE slug = 'gold-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-01-4', 4 FROM products WHERE slug = 'gold-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-01-5', 5 FROM products WHERE slug = 'gold-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-02-1', 1 FROM products WHERE slug = 'crimson-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-02-2', 2 FROM products WHERE slug = 'crimson-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-02-3', 3 FROM products WHERE slug = 'crimson-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-02-4', 4 FROM products WHERE slug = 'crimson-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-02-5', 5 FROM products WHERE slug = 'crimson-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-04-1', 1 FROM products WHERE slug = 'burgundy-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-04-2', 2 FROM products WHERE slug = 'burgundy-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-04-3', 3 FROM products WHERE slug = 'burgundy-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-04-4', 4 FROM products WHERE slug = 'burgundy-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-04-5', 5 FROM products WHERE slug = 'burgundy-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-04-6', 6 FROM products WHERE slug = 'burgundy-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-05-1', 1 FROM products WHERE slug = 'sand-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-05-2', 2 FROM products WHERE slug = 'sand-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-05-3', 3 FROM products WHERE slug = 'sand-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-05-4', 4 FROM products WHERE slug = 'sand-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-05-5', 5 FROM products WHERE slug = 'sand-agbada';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-08-1', 1 FROM products WHERE slug = 'red-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-08-2', 2 FROM products WHERE slug = 'red-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-08-3', 3 FROM products WHERE slug = 'red-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-08-4', 4 FROM products WHERE slug = 'red-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-09-1', 1 FROM products WHERE slug = 'purple-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-09-2', 2 FROM products WHERE slug = 'purple-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-09-3', 3 FROM products WHERE slug = 'purple-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-11-1', 1 FROM products WHERE slug = 'burgundy-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-11-2', 2 FROM products WHERE slug = 'burgundy-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-11-3', 3 FROM products WHERE slug = 'burgundy-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-11-4', 4 FROM products WHERE slug = 'burgundy-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-11-5', 5 FROM products WHERE slug = 'burgundy-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-12-1', 1 FROM products WHERE slug = 'camel-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-12-2', 2 FROM products WHERE slug = 'camel-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-12-3', 3 FROM products WHERE slug = 'camel-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-12-4', 4 FROM products WHERE slug = 'camel-asooke-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-03-1', 1 FROM products WHERE slug = 'white-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-03-2', 2 FROM products WHERE slug = 'white-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-03-3', 3 FROM products WHERE slug = 'white-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-03-4', 4 FROM products WHERE slug = 'white-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-06-1', 1 FROM products WHERE slug = 'ivory-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-06-2', 2 FROM products WHERE slug = 'ivory-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-06-3', 3 FROM products WHERE slug = 'ivory-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-06-4', 4 FROM products WHERE slug = 'ivory-buba-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-07-1', 1 FROM products WHERE slug = 'wedding-ensemble';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-07-2', 2 FROM products WHERE slug = 'wedding-ensemble';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-07-3', 3 FROM products WHERE slug = 'wedding-ensemble';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/asoke-look-07-4', 4 FROM products WHERE slug = 'wedding-ensemble';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-01', 1 FROM products WHERE slug = 'adire-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-02', 2 FROM products WHERE slug = 'adire-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-03', 3 FROM products WHERE slug = 'adire-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-10', 1 FROM products WHERE slug = 'sky-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-21', 2 FROM products WHERE slug = 'sky-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-22', 3 FROM products WHERE slug = 'sky-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-14', 1 FROM products WHERE slug = 'sage-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-18', 2 FROM products WHERE slug = 'sage-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-19', 3 FROM products WHERE slug = 'sage-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-20', 4 FROM products WHERE slug = 'sage-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-27', 1 FROM products WHERE slug = 'blush-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-28', 2 FROM products WHERE slug = 'blush-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-29', 3 FROM products WHERE slug = 'blush-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-30', 4 FROM products WHERE slug = 'blush-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-07', 1 FROM products WHERE slug = 'striped-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-09', 2 FROM products WHERE slug = 'striped-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-08', 1 FROM products WHERE slug = 'ivory-black-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-11', 1 FROM products WHERE slug = 'onyx-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-12', 2 FROM products WHERE slug = 'onyx-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-13', 3 FROM products WHERE slug = 'onyx-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-23', 1 FROM products WHERE slug = 'ivory-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-24', 2 FROM products WHERE slug = 'ivory-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-25', 3 FROM products WHERE slug = 'ivory-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-26', 4 FROM products WHERE slug = 'ivory-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-31', 1 FROM products WHERE slug = 'indigo-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-32', 2 FROM products WHERE slug = 'indigo-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-33', 3 FROM products WHERE slug = 'indigo-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-04', 1 FROM products WHERE slug = 'cream-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-05', 2 FROM products WHERE slug = 'cream-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-06', 3 FROM products WHERE slug = 'cream-linen-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-15', 1 FROM products WHERE slug = 'botanical-short-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-16', 2 FROM products WHERE slug = 'botanical-short-set';
INSERT INTO product_images (product_id, r2_key, position)
SELECT id, 'oyecouture/casual-17', 3 FROM products WHERE slug = 'botanical-short-set';
