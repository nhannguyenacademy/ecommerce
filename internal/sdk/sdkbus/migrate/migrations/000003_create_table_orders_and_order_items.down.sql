-- order_items table ---------------------------------------------

DROP INDEX  IF EXISTS product_id_index;

DROP INDEX  IF EXISTS order_id_index;

ALTER TABLE order_items DROP CONSTRAINT fk_order_id;

ALTER TABLE order_items DROP CONSTRAINT fk_product_id;

DROP TABLE IF EXISTS order_items;

-- orders table ---------------------------------------------------

ALTER TABLE orders DROP CONSTRAINT fk_user_id;

DROP TABLE IF EXISTS orders;