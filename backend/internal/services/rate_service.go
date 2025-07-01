package services

import (
	"crypto-exchange-backend/internal/config"
	"crypto-exchange-backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RateService struct {
	db     *gorm.DB
	config *config.Config
	client *http.Client
}

// CoinGeckoResponse структура ответа от CoinGecko API
type CoinGeckoResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
		RUB float64 `json:"rub"`
	} `json:"bitcoin"`
	Tether struct {
		USD float64 `json:"usd"`
		RUB float64 `json:"rub"`
	} `json:"tether"`
}

func NewRateService(db *gorm.DB, cfg *config.Config) *RateService {
	return &RateService{
		db:     db,
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.API.RequestTimeout) * time.Second,
		},
	}
}

// GetExchangeRate возвращает курс обмена между двумя валютами
func (s *RateService) GetExchangeRate(fromCode, toCode string) (*models.ExchangeRateResponse, error) {
	var rate models.ExchangeRate
	
	// Ищем курс в базе данных
	err := s.db.Preload("FromCurrency").Preload("ToCurrency").
		Joins("JOIN currencies fc ON fc.id = exchange_rates.from_currency_id").
		Joins("JOIN currencies tc ON tc.id = exchange_rates.to_currency_id").
		Where("fc.code = ? AND tc.code = ? AND exchange_rates.is_active = ?", fromCode, toCode, true).
		First(&rate).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("exchange rate not found")
		}
		return nil, err
	}
	
	response := &models.ExchangeRateResponse{
		FromCurrency: rate.FromCurrency,
		ToCurrency:   rate.ToCurrency,
		Rate:         rate.Rate,
		ReverseRate:  rate.ReverseRate,
		MinAmount:    rate.FromCurrency.MinAmount,
		MaxAmount:    rate.FromCurrency.MaxAmount,
		UpdatedAt:    rate.UpdatedAt,
	}
	
	return response, nil
}

// GetAllExchangeRates возвращает все активные курсы обмена
func (s *RateService) GetAllExchangeRates() ([]models.ExchangeRateResponse, error) {
	var rates []models.ExchangeRate
	
	err := s.db.Preload("FromCurrency").Preload("ToCurrency").
		Where("is_active = ?", true).Find(&rates).Error
	
	if err != nil {
		return nil, err
	}
	
	var responses []models.ExchangeRateResponse
	for _, rate := range rates {
		response := models.ExchangeRateResponse{
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         rate.Rate,
			ReverseRate:  rate.ReverseRate,
			MinAmount:    rate.FromCurrency.MinAmount,
			MaxAmount:    rate.FromCurrency.MaxAmount,
			UpdatedAt:    rate.UpdatedAt,
		}
		responses = append(responses, response)
	}
	
	return responses, nil
}

// Calculate вычисляет сумму обмена с учетом комиссии
func (s *RateService) Calculate(fromCode, toCode string, fromAmount float64) (*models.CalculateResponse, error) {
	// Получаем курс обмена
	exchangeRate, err := s.GetExchangeRate(fromCode, toCode)
	if err != nil {
		return nil, err
	}
	
	// Проверяем минимальную сумму
	if fromAmount < exchangeRate.MinAmount {
		return nil, fmt.Errorf("amount is below minimum: %f", exchangeRate.MinAmount)
	}
	
	// Проверяем максимальную сумму
	if exchangeRate.MaxAmount != nil && fromAmount > *exchangeRate.MaxAmount {
		return nil, fmt.Errorf("amount exceeds maximum: %f", *exchangeRate.MaxAmount)
	}
	
	// Получаем комиссию из настроек (2%)
	commissionRate, err := s.getCommissionRate()
	if err != nil {
		return nil, err
	}
	
	// Вычисляем сумму без комиссии
	baseToAmount := fromAmount * exchangeRate.Rate
	
	// Вычисляем комиссию
	commission := baseToAmount * (commissionRate / 100)
	
	// Итоговая сумма к получению (уже с вычетом комиссии)
	toAmount := baseToAmount - commission
	
	response := &models.CalculateResponse{
		FromAmount:     fromAmount,
		ToAmount:       toAmount,
		Rate:           exchangeRate.Rate,
		Commission:     commission,
		CommissionRate: commissionRate,
		ExchangeRate:   exchangeRate,
	}
	
	return response, nil
}

