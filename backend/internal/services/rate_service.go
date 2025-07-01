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

type CoinGeckoResponse struct {
	Bitcoin struct {
		USD float64 `json:"usd"`
		RUB float64 `json:"rub"`
	} `json:"bitcoin"`
	Ethereum struct {
		USD float64 `json:"usd"`
		RUB float64 `json:"rub"`
	} `json:"ethereum"`
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
	
	// Получаем комиссию из настроек
	commissionRate, err := s.getCommissionRate()
	if err != nil {
		return nil, err
	}
	
	// Вычисляем сумму без комиссии
	baseToAmount := fromAmount * exchangeRate.Rate
	
	// Вычисляем комиссию
	commission := baseToAmount * (commissionRate / 100)
	
	// Итоговая сумма к получению
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
	
	logrus.Info("Starting rate updater")
	
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

// UpdateRatesFromAPI обновляет курсы из внешнего API
func (s *RateService) UpdateRatesFromAPI() error {
	logrus.Debug("Updating exchange rates from CoinGecko API")
	
	// Запрос к CoinGecko API
	url := fmt.Sprintf("%s/simple/price?ids=bitcoin,ethereum,tether&vs_currencies=usd,rub", s.config.API.CoinGeckoURL)
	
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
	
	// Обновляем курсы
	rates := []struct {
		from, to string
		rate     float64
	}{
		// BTC курсы
		{"BTC", "USDT_TRC20", 1.0 / apiResponse.Bitcoin.USD}, // BTC -> USDT
		{"USDT_TRC20", "BTC", apiResponse.Bitcoin.USD},       // USDT -> BTC
		{"BTC", "RUB_TBANK", apiResponse.Bitcoin.RUB},        // BTC -> RUB
		{"RUB_TBANK", "BTC", 1.0 / apiResponse.Bitcoin.RUB},  // RUB -> BTC
		
		// ETH курсы
		{"ETH", "USDT_TRC20", 1.0 / apiResponse.Ethereum.USD}, // ETH -> USDT
		{"USDT_TRC20", "ETH", apiResponse.Ethereum.USD},       // USDT -> ETH
		{"ETH", "RUB_TBANK", apiResponse.Ethereum.RUB},        // ETH -> RUB
		{"RUB_TBANK", "ETH", 1.0 / apiResponse.Ethereum.RUB},  // RUB -> ETH
		
		// USDT курсы
		{"USDT_TRC20", "RUB_TBANK", apiResponse.Tether.RUB},       // USDT -> RUB
		{"RUB_TBANK", "USDT_TRC20", 1.0 / apiResponse.Tether.RUB}, // RUB -> USDT
	}
	
	for _, r := range rates {
		fromID, fromExists := currencies[r.from]
		toID, toExists := currencies[r.to]
		
		if fromExists && toExists {
			if err := s.UpdateExchangeRate(fromID, toID, r.rate, "coingecko"); err != nil {
				logrus.WithError(err).Errorf("Failed to update rate %s -> %s", r.from, r.to)
			}
		}
	}
	
	logrus.Debug("Successfully updated exchange rates")
	return nil
}

// getCommissionRate получает процент комиссии из настроек
func (s *RateService) getCommissionRate() (float64, error) {
	var setting models.SystemSetting
	
	err := s.db.Where("setting_key = ? AND is_active = ?", "commission_percent", true).
		First(&setting).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1.5, nil // Дефолтная комиссия 1.5%
		}
		return 0, err
	}
	
	commission, err := strconv.ParseFloat(setting.SettingValue, 64)
	if err != nil {
		return 1.5, nil // Дефолтная комиссия при ошибке парсинга
	}
	
	return commission, nil
}