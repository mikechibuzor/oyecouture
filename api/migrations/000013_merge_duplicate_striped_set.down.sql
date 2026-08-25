-- Split the duplicate back out.
UPDATE products
SET active = TRUE
WHERE slug = 'ivory-black-set';

UPDATE product_images
SET
    product_id = (SELECT id FROM products WHERE slug = 'ivory-black-set'),
    position = 1
WHERE r2_key = 'oyecouture/casual-08';
