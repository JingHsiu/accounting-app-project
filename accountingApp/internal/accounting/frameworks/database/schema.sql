-- Create wallets table
CREATE TABLE IF NOT EXISTS wallets (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(15) NOT NULL CHECK (type IN ('CASH', 'BANK', 'CREDIT', 'INVESTMENT')),
    currency CHAR(3) NOT NULL,
    balance_amount BIGINT NOT NULL DEFAULT 0,
    balance_currency CHAR(3) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_wallet_currency CHECK (currency = balance_currency)
);

-- Create expense_categories table
CREATE TABLE IF NOT EXISTS expense_categories (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, name)
);

-- Create expense_subcategories table  
CREATE TABLE IF NOT EXISTS expense_subcategories (
    id VARCHAR(36) PRIMARY KEY,
    parent_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    
    FOREIGN KEY (parent_id) REFERENCES expense_categories(id) ON DELETE CASCADE,
    UNIQUE(parent_id, name)
);

-- Create income_categories table
CREATE TABLE IF NOT EXISTS income_categories (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, name)
);

-- Create income_subcategories table
CREATE TABLE IF NOT EXISTS income_subcategories (
    id VARCHAR(36) PRIMARY KEY,
    parent_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    
    FOREIGN KEY (parent_id) REFERENCES income_categories(id) ON DELETE CASCADE,
    UNIQUE(parent_id, name)
);

-- Create expense_records table
CREATE TABLE IF NOT EXISTS expense_records (
    id VARCHAR(36) PRIMARY KEY,
    wallet_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES expense_categories(id)
);

-- Create income_records table
CREATE TABLE IF NOT EXISTS income_records (
    id VARCHAR(36) PRIMARY KEY,
    wallet_id VARCHAR(36) NOT NULL,
    category_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES income_categories(id)
);

-- Create transfers table
CREATE TABLE IF NOT EXISTS transfers (
    id VARCHAR(36) PRIMARY KEY,
    from_wallet_id VARCHAR(36) NOT NULL,
    to_wallet_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL,
    fee_amount BIGINT NOT NULL DEFAULT 0 CHECK (fee_amount >= 0),
    fee_currency CHAR(3) NOT NULL,
    description TEXT,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    FOREIGN KEY (from_wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (to_wallet_id) REFERENCES wallets(id),
    CHECK (from_wallet_id != to_wallet_id),
    CHECK (currency = fee_currency)
);

-- Create indexes for better query performance
CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_expense_categories_user_id ON expense_categories(user_id);
CREATE INDEX idx_income_categories_user_id ON income_categories(user_id);
CREATE INDEX idx_expense_records_wallet_id ON expense_records(wallet_id);
CREATE INDEX idx_expense_records_date ON expense_records(date);
CREATE INDEX idx_income_records_wallet_id ON income_records(wallet_id);
CREATE INDEX idx_income_records_date ON income_records(date);
CREATE INDEX idx_transfers_from_wallet ON transfers(from_wallet_id);
CREATE INDEX idx_transfers_to_wallet ON transfers(to_wallet_id);
CREATE INDEX idx_transfers_date ON transfers(date);