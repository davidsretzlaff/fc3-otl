-- Criação da tabela customers
CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    document VARCHAR(50) NOT NULL UNIQUE,
    status ENUM('active', 'inactive', 'blocked') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_customers_email (email),
    INDEX idx_customers_document (document),
    INDEX idx_customers_status (status),
    INDEX idx_customers_created_at (created_at)
); 