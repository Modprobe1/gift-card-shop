package services

import (
	"crypto-exchange-backend/internal/models"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderService struct {
	db              *gorm.DB
	currencyService *CurrencyService
	rateService     *RateService
}

func NewOrderService(db *gorm.DB, currencyService *CurrencyService, rateService *RateService) *OrderService {
	return &OrderService{
		db:              db,
		currencyService: currencyService,
		rateService:     rateService,
	}
}

// CreateOrder создает новую заявку на обмен
func (s *OrderService) CreateOrder(req *models.CreateOrderRequest, ipAddress, userAgent string) (*models.ExchangeOrder, error) {
	// Получаем валюты
	fromCurrency, err := s.currencyService.GetCurrencyByCode(req.FromCurrencyCode)
	if err != nil {
		return nil, fmt.Errorf("from currency not found: %w", err)
	}

	toCurrency, err := s.currencyService.GetCurrencyByCode(req.ToCurrencyCode)
	if err != nil {
		return nil, fmt.Errorf("to currency not found: %w", err)
	}

	// Вычисляем сумму обмена
	calculation, err := s.rateService.Calculate(req.FromCurrencyCode, req.ToCurrencyCode, req.FromAmount)
	if err != nil {
		return nil, fmt.Errorf("calculation failed: %w", err)
	}

	// Получаем время истечения заявки
	expiryMinutes, err := s.getOrderExpiryMinutes()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	// Создаем заявку
	order := &models.ExchangeOrder{
		OrderNumber:      s.generateOrderNumber(),
		FromCurrencyID:   fromCurrency.ID,
		ToCurrencyID:     toCurrency.ID,
		FromAmount:       req.FromAmount,
		ToAmount:         calculation.ToAmount,
		Rate:             calculation.Rate,
		StatusID:         1, // pending
		ClientPhone:      req.ClientPhone,
		ClientName:       req.ClientName,
		ClientEmail:      req.ClientEmail,
		ClientTelegram:   req.ClientTelegram,
		RecipientWallet:  req.RecipientWallet,
		RecipientDetails: req.RecipientDetails,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		ExpiresAt:        &expiresAt,
	}

	// Сохраняем заявку в транзакции
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Логируем создание заявки
		log := &models.OperationLog{
			OrderID:     &order.ID,
			Action:      "order_created",
			Description: "Order created by user",
			NewData:     s.orderToJSON(order),
			IPAddress:   ipAddress,
			UserAgent:   userAgent,
		}

		return tx.Create(log).Error
	})

	if err != nil {
		return nil, err
	}

	// Загружаем связанные данные
	if err := s.db.Preload("FromCurrency").Preload("ToCurrency").Preload("Status").
		First(order, order.ID).Error; err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
		"from":         req.FromCurrencyCode,
		"to":           req.ToCurrencyCode,
		"amount":       req.FromAmount,
	}).Info("Order created")

	return order, nil
}

