CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    document_number VARCHAR(20) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    operation_type_id INTEGER NOT NULL, 
    amount NUMERIC(15, 2) NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_account 
        FOREIGN KEY(account_id) 
        REFERENCES accounts(account_id) 
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_transactions_account_id ON transactions(account_id);