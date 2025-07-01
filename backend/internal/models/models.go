package models

import (
	"time"

	"gorm.io/gorm"
)

// Currency представляет валюту в системе
type Currency struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Code      string         `json:"code" gorm:"uniqueIndex;size:10;not null"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Symbol    string         `json:"symbol" gorm:"size:10"`
	Network   string         `json:"network" gorm:"size:50"`
	IconURL   string         `json:"icon_url" gorm:"size:255"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	MinAmount float64        `json:"min_amount" gorm:"type:decimal(20,8);default:0"`
	MaxAmount *float64       `json:"max_amount" gorm:"type:decimal(20,8)"`
	Decimals  int            `json:"decimals" gorm:"default:8"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ExchangeRate представляет курс обмена между валютами
type ExchangeRate struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	FromCurrencyID uint           `json:"from_currency_id" gorm:"not null"`
	ToCurrencyID   uint           `json:"to_currency_id" gorm:"not null"`
	Rate           float64        `json:"rate" gorm:"type:decimal(20,8);not null"`
	ReverseRate    float64        `json:"reverse_rate" gorm:"type:decimal(20,8);not null"`
	Source         string         `json:"source" gorm:"size:50;not null"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	FromCurrency Currency `json:"from_currency" gorm:"foreignKey:FromCurrencyID"`
	ToCurrency   Currency `json:"to_currency" gorm:"foreignKey:ToCurrencyID"`
}

// OrderStatus представляет статус заявки
type OrderStatus struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Description string         `json:"description" gorm:"size:255"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// ExchangeOrder представляет заявку на обмен
type ExchangeOrder struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OrderNumber    string         `json:"order_number" gorm:"uniqueIndex;size:20;not null"`
	FromCurrencyID uint           `json:"from_currency_id" gorm:"not null"`
	ToCurrencyID   uint           `json:"to_currency_id" gorm:"not null"`
	FromAmount     float64        `json:"from_amount" gorm:"type:decimal(20,8);not null"`
	ToAmount       float64        `json:"to_amount" gorm:"type:decimal(20,8);not null"`
	Rate           float64        `json:"rate" gorm:"type:decimal(20,8);not null"`
	StatusID       uint           `json:"status_id" gorm:"not null;default:1"`

	// Контактная информация клиента
	ClientPhone    string `json:"client_phone" gorm:"size:20"`
	ClientName     string `json:"client_name" gorm:"size:255"`
	ClientEmail    string `json:"client_email" gorm:"size:255;index"`
	ClientTelegram string `json:"client_telegram" gorm:"size:100"`

	// Реквизиты для получения
	RecipientWallet  string `json:"recipient_wallet" gorm:"size:255"`
	RecipientDetails string `json:"recipient_details" gorm:"type:text"`

	// Системная информация
	IPAddress string `json:"ip_address" gorm:"size:45"`
	UserAgent string `json:"user_agent" gorm:"type:text"`
	Referrer  string `json:"referrer" gorm:"size:255"`

	// Временные метки
	ExpiresAt   *time.Time     `json:"expires_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Связи
	FromCurrency Currency    `json:"from_currency" gorm:"foreignKey:FromCurrencyID"`
	ToCurrency   Currency    `json:"to_currency" gorm:"foreignKey:ToCurrencyID"`
	Status       OrderStatus `json:"status" gorm:"foreignKey:StatusID"`
}

// OperationLog представляет лог операций
type OperationLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OrderID     *uint     `json:"order_id"`
	Action      string    `json:"action" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	OldData     string    `json:"old_data" gorm:"type:json"`
	NewData     string    `json:"new_data" gorm:"type:json"`
	IPAddress   string    `json:"ip_address" gorm:"size:45"`
	UserAgent   string    `json:"user_agent" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`

	// Связи
	Order *ExchangeOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}

// SystemSetting представляет настройки системы
type SystemSetting struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	SettingKey   string         `json:"setting_key" gorm:"uniqueIndex;size:100;not null"`
	SettingValue string         `json:"setting_value" gorm:"type:text"`
	Description  string         `json:"description" gorm:"size:255"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateOrderRequest представляет запрос на создание заявки
type CreateOrderRequest struct {
	FromCurrencyCode string  `json:"from_currency" binding:"required"`
	ToCurrencyCode   string  `json:"to_currency" binding:"required"`
	FromAmount       float64 `json:"from_amount" binding:"required,gt=0"`
	
	// Контактная информация
	ClientPhone    string `json:"client_phone" binding:"required"`
	ClientName     string `json:"client_name" binding:"required"`
	ClientEmail    string `json:"client_email" binding:"required,email"`
	ClientTelegram string `json:"client_telegram"`
	
	// Реквизиты для получения
	RecipientWallet  string `json:"recipient_wallet" binding:"required"`
	RecipientDetails string `json:"recipient_details"`
}

// ExchangeRateResponse представляет ответ с курсом обмена
type ExchangeRateResponse struct {
	FromCurrency Currency `json:"from_currency"`
	ToCurrency   Currency `json:"to_currency"`
	Rate         float64  `json:"rate"`
	ReverseRate  float64  `json:"reverse_rate"`
	MinAmount    float64  `json:"min_amount"`
	MaxAmount    *float64 `json:"max_amount"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CalculateResponse представляет ответ калькулятора
type CalculateResponse struct {
	FromAmount     float64                `json:"from_amount"`
	ToAmount       float64                `json:"to_amount"`
	Rate           float64                `json:"rate"`
	Commission     float64                `json:"commission"`
	CommissionRate float64                `json:"commission_rate"`
	ExchangeRate   *ExchangeRateResponse  `json:"exchange_rate"`
}