-- Отключаем проверку foreign keys
SET FOREIGN_KEY_CHECKS = 0;

-- Удаляем все таблицы если они существуют
DROP TABLE IF EXISTS `operation_logs`;
DROP TABLE IF EXISTS `exchange_orders`;
DROP TABLE IF EXISTS `exchange_rates`;
DROP TABLE IF EXISTS `system_settings`;
DROP TABLE IF EXISTS `order_statuses`;
DROP TABLE IF EXISTS `currencies`;

-- Включаем обратно проверку foreign keys
SET FOREIGN_KEY_CHECKS = 1;

-- Теперь таблицы будут созданы заново при запуске backend