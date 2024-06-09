BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR UNIQUE,
    email VARCHAR,
    password VARCHAR,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS merchants (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    merchant_category VARCHAR NOT NULL,
    image_url VARCHAR NOT NULL,
    location_lat FLOAT NOT NULL,
    location_long FLOAT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    merchant_id UUID REFERENCES merchants(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    product_category VARCHAR NOT NULL,
    price INT NOT NULL,
    image_url VARCHAR NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS estimates (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    total_price INT NOT NULL,
    delivery_time INT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS estimate_users_locations (
    id UUID PRIMARY KEY,
    estimate_id UUID REFERENCES estimates(id) ON DELETE CASCADE,
    location_lat FLOAT NOT NULL,
    location_long FLOAT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS estimate_orders (
    id UUID PRIMARY KEY,
    estimate_id UUID REFERENCES estimates(id) ON DELETE CASCADE,
    merchant_id UUID REFERENCES merchants(id) ON DELETE CASCADE,
    is_starting_point BOOLEAN DEFAULT FALSE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS estimate_order_items (
    id UUID PRIMARY KEY,
    estimate_order_id UUID REFERENCES estimate_orders(id) ON DELETE CASCADE,
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    estimate_id UUID REFERENCES estimates(id) ON DELETE CASCADE,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "btree_gist";

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_id ON users (id);
CREATE INDEX IF NOT EXISTS idx_merchants_id ON merchants (id);
CREATE INDEX IF NOT EXISTS idx_merchants_user_id ON merchants (user_id);
CREATE INDEX IF NOT EXISTS idx_merchants_merchant_category ON merchants (merchant_category);
CREATE INDEX IF NOT EXISTS idx_merchants_name ON merchants (name);
-- Creating an index on the merchants table for geolocation searches
CREATE INDEX IF NOT EXISTS idx_merchants_location ON merchants USING gist (
    point(location_lat, location_long)
);

CREATE INDEX IF NOT EXISTS idx_items_id ON items (id);
CREATE INDEX IF NOT EXISTS idx_items_user_id ON items (user_id);
CREATE INDEX IF NOT EXISTS idx_items_product_category ON items (product_category);
CREATE INDEX IF NOT EXISTS idx_estimates_id ON estimates (id);
CREATE INDEX IF NOT EXISTS idx_estimates_user_id ON estimates (user_id);



COMMIT TRANSACTION;