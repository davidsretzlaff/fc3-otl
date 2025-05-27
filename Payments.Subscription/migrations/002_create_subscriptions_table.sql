-- Criação da tabela subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    plan_id VARCHAR(36) NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    status ENUM('pending', 'active', 'inactive', 'cancelled', 'suspended') NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_subscriptions_plan_id (plan_id),
    INDEX idx_subscriptions_customer_id (customer_id),
    INDEX idx_subscriptions_status (status),
    INDEX idx_subscriptions_created_at (created_at),
    INDEX idx_subscriptions_customer_status (customer_id, status)
); 