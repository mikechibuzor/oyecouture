-- The striped linen shirt was seeded as two products. When the client's
-- photography was first grouped, one frame of the shoot read as a plain shirt
-- at thumbnail size and was split off as "The Ivory & Black Set"; it is the
-- same piece as "The Striped Linen Set", shot from a different angle.
--
-- The stray frame moves onto the surviving product and the duplicate is
-- deactivated rather than deleted, so the change is reversible and any future
-- order history keeps its product link.
UPDATE product_images
SET
    product_id = (SELECT id FROM products WHERE slug = 'striped-linen-set'),
    -- The subquery reads the pre-statement snapshot, so this appends after the
    -- images the surviving product already has.
    position = COALESCE(
        (
            SELECT MAX(existing.position)
            FROM product_images existing
            JOIN products owner ON owner.id = existing.product_id
            WHERE owner.slug = 'striped-linen-set'
        ),
        0
    ) + 1
WHERE r2_key = 'oyecouture/casual-08';

UPDATE products
SET active = FALSE
WHERE slug = 'ivory-black-set';
