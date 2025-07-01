package handlers

import (
	"crypto-exchange-backend/internal/models"
	"crypto-exchange-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	currencyService *services.CurrencyService
	rateService     *services.RateService
	orderService    *services.OrderService
}

func NewHandler(currencyService *services.CurrencyService, rateService *services.RateService, orderService *services.OrderService) *Handler {
	return &Handler{
		currencyService: currencyService,
		rateService:     rateService,
		orderService:    orderService,
	}
}

// GetCurrencies возвращает список всех валют
func (h *Handler) GetCurrencies(c *gin.Context) {
	currencies, err := h.currencyService.GetAllCurrencies()
	if err != nil {
		logrus.WithError(err).Error("Failed to get currencies")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get currencies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    currencies,
	})
}

// GetCurrency возвращает валюту по коду
func (h *Handler) GetCurrency(c *gin.Context) {
	code := c.Param("code")

	currency, err := h.currencyService.GetCurrencyByCode(code)
	if err != nil {
		if err.Error() == "currency not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Currency not found"})
			return
		}
		logrus.WithError(err).Error("Failed to get currency")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get currency"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    currency,
	})
}

// GetExchangeRates возвращает все курсы обмена
func (h *Handler) GetExchangeRates(c *gin.Context) {
	rates, err := h.rateService.GetAllExchangeRates()
	if err != nil {
		logrus.WithError(err).Error("Failed to get exchange rates")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exchange rates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rates,
	})
}

// GetExchangeRate возвращает курс обмена между двумя валютами
func (h *Handler) GetExchangeRate(c *gin.Context) {
	from := c.Param("from")
	to := c.Param("to")

	rate, err := h.rateService.GetExchangeRate(from, to)
	if err != nil {
		if err.Error() == "exchange rate not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exchange rate not found"})
			return
		}
		logrus.WithError(err).Error("Failed to get exchange rate")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exchange rate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rate,
	})
}

// Calculate вычисляет сумму обмена
func (h *Handler) Calculate(c *gin.Context) {
	var req struct {
		FromCurrency string  `json:"from_currency" binding:"required"`
		ToCurrency   string  `json:"to_currency" binding:"required"`
		FromAmount   float64 `json:"from_amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	calculation, err := h.rateService.Calculate(req.FromCurrency, req.ToCurrency, req.FromAmount)
	if err != nil {
		logrus.WithError(err).Error("Failed to calculate exchange")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    calculation,
	})
}

// CreateOrder создает новую заявку на обмен
func (h *Handler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем IP адрес и User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	order, err := h.orderService.CreateOrder(&req, ipAddress, userAgent)
	if err != nil {
		logrus.WithError(err).Error("Failed to create order")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    order,
	})
}

// GetOrder возвращает заявку по номеру
func (h *Handler) GetOrder(c *gin.Context) {
	orderNumber := c.Param("number")

	order, err := h.orderService.GetOrderByNumber(orderNumber)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		logrus.WithError(err).Error("Failed to get order")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    order,
	})
}

// GetOrders возвращает список заявок (админский эндпоинт)
func (h *Handler) GetOrders(c *gin.Context) {
	// Параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	orders, total, err := h.orderService.GetOrders(page, limit, status)
	if err != nil {
		logrus.WithError(err).Error("Failed to get orders")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"orders": orders,
			"total":  total,
			"page":   page,
			"limit":  limit,
		},
	})
}

// UpdateOrderStatus обновляет статус заявки (админский эндпоинт)
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем IP адрес и User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err = h.orderService.UpdateOrderStatus(uint(orderID), req.Status, ipAddress, userAgent)
	if err != nil {
		logrus.WithError(err).Error("Failed to update order status")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order status updated successfully",
	})
}

// GetStatistics возвращает статистику (админский эндпоинт)
func (h *Handler) GetStatistics(c *gin.Context) {
	stats, err := h.orderService.GetStatistics()
	if err != nil {
		logrus.WithError(err).Error("Failed to get statistics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}