// UpdateExchangeRate обновляет курс обмена
func (s *RateService) UpdateExchangeRate(fromCurrencyID, toCurrencyID uint, rate float64, source string) error {
	reverseRate := 1.0 / rate
	
	exchangeRate := models.ExchangeRate{
		FromCurrencyID: fromCurrencyID,
		ToCurrencyID:   toCurrencyID,
		Rate:           rate,
		ReverseRate:    reverseRate,
		Source:         source,
		IsActive:       true,
	}
	
	// Используем ON DUPLICATE KEY UPDATE логику
	return s.db.Save(&exchangeRate).Error
}

// StartRateUpdater запускает периодическое обновление курсов
func (s *RateService) StartRateUpdater() {
	ticker := time.NewTicker(time.Duration(s.config.API.UpdateInterval) * time.Second)
	defer ticker.Stop()
	
	logrus.Info("Starting rate updater for USDT-BTC and USDT-TBANK pairs")
	
	// Обновляем курсы сразу при запуске
	if err := s.UpdateRatesFromAPI(); err != nil {
		logrus.WithError(err).Error("Failed to update rates on startup")
	}
	
	for {
		select {
		case <-ticker.C:
			if err := s.UpdateRatesFromAPI(); err != nil {
				logrus.WithError(err).Error("Failed to update rates")
			}
		}
	}
}

// UpdateRatesFromAPI обновляет курсы из CoinGecko API с добавлением 2% профита
func (s *RateService) UpdateRatesFromAPI() error {
	logrus.Debug("Updating exchange rates from CoinGecko API")
	
	// Запрос к CoinGecko API (только bitcoin и tether)
	url := fmt.Sprintf("%s/simple/price?ids=bitcoin,tether&vs_currencies=usd,rub", s.config.API.CoinGeckoURL)
	
	resp, err := s.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch rates from API: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	var apiResponse CoinGeckoResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}
	
	// Получаем валюты из базы
	currencies := make(map[string]uint)
	var currencyList []models.Currency
	if err := s.db.Find(&currencyList).Error; err != nil {
		return fmt.Errorf("failed to fetch currencies: %w", err)
	}
	
	for _, currency := range currencyList {
		currencies[currency.Code] = currency.ID
	}
	
	// Применяем 2% профит к курсам
	profitMargin := 1.02 // 2% профит
	
	// Обновляем только нужные курсы
	rates := []struct {
		from, to string
		rate     float64
	}{
		// USDT -> BTC пара (основная)
		{"USDT_TRC20", "BTC", (1.0 / apiResponse.Bitcoin.USD) / profitMargin}, // USDT -> BTC (с вычетом профита)
		{"BTC", "USDT_TRC20", apiResponse.Bitcoin.USD * profitMargin},         // BTC -> USDT (с добавлением профита)
		
		// USDT -> TBANK пара (дополнительная)
		{"USDT_TRC20", "RUB_TBANK", apiResponse.Tether.RUB * profitMargin},        // USDT -> RUB (с добавлением профита)
		{"RUB_TBANK", "USDT_TRC20", (1.0 / apiResponse.Tether.RUB) / profitMargin}, // RUB -> USDT (с вычетом профита)
	}
	
	for _, r := range rates {
		fromID, fromExists := currencies[r.from]
		toID, toExists := currencies[r.to]
		
		if fromExists && toExists {
			if err := s.UpdateExchangeRate(fromID, toID, r.rate, "coingecko"); err != nil {
				logrus.WithError(err).Errorf("Failed to update rate %s -> %s", r.from, r.to)
			} else {
				logrus.Debugf("Updated rate %s -> %s: %.8f", r.from, r.to, r.rate)
			}
		}
	}
	
	logrus.Info("Successfully updated exchange rates with 2% profit margin")
	return nil
}

// getCommissionRate получает процент комиссии из настроек (должно быть 2%)
func (s *RateService) getCommissionRate() (float64, error) {
	var setting models.SystemSetting
	
	err := s.db.Where("setting_key = ? AND is_active = ?", "commission_percent", true).
		First(&setting).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 2.0, nil // Дефолтная комиссия 2%
		}
		return 0, err
	}
	
	commission, err := strconv.ParseFloat(setting.SettingValue, 64)
	if err != nil {
		return 2.0, nil // Дефолтная комиссия при ошибке парсинга
	}
	
	return commission, nil
}