DROP INDEX  IF EXISTS order_id_index;

DROP INDEX  IF EXISTS partner_order_id_index;

DROP INDEX  IF EXISTS partner_transaction_id_index;

DROP INDEX  IF EXISTS order_id_status_index;

ALTER TABLE payments DROP CONSTRAINT fk_order_id;

DROP TABLE IF EXISTS payments;