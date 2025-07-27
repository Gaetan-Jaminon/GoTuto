-- Drop constraints first
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS check_positive_amount;
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS check_invoice_status;

-- Drop indexes
DROP INDEX IF EXISTS idx_invoices_deleted_at;
DROP INDEX IF EXISTS idx_invoices_due_date;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_invoice_number;
DROP INDEX IF EXISTS idx_invoices_client_id;

-- Drop tables (invoices first due to foreign key)
DROP TABLE IF EXISTS invoices;

-- Drop client indexes
DROP INDEX IF EXISTS idx_clients_deleted_at;
DROP INDEX IF EXISTS idx_clients_email;

-- Drop clients table
DROP TABLE IF EXISTS clients;