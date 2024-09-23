CREATE TABLE IF NOT EXISTS products (
    product_id    UUID        NOT NULL,
    name          TEXT        NOT NULL,
    description   TEXT            NULL,
    image_url      TEXT            NULL,
    price         BIGINT      NOT NULL,
    quantity      INT         NOT NULL,
    date_created  TIMESTAMP   NOT NULL,
    date_updated  TIMESTAMP   NOT NULL,

    PRIMARY KEY (product_id)
);