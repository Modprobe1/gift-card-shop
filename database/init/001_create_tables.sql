-- Создание таблиц для криптообменника

-- Таблица валют
CREATE TABLE currencies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10),
    network VARCHAR(50),
    icon_url VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    min_amount DECIMAL(20, 8) DEFAULT 0,
    max_amount DECIMAL(20, 8),
    decimals INT DEFAULT 8,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Таблица курсов валют
CREATE TABLE exchange_rates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_currency_id INT NOT NULL,
    to_currency_id INT NOT NULL,
    rate DECIMAL(20, 8) NOT NULL,
    reverse_rate DECIMAL(20, 8) NOT NULL,
    source VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (from_currency_id) REFERENCES currencies(id),
    FOREIGN KEY (to_currency_id) REFERENCES currencies(id),
    UNIQUE KEY unique_pair (from_currency_id, to_currency_id)
);

-- Статусы заявок
CREATE TABLE order_statuses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255)
);

-- Таблица заявок на обмен
CREATE TABLE exchange_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_number VARCHAR(20) NOT NULL UNIQUE,
    from_currency_id INT NOT NULL,
    to_currency_id INT NOT NULL,
    from_amount DECIMAL(20, 8) NOT NULL,
    to_amount DECIMAL(20, 8) NOT NULL,
    rate DECIMAL(20, 8) NOT NULL,
    status_id INT NOT NULL DEFAULT 1,
    
    -- Контактная информация клиента
    client_phone VARCHAR(20),
    client_name VARCHAR(255),
    client_email VARCHAR(255),
    client_telegram VARCHAR(100),
    
    -- Реквизиты для получения
    recipient_wallet VARCHAR(255),
    recipient_details TEXT,
    
    -- Системная информация
    ip_address VARCHAR(45),
    user_agent TEXT,
    referrer VARCHAR(255),
    
    -- Временные метки
    expires_at TIMESTAMP,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (from_currency_id) REFERENCES currencies(id),
    FOREIGN KEY (to_currency_id) REFERENCES currencies(id),
    FOREIGN KEY (status_id) REFERENCES order_statuses(id),
    
    INDEX idx_order_number (order_number),
    INDEX idx_status (status_id),
    INDEX idx_created_at (created_at),
    INDEX idx_client_email (client_email)
);

-- Таблица логов операций
CREATE TABLE operation_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT,
    action VARCHAR(100) NOT NULL,
    description TEXT,
    old_data JSON,
    new_data JSON,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (order_id) REFERENCES exchange_orders(id),
    INDEX idx_order_id (order_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
);

-- Таблица настроек системы
CREATE TABLE system_settings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    setting_key VARCHAR(100) NOT NULL UNIQUE,
    setting_value TEXT,
    description VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Заполнение начальными данными

-- Статусы заявок
INSERT INTO order_statuses (name, description) VALUES
('pending', 'Ожидает обработки'),
('confirmed', 'Подтверждена'),
('processing', 'В обработке'),
('completed', 'Завершена'),
('cancelled', 'Отменена'),
('expired', 'Истекла');

-- Валюты
INSERT INTO currencies (code, name, symbol, network, min_amount, decimals) VALUES
('USDT_TRC20', 'Tether USDT', 'USDT', 'TRC20', 10.00000000, 6),
('RUB_TBANK', 'Российский рубль', 'RUB', 'Т-Банк', 500.00000000, 2),
('BTC', 'Bitcoin', 'BTC', 'Bitcoin', 0.00010000, 8),
('ETH', 'Ethereum', 'ETH', 'Ethereum', 0.01000000, 18),
('USDT_ERC20', 'Tether USDT', 'USDT', 'ERC20', 10.00000000, 6),
('USDC', 'USD Coin', 'USDC', 'ERC20', 10.00000000, 6);

-- Системные настройки
INSERT INTO system_settings (setting_key, setting_value, description) VALUES
('site_name', 'CryptoExchange', 'Название сайта'),
('admin_email', 'admin@cryptoexchange.com', 'Email администратора'),
('order_expiry_minutes', '30', 'Время жизни заявки в минутах'),
('min_order_amount_usd', '10', 'Минимальная сумма заявки в USD'),
('max_order_amount_usd', '50000', 'Максимальная сумма заявки в USD'),
('commission_percent', '1.5', 'Комиссия в процентах'),
('api_update_interval', '60', 'Интервал обновления курсов в секундах');