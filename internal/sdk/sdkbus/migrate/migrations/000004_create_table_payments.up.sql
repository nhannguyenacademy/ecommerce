CREATE TABLE IF NOT EXISTS payments (
    payment_id                UUID        NOT NULL,
    order_id                  UUID        NOT NULL,
    partner                   TEXT        NOT NULL,
    partner_order_id          TEXT        NOT NULL,
    partner_transaction_id    TEXT UNIQUE     NULL,
    status                    TEXT        NOT NULL,
    currency                  TEXT        NOT NULL,
    date_created              TIMESTAMP   NOT NULL,
    date_updated              TIMESTAMP   NOT NULL,

    PRIMARY KEY (payment_id)
);

CREATE INDEX order_id_index ON payments (order_id);

CREATE INDEX partner_order_id_index ON payments (partner_order_id);

CREATE INDEX partner_transaction_id_index ON payments (partner_transaction_id);

CREATE INDEX order_id_status_index ON payments (order_id, status);

ALTER TABLE payments ADD CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders (order_id);