-- Boucherie API — Full schema

CREATE TABLE IF NOT EXISTS clients (
    id         TEXT PRIMARY KEY,
    name       TEXT    NOT NULL,
    phone      TEXT    NOT NULL,
    email      TEXT    DEFAULT '',
    avatar     TEXT    DEFAULT '',
    total_credit REAL  DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    category     TEXT NOT NULL CHECK(category IN ('boeuf','agneau','poulet','veau','charcuterie')),
    price_per_kg REAL NOT NULL CHECK(price_per_kg > 0),
    image        TEXT DEFAULT '',
    in_stock     INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS sales (
    id            TEXT PRIMARY KEY,
    client_id     TEXT NOT NULL REFERENCES clients(id),
    client_name   TEXT NOT NULL,
    total         REAL NOT NULL,
    paid_amount   REAL NOT NULL DEFAULT 0,
    credit_amount REAL NOT NULL DEFAULT 0,
    date          DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sale_items (
    id           TEXT PRIMARY KEY,
    sale_id      TEXT NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    product_id   TEXT NOT NULL,
    product_name TEXT NOT NULL,
    quantity     REAL NOT NULL CHECK(quantity > 0),
    subtotal     REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS credits (
    id               TEXT PRIMARY KEY,
    client_id        TEXT NOT NULL REFERENCES clients(id),
    client_name      TEXT NOT NULL,
    sale_id          TEXT NOT NULL REFERENCES sales(id),
    amount           REAL NOT NULL,
    remaining_amount REAL NOT NULL,
    status           TEXT NOT NULL DEFAULT 'en_cours' CHECK(status IN ('en_cours','en_retard','paye')),
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    due_date         DATETIME
);

CREATE TABLE IF NOT EXISTS payments (
    id        TEXT PRIMARY KEY,
    credit_id TEXT NOT NULL REFERENCES credits(id) ON DELETE CASCADE,
    amount    REAL NOT NULL CHECK(amount > 0),
    date      DATETIME DEFAULT CURRENT_TIMESTAMP,
    method    TEXT NOT NULL DEFAULT 'cash' CHECK(method IN ('cash','carte','virement'))
);

CREATE TABLE IF NOT EXISTS orders (
    id           TEXT PRIMARY KEY,
    client_id    TEXT NOT NULL REFERENCES clients(id),
    client_name  TEXT NOT NULL,
    client_phone TEXT NOT NULL,
    pickup_date  DATETIME NOT NULL,
    notes        TEXT DEFAULT '',
    status       TEXT NOT NULL DEFAULT 'en_attente' CHECK(status IN ('en_attente','confirmee','prete','livree','annulee')),
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id           TEXT PRIMARY KEY,
    order_id     TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id   TEXT NOT NULL,
    product_name TEXT NOT NULL,
    quantity     REAL NOT NULL CHECK(quantity > 0)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_sales_client    ON sales(client_id);
CREATE INDEX IF NOT EXISTS idx_sales_date      ON sales(date);
CREATE INDEX IF NOT EXISTS idx_credits_client  ON credits(client_id);
CREATE INDEX IF NOT EXISTS idx_credits_status  ON credits(status);
CREATE INDEX IF NOT EXISTS idx_orders_status   ON orders(status);
CREATE INDEX IF NOT EXISTS idx_payments_credit ON payments(credit_id);
