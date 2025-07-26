-- Drop triggers
DROP TRIGGER IF EXISTS update_invoices_updated_at ON invoices;
DROP TRIGGER IF EXISTS update_clients_updated_at ON clients;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_invoices_deleted_at;
DROP INDEX IF EXISTS idx_invoices_number;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_client_id;
DROP INDEX IF EXISTS idx_clients_deleted_at;
DROP INDEX IF EXISTS idx_clients_email;

-- Drop tables (in reverse order due to foreign key constraints)
DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS clients;