// GetOrderByNumber возвращает заявку по номеру
func (s *OrderService) GetOrderByNumber(orderNumber string) (*models.ExchangeOrder, error) {
	var order models.ExchangeOrder

	err := s.db.Preload("FromCurrency").Preload("ToCurrency").Preload("Status").
		Where("order_number = ?", orderNumber).First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// GetOrderByID возвращает заявку по ID
func (s *OrderService) GetOrderByID(id uint) (*models.ExchangeOrder, error) {
	var order models.ExchangeOrder

	err := s.db.Preload("FromCurrency").Preload("ToCurrency").Preload("Status").
		First(&order, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return &order, nil
}

// GetOrders возвращает список заявок с пагинацией
func (s *OrderService) GetOrders(page, limit int, status string) ([]models.ExchangeOrder, int64, error) {
	var orders []models.ExchangeOrder
	var total int64

	query := s.db.Model(&models.ExchangeOrder{}).
		Preload("FromCurrency").
		Preload("ToCurrency").
		Preload("Status")

	// Фильтр по статусу
	if status != "" {
		query = query.Joins("JOIN order_statuses ON order_statuses.id = exchange_orders.status_id").
			Where("order_statuses.name = ?", status)
	}

	// Считаем общее количество
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Получаем заявки с пагинацией
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateOrderStatus обновляет статус заявки
func (s *OrderService) UpdateOrderStatus(orderID uint, statusName string, ipAddress, userAgent string) error {
	// Получаем статус
	var status models.OrderStatus
	if err := s.db.Where("name = ?", statusName).First(&status).Error; err != nil {
		return fmt.Errorf("status not found: %w", err)
	}

	// Получаем заявку
	order, err := s.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	oldStatus := order.Status.Name

	// Обновляем в транзакции
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Обновляем статус
		if err := tx.Model(order).Update("status_id", status.ID).Error; err != nil {
			return err
		}

		// Если статус "completed", устанавливаем время завершения
		if statusName == "completed" {
			now := time.Now()
			if err := tx.Model(order).Update("completed_at", now).Error; err != nil {
				return err
			}
		}

		// Логируем изменение статуса
		log := &models.OperationLog{
			OrderID:     &order.ID,
			Action:      "status_updated",
			Description: fmt.Sprintf("Status changed from %s to %s", oldStatus, statusName),
			OldData:     fmt.Sprintf(`{"status": "%s"}`, oldStatus),
			NewData:     fmt.Sprintf(`{"status": "%s"}`, statusName),
			IPAddress:   ipAddress,
			UserAgent:   userAgent,
		}

		return tx.Create(log).Error
	})
}

// GetStatistics возвращает статистику заявок
func (s *OrderService) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Общее количество заявок
	var totalOrders int64
	if err := s.db.Model(&models.ExchangeOrder{}).Count(&totalOrders).Error; err != nil {
		return nil, err
	}
	stats["total_orders"] = totalOrders

	// Заявки по статусам
	var statusStats []struct {
		StatusName string `json:"status_name"`
		Count      int64  `json:"count"`
	}

	err := s.db.Model(&models.ExchangeOrder{}).
		Select("order_statuses.name as status_name, COUNT(*) as count").
		Joins("JOIN order_statuses ON order_statuses.id = exchange_orders.status_id").
		Group("order_statuses.name").
		Scan(&statusStats).Error

	if err != nil {
		return nil, err
	}
	stats["status_stats"] = statusStats

	// Заявки за сегодня
	today := time.Now().Truncate(24 * time.Hour)
	var todayOrders int64
	if err := s.db.Model(&models.ExchangeOrder{}).
		Where("created_at >= ?", today).Count(&todayOrders).Error; err != nil {
		return nil, err
	}
	stats["today_orders"] = todayOrders

	// Заявки за неделю
	weekAgo := time.Now().AddDate(0, 0, -7)
	var weekOrders int64
	if err := s.db.Model(&models.ExchangeOrder{}).
		Where("created_at >= ?", weekAgo).Count(&weekOrders).Error; err != nil {
		return nil, err
	}
	stats["week_orders"] = weekOrders

	return stats, nil
}

// CheckExpiredOrders проверяет и помечает истекшие заявки
func (s *OrderService) CheckExpiredOrders() error {
	now := time.Now()
	
	// Получаем ID статуса "expired"
	var expiredStatus models.OrderStatus
	if err := s.db.Where("name = ?", "expired").First(&expiredStatus).Error; err != nil {
		return err
	}

	// Обновляем истекшие заявки
	result := s.db.Model(&models.ExchangeOrder{}).
		Where("expires_at < ? AND status_id = ?", now, 1). // 1 = pending
		Update("status_id", expiredStatus.ID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		logrus.WithField("count", result.RowsAffected).Info("Marked orders as expired")
	}

	return nil
}

// generateOrderNumber генерирует уникальный номер заявки
func (s *OrderService) generateOrderNumber() string {
	timestamp := time.Now().Unix()
	random := rand.Intn(9999)
	return fmt.Sprintf("ORD%d%04d", timestamp, random)
}

// getOrderExpiryMinutes получает время жизни заявки из настроек
func (s *OrderService) getOrderExpiryMinutes() (int, error) {
	var setting models.SystemSetting

	err := s.db.Where("setting_key = ? AND is_active = ?", "order_expiry_minutes", true).
		First(&setting).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 30, nil // Дефолтное время 30 минут
		}
		return 0, err
	}

	minutes, err := strconv.Atoi(setting.SettingValue)
	if err != nil {
		return 30, nil // Дефолтное время при ошибке парсинга
	}

	return minutes, nil
}

// orderToJSON конвертирует заявку в JSON для логирования
func (s *OrderService) orderToJSON(order *models.ExchangeOrder) string {
	// Упрощенная версия для логирования
	return fmt.Sprintf(`{
		"id": %d,
		"order_number": "%s",
		"from_currency_id": %d,
		"to_currency_id": %d,
		"from_amount": %f,
		"to_amount": %f,
		"status_id": %d
	}`, order.ID, order.OrderNumber, order.FromCurrencyID, order.ToCurrencyID,
		order.FromAmount, order.ToAmount, order.StatusID)
}