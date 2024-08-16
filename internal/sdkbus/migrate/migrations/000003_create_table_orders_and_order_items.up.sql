-- orders table ---------------------------------------------------
CREATE TABLE IF NOT EXISTS orders (
    order_id      UUID        NOT NULL,
    user_id       UUID        NOT NULL,
    amount        BIGINT      NOT NULL,
    status        TEXT        NOT NULL,
    date_created  TIMESTAMP   NOT NULL,

    PRIMARY KEY (order_id)
);

ALTER TABLE orders ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (user_id);

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

ALTER TABLE order_items ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders (order_id);

ALTER TABLE order_items ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products (product_id);