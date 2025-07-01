package services

import (
	"crypto-exchange-backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

type CurrencyService struct {
	db *gorm.DB
}

func NewCurrencyService(db *gorm.DB) *CurrencyService {
	return &CurrencyService{db: db}
}

// GetAllCurrencies возвращает все активные валюты
func (s *CurrencyService) GetAllCurrencies() ([]models.Currency, error) {
	var currencies []models.Currency
	
	if err := s.db.Where("is_active = ?", true).Find(&currencies).Error; err != nil {
		return nil, err
	}
	
	return currencies, nil
}

// GetCurrencyByCode возвращает валюту по коду
func (s *CurrencyService) GetCurrencyByCode(code string) (*models.Currency, error) {
	var currency models.Currency
	
	if err := s.db.Where("code = ? AND is_active = ?", code, true).First(&currency).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("currency not found")
		}
		return nil, err
	}
	
	return &currency, nil
}

// GetCurrencyByID возвращает валюту по ID
func (s *CurrencyService) GetCurrencyByID(id uint) (*models.Currency, error) {
	var currency models.Currency
	
	if err := s.db.Where("id = ? AND is_active = ?", id, true).First(&currency).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("currency not found")
		}
		return nil, err
	}
	
	return &currency, nil
}

// CreateCurrency создает новую валюту
func (s *CurrencyService) CreateCurrency(currency *models.Currency) error {
	return s.db.Create(currency).Error
}

// UpdateCurrency обновляет валюту
func (s *CurrencyService) UpdateCurrency(currency *models.Currency) error {
	return s.db.Save(currency).Error
}

// DeleteCurrency удаляет валюту (мягкое удаление)
func (s *CurrencyService) DeleteCurrency(id uint) error {
	return s.db.Delete(&models.Currency{}, id).Error
}