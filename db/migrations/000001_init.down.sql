DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_merchants_id;
DROP INDEX IF EXISTS idx_merchants_user_id;
DROP INDEX IF EXISTS idx_merchants_merchant_category;
DROP INDEX IF EXISTS idx_items_id;
DROP INDEX IF EXISTS idx_items_user_id;
DROP INDEX IF EXISTS idx_items_product_category;

DROP TABLE IF EXISTS items;

DROP TABLE IF EXISTS merchants;

DROP TABLE IF EXISTS users;