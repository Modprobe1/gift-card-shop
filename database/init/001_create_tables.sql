-- Удаление существующих таблиц (если они есть)
DROP TABLE IF EXISTS operation_logs;
DROP TABLE IF EXISTS exchange_orders;
DROP TABLE IF EXISTS order_statuses;
DROP TABLE IF EXISTS exchange_rates;
DROP TABLE IF EXISTS system_settings;
DROP TABLE IF EXISTS currencies;

-- Создание таблицы валют
CREATE TABLE currencies (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    network VARCHAR(50),
    min_amount DECIMAL(18,8) DEFAULT 0,
    max_amount DECIMAL(18,8) NULL,
    is_crypto BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Создание таблицы курсов обмена
CREATE TABLE exchange_rates (
    id INT PRIMARY KEY AUTO_INCREMENT,
    from_currency_id INT NOT NULL,
    to_currency_id INT NOT NULL,
    rate DECIMAL(18,8) NOT NULL,
    reverse_rate DECIMAL(18,8) NOT NULL,
    source VARCHAR(50) DEFAULT 'manual',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (from_currency_id) REFERENCES currencies(id),
    FOREIGN KEY (to_currency_id) REFERENCES currencies(id),
    UNIQUE KEY unique_currency_pair (from_currency_id, to_currency_id)
);

-- Создание таблицы статусов заказов
CREATE TABLE order_statuses (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(7) DEFAULT '#007bff',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы заказов на обмен
CREATE TABLE exchange_orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_number VARCHAR(20) UNIQUE NOT NULL,
    from_currency_id INT NOT NULL,
    to_currency_id INT NOT NULL,
    from_amount DECIMAL(18,8) NOT NULL,
    to_amount DECIMAL(18,8) NOT NULL,
    rate DECIMAL(18,8) NOT NULL,
    commission DECIMAL(18,8) NOT NULL,
    commission_rate DECIMAL(5,2) NOT NULL,
    status_id INT NOT NULL,
    client_name VARCHAR(255) NOT NULL,
    client_phone VARCHAR(50) NOT NULL,
    client_email VARCHAR(255) NOT NULL,
    client_telegram VARCHAR(100),
    recipient_wallet TEXT NOT NULL,
    recipient_details TEXT,
    expires_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (from_currency_id) REFERENCES currencies(id),
    FOREIGN KEY (to_currency_id) REFERENCES currencies(id),
    FOREIGN KEY (status_id) REFERENCES order_statuses(id),
    INDEX idx_order_number (order_number),
    INDEX idx_client_email (client_email),
    INDEX idx_status (status_id),
    INDEX idx_created_at (created_at)
);

-- Создание таблицы системных настроек
CREATE TABLE system_settings (
    id INT PRIMARY KEY AUTO_INCREMENT,
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Создание таблицы логов операций
CREATE TABLE operation_logs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT,
    operation_type VARCHAR(50) NOT NULL,
    description TEXT,
    old_status_id INT,
    new_status_id INT,
    user_id INT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES exchange_orders(id),
    FOREIGN KEY (old_status_id) REFERENCES order_statuses(id),
    FOREIGN KEY (new_status_id) REFERENCES order_statuses(id),
    INDEX idx_order_id (order_id),
    INDEX idx_created_at (created_at)
);

-- Вставка валют (только USDT, BTC, RUB_TBANK)
INSERT INTO currencies (name, code, symbol, network, min_amount, max_amount, is_crypto, is_active) VALUES
('Tether USD', 'USDT_TRC20', 'USDT', 'TRC20', 10.00, 50000.00, TRUE, TRUE),
('Bitcoin', 'BTC', 'BTC', 'Bitcoin', 0.001, 10.0, TRUE, TRUE),
('Рубли Т-Банк', 'RUB_TBANK', 'RUB', 'T-Bank', 1000.00, 2000000.00, FALSE, TRUE);

-- Вставка статусов заказов
INSERT INTO order_statuses (name, description, color) VALUES
('pending', 'Ожидает оплаты', '#ffc107'),
('confirmed', 'Оплата подтверждена', '#17a2b8'),
('processing', 'В обработке', '#007bff'),
('completed', 'Завершен', '#28a745'),
('cancelled', 'Отменен', '#dc3545'),
('expired', 'Истек срок', '#6c757d');

-- Вставка системных настроек с комиссией 2%
INSERT INTO system_settings (setting_key, setting_value, description) VALUES
('commission_percent', '2.0', 'Процент комиссии за обмен'),
('order_expiry_hours', '24', 'Время жизни заказа в часах'),
('min_confirmations_btc', '3', 'Минимальное количество подтверждений для BTC'),
('min_confirmations_usdt', '12', 'Минимальное количество подтверждений для USDT'),
('support_email', 'support@crypto-exchange.com', 'Email поддержки'),
('support_telegram', '@crypto_exchange_support', 'Telegram поддержки');

-- Создание начальных курсов обмена (будут обновляться из API)
-- USDT -> BTC
INSERT INTO exchange_rates (from_currency_id, to_currency_id, rate, reverse_rate, source) VALUES
((SELECT id FROM currencies WHERE code = 'USDT_TRC20'), (SELECT id FROM currencies WHERE code = 'BTC'), 0.000014, 71428.57, 'manual');

-- BTC -> USDT  
INSERT INTO exchange_rates (from_currency_id, to_currency_id, rate, reverse_rate, source) VALUES
((SELECT id FROM currencies WHERE code = 'BTC'), (SELECT id FROM currencies WHERE code = 'USDT_TRC20'), 71428.57, 0.000014, 'manual');

-- USDT -> RUB_TBANK
INSERT INTO exchange_rates (from_currency_id, to_currency_id, rate, reverse_rate, source) VALUES
((SELECT id FROM currencies WHERE code = 'USDT_TRC20'), (SELECT id FROM currencies WHERE code = 'RUB_TBANK'), 96.50, 0.0104, 'manual');

-- RUB_TBANK -> USDT
INSERT INTO exchange_rates (from_currency_id, to_currency_id, rate, reverse_rate, source) VALUES
((SELECT id FROM currencies WHERE code = 'RUB_TBANK'), (SELECT id FROM currencies WHERE code = 'USDT_TRC20'), 0.0104, 96.50, 'manual');