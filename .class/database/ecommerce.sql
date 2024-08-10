-- users table ---------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
   user_id               UUID        NOT NULL,
   name                  TEXT        NOT NULL,
   email                 TEXT UNIQUE NOT NULL,
   roles                 TEXT[]      NOT NULL,
   password_hash         TEXT        NOT NULL,
   enabled               BOOLEAN     NOT NULL,
   email_confirm_token   TEXT UNIQUE     NULL,
   date_created          TIMESTAMP   NOT NULL,
   date_updated          TIMESTAMP   NOT NULL,

   PRIMARY KEY (user_id)
);

-- products table ---------------------------------------------------
CREATE TABLE IF NOT EXISTS products (
  product_id    UUID        NOT NULL,
  name          TEXT        NOT NULL,
  description   TEXT            NULL,
  imageURL      TEXT            NULL,
  price         BIGINT      NOT NULL,
  quantity      INT         NOT NULL,
  date_created  TIMESTAMP   NOT NULL,
  date_updated  TIMESTAMP   NOT NULL,

  PRIMARY KEY (product_id)
);

-- orders table ---------------------------------------------------
CREATE TABLE IF NOT EXISTS orders (
  order_id      UUID        NOT NULL,
  user_id       UUID        NOT NULL,
  amount        BIGINT      NOT NULL,
  status        TEXT        NOT NULL,
  date_created  TIMESTAMP   NOT NULL,

  PRIMARY KEY (order_id)
);

ALTER TABLE orders ADD FOREIGN KEY (user_id) REFERENCES users (user_id);

-- order_items table ---------------------------------------------------
CREATE TABLE IF NOT EXISTS order_items (
  order_item_id     UUID        NOT NULL,
  order_id          UUID        NOT NULL,
  product_id        UUID        NOT NULL,
  price             BIGINT      NOT NULL,
  quantity          INT         NOT NULL,
  date_created      TIMESTAMP   NOT NULL,

  PRIMARY KEY (order_item_id)
);

CREATE INDEX product_id_index ON order_items (product_id);

CREATE INDEX order_id_index ON order_items (order_id);

ALTER TABLE order_items ADD FOREIGN KEY (order_id) REFERENCES orders (order_id);

ALTER TABLE order_items ADD FOREIGN KEY (product_id) REFERENCES products (product_id);

-- payments table ---------------------------------------------------
CREATE TABLE IF NOT EXISTS payments (
  payment_id                UUID        NOT NULL,
  order_id                  UUID        NOT NULL,
  partner                   TEXT        NOT NULL,
  partner_order_id          TEXT        NOT NULL,
  partner_transaction_id    TEXT UNIQUE     NULL,
  status                    TEXT        NOT NULL,
  currency                  TEXT        NOT NULL,
  date_created              TIMESTAMP   NOT NULL,

  PRIMARY KEY (payment_id)
);

CREATE INDEX order_id_index ON payments (order_id);

CREATE INDEX partner_order_id_index ON payments (partner_order_id);

CREATE INDEX partner_transaction_id_index ON payments (partner_transaction_id);

CREATE INDEX order_id_status_index ON payments (order_id, status);

ALTER TABLE payments ADD FOREIGN KEY (order_id) REFERENCES orders (order_id);
