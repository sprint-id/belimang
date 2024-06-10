DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_is_admin;
DROP INDEX IF EXISTS idx_merchants_id;
DROP INDEX IF EXISTS idx_merchants_merchant_category;
DROP INDEX IF EXISTS idx_merchants_name;
DROP INDEX IF EXISTS idx_merchants_location;
DROP INDEX IF EXISTS idx_items_id;
DROP INDEX IF EXISTS idx_items_product_category;
DROP INDEX IF EXISTS idx_estimates_id;
DROP INDEX IF EXISTS idx_estimates_user_id;

-- Drop extensions
DROP EXTENSION IF EXISTS "btree_gist";

DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS estimate_order_items;

DROP TABLE IF EXISTS estimate_orders;

DROP TABLE IF EXISTS estimate_users_locations;

DROP TABLE IF EXISTS estimates;

DROP TABLE IF EXISTS items;

DROP TABLE IF EXISTS merchants;

DROP TABLE IF EXISTS users;