-- Ensure we're in the catalog schema
SET search_path TO catalog;

-- Drop trigger and function
DROP TRIGGER IF EXISTS ensure_single_primary_category ON product_categories;
DROP FUNCTION IF EXISTS check_single_primary_category();

-- Drop constraints
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_product_status;
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_product_price;

-- Drop indexes
DROP INDEX IF EXISTS idx_product_variants_deleted_at;
DROP INDEX IF EXISTS idx_product_variants_sku;
DROP INDEX IF EXISTS idx_product_variants_product_id;
DROP INDEX IF EXISTS idx_product_categories_category_id;
DROP INDEX IF EXISTS idx_products_deleted_at;
DROP INDEX IF EXISTS idx_products_status;
DROP INDEX IF EXISTS idx_products_brand_id;
DROP INDEX IF EXISTS idx_products_sku;
DROP INDEX IF EXISTS idx_brands_deleted_at;
DROP INDEX IF EXISTS idx_categories_deleted_at;
DROP INDEX IF EXISTS idx_categories_is_active;
DROP INDEX IF EXISTS idx_categories_parent_id;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS categories;