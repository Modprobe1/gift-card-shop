package database

import (
	"crypto-exchange-backend/internal/config"
	"crypto-exchange-backend/internal/models"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect устанавливает соединение с базой данных
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// Настройки GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Получаем базовый объект *sql.DB для настройки пула соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// Migrate выполняет миграции базы данных
func Migrate(db *gorm.DB) error {
	// Автоматические миграции для всех моделей
	if err := db.AutoMigrate(
		&models.Currency{},
		&models.ExchangeRate{},
		&models.OrderStatus{},
		&models.ExchangeOrder{},
		&models.OperationLog{},
		&models.SystemSetting{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Заполняем начальными данными, если нужно
	if err := seedData(db); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	return nil
}

// seedData заполняет базу начальными данными
func seedData(db *gorm.DB) error {
	// Проверяем, есть ли уже данные
	var count int64
	if err := db.Model(&models.OrderStatus{}).Count(&count).Error; err != nil {
		return err
	}

	// Если данные уже есть, пропускаем
	if count > 0 {
		return nil
	}

	// Создаем статусы заявок
	statuses := []models.OrderStatus{
		{Name: "pending", Description: "Ожидает обработки"},
		{Name: "confirmed", Description: "Подтверждена"},
		{Name: "processing", Description: "В обработке"},
		{Name: "completed", Description: "Завершена"},
		{Name: "cancelled", Description: "Отменена"},
		{Name: "expired", Description: "Истекла"},
	}

	for _, status := range statuses {
		if err := db.Create(&status).Error; err != nil {
			return fmt.Errorf("failed to create status %s: %w", status.Name, err)
		}
	}

	// Создаем валюты
	currencies := []models.Currency{
		{
			Code:      "USDT_TRC20",
			Name:      "Tether USDT",
			Symbol:    "USDT",
			Network:   "TRC20",
			MinAmount: 10.0,
			Decimals:  6,
			IsActive:  true,
		},
		{
			Code:      "RUB_TBANK",
			Name:      "Российский рубль",
			Symbol:    "RUB",
			Network:   "Т-Банк",
			MinAmount: 500.0,
			Decimals:  2,
			IsActive:  true,
		},
		{
			Code:      "BTC",
			Name:      "Bitcoin",
			Symbol:    "BTC",
			Network:   "Bitcoin",
			MinAmount: 0.0001,
			Decimals:  8,
			IsActive:  true,
		},
		{
			Code:      "ETH",
			Name:      "Ethereum",
			Symbol:    "ETH",
			Network:   "Ethereum",
			MinAmount: 0.01,
			Decimals:  18,
			IsActive:  true,
		},
	}

	for _, currency := range currencies {
		if err := db.Create(&currency).Error; err != nil {
			return fmt.Errorf("failed to create currency %s: %w", currency.Code, err)
		}
	}

	// Создаем системные настройки
	settings := []models.SystemSetting{
		{SettingKey: "site_name", SettingValue: "CryptoExchange", Description: "Название сайта"},
		{SettingKey: "admin_email", SettingValue: "admin@cryptoexchange.com", Description: "Email администратора"},
		{SettingKey: "order_expiry_minutes", SettingValue: "30", Description: "Время жизни заявки в минутах"},
		{SettingKey: "commission_percent", SettingValue: "1.5", Description: "Комиссия в процентах"},
		{SettingKey: "api_update_interval", SettingValue: "60", Description: "Интервал обновления курсов в секундах"},
	}

	for _, setting := range settings {
		if err := db.Create(&setting).Error; err != nil {
			return fmt.Errorf("failed to create setting %s: %w", setting.SettingKey, err)
		}
	}

	return nil
}