
-- users table ---------------------------------------------------
CREATE TABLE "users" (
  "user_id" uuid PRIMARY KEY,
  "name" text NOT NULL,
  "email" text UNIQUE NOT NULL,
  "roles" text[] NOT NULL,
  "password_hash" text NOT NULL,
  "enabled" boolean NOT NULL,
  "email_confirm_token" text UNIQUE,
  "date_created" timestamp NOT NULL,
  "date_updated" timestamp NOT NULL
);

-- products table ---------------------------------------------------
CREATE TABLE "products" (
  "product_id" uuid PRIMARY KEY,
  "name" text NOT NULL,
  "desc" text,
  "imageURL" text,
  "price" bigint NOT NULL,
  "quantity" int NOT NULL,
  "date_created" timestamp NOT NULL,
  "date_updated" timestamp NOT NULL
);

-- orders table ---------------------------------------------------
CREATE TABLE "orders" (
  "order_id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "amount" bigint NOT NULL,
  "status" text NOT NULL,
  "date_created" timestamp NOT NULL
);

ALTER TABLE "orders" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

-- order_items table ---------------------------------------------------
CREATE TABLE "order_items" (
  "order_item_id" uuid PRIMARY KEY,
  "order_id" uuid NOT NULL,
  "product_id" uuid NOT NULL,
  "price" bigint NOT NULL,
  "quantity" int NOT NULL,
  "date_created" timestamp NOT NULL
);

CREATE INDEX "product_id_index" ON "order_items" ("product_id");

CREATE INDEX "order_id_index" ON "order_items" ("order_id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("order_id");

ALTER TABLE "order_items" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("product_id");

-- payments table ---------------------------------------------------
CREATE TABLE "payments" (
  "payment_id" uuid PRIMARY KEY,
  "order_id" uuid NOT NULL,
  "partner" text NOT NULL,
  "partner_order_id" text NOT NULL,
  "partner_transaction_id" text UNIQUE,
  "status" text NOT NULL,
  "currency" text NOT NULL,
  "date_created" timestamp NOT NULL
);

CREATE INDEX "order_id_index" ON "payments" ("order_id");

CREATE INDEX "partner_order_id_index" ON "payments" ("partner_order_id");

CREATE INDEX "partner_transaction_id_index" ON "payments" ("partner_transaction_id");

CREATE INDEX "order_id_status_index" ON "payments" ("order_id", "status");

ALTER TABLE "payments" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("order_id");
