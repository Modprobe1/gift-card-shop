import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { exchangeApi, Currency, CalculateResponse, CreateOrderRequest } from '../services/api';
import './ExchangeForm.css';

// Доступные направления обмена
const EXCHANGE_DIRECTIONS = [
  { from: 'USDT_TRC20', to: 'BTC', label: 'USDT → BTC' },
  { from: 'BTC', to: 'USDT_TRC20', label: 'BTC → USDT' },
  { from: 'USDT_TRC20', to: 'RUB_TBANK', label: 'USDT → RUB' },
  { from: 'RUB_TBANK', to: 'USDT_TRC20', label: 'RUB → USDT' },
];

const ExchangeForm: React.FC = () => {
  const navigate = useNavigate();
  
  // Состояния
  const [currencies, setCurrencies] = useState<Currency[]>([]);
  const [selectedDirection, setSelectedDirection] = useState(0); // Индекс выбранного направления
  const [fromAmount, setFromAmount] = useState<string>('');
  const [calculation, setCalculation] = useState<CalculateResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  // Данные формы
  const [formData, setFormData] = useState({
    clientName: '',
    clientPhone: '',
    clientEmail: '',
    clientTelegram: '',
    recipientWallet: '',
    recipientDetails: '',
  });

  // Загрузка валют при монтировании
  useEffect(() => {
    loadCurrencies();
  }, []);

  // Пересчет при изменении направления или суммы
  useEffect(() => {
    if (fromAmount && parseFloat(fromAmount) > 0) {
      calculateExchange();
    } else {
      setCalculation(null);
    }
  }, [selectedDirection, fromAmount]);

  const loadCurrencies = async () => {
    try {
      const data = await exchangeApi.getCurrencies();
      setCurrencies(data);
    } catch (err) {
      setError('Ошибка загрузки валют');
    }
  };

  const calculateExchange = async () => {
    if (!fromAmount || parseFloat(fromAmount) <= 0) return;
    
    const direction = EXCHANGE_DIRECTIONS[selectedDirection];
    
    try {
      setLoading(true);
      const result = await exchangeApi.calculate({
        from_currency: direction.from,
        to_currency: direction.to,
        from_amount: parseFloat(fromAmount),
      });
      setCalculation(result);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка расчета');
      setCalculation(null);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!calculation) {
      setError('Сначала выполните расчет обмена');
      return;
    }

    if (!formData.clientName || !formData.clientPhone || !formData.clientEmail || !formData.recipientWallet) {
      setError('Заполните все обязательные поля');
      return;
    }

    try {
      setLoading(true);
      
      const direction = EXCHANGE_DIRECTIONS[selectedDirection];
      const orderRequest: CreateOrderRequest = {
        from_currency: direction.from,
        to_currency: direction.to,
        from_amount: parseFloat(fromAmount),
        client_name: formData.clientName,
        client_phone: formData.clientPhone,
        client_email: formData.clientEmail,
        client_telegram: formData.clientTelegram,
        recipient_wallet: formData.recipientWallet,
        recipient_details: formData.recipientDetails,
      };

      const order = await exchangeApi.createOrder(orderRequest);
      
      // Перенаправляем на страницу заявки
      navigate(`/order/${order.order_number}/${order.id}`);
      
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка создания заявки');
    } finally {
      setLoading(false);
    }
  };

  const getCurrencyInfo = (code: string) => {
    return currencies.find(c => c.code === code);
  };

  const getMinAmount = () => {
    const direction = EXCHANGE_DIRECTIONS[selectedDirection];
    const currency = getCurrencyInfo(direction.from);
    return currency?.min_amount || 0;
  };

  const getCurrentDirection = () => EXCHANGE_DIRECTIONS[selectedDirection];

  return (
    <div className="exchange-container">
      <div className="exchange-header">
        <h1>Обмен криптовалюты</h1>
        {calculation && (
          <div className="current-rate">
            Курс: 1 {getCurrencyInfo(getCurrentDirection().from)?.symbol} = {calculation.rate.toFixed(6)} {getCurrencyInfo(getCurrentDirection().to)?.symbol}
          </div>
        )}
      </div>

      <div className="exchange-form-container">
        {/* Кнопки выбора направления */}
        <div className="direction-buttons">
          {EXCHANGE_DIRECTIONS.map((direction, index) => (
            <button
              key={index}
              type="button"
              className={`direction-btn ${selectedDirection === index ? 'active' : ''}`}
              onClick={() => setSelectedDirection(index)}
            >
              {direction.label}
            </button>
          ))}
        </div>

        <form onSubmit={handleSubmit} className="exchange-form">
          
          {/* Блок "Отдаете" */}
          <div className="exchange-block from-block">
            <label className="exchange-label">
              Отдаете ({getCurrencyInfo(getCurrentDirection().from)?.symbol}):
            </label>
            <div className="exchange-input-container">
              <input
                type="number"
                value={fromAmount}
                onChange={(e) => setFromAmount(e.target.value)}
                placeholder="Введите сумму"
                className="amount-input"
                min={getMinAmount()}
                step="any"
                required
              />
              <div className="currency-display">
                {getCurrencyInfo(getCurrentDirection().from)?.name}
              </div>
            </div>
            {getMinAmount() > 0 && (
              <div className="min-amount">Минимум: {getMinAmount()} {getCurrencyInfo(getCurrentDirection().from)?.symbol}</div>
            )}
          </div>

          {/* Блок "Получите" */}
          <div className="exchange-block to-block">
            <label className="exchange-label">
              Получите ({getCurrencyInfo(getCurrentDirection().to)?.symbol}):
            </label>
            <div className="exchange-input-container">
              <div className="amount-display">
                {loading ? 'Рассчитываем...' : (calculation ? calculation.to_amount.toFixed(6) : '0')}
              </div>
              <div className="currency-display">
                {getCurrencyInfo(getCurrentDirection().to)?.name}
              </div>
            </div>
            {calculation && (
              <div className="commission-info">
                Комиссия {calculation.commission_rate}%: {calculation.commission.toFixed(6)} {getCurrencyInfo(getCurrentDirection().to)?.symbol}
              </div>
            )}
          </div>

          {/* Поля контактной информации */}
          <div className="contact-fields">
            <h3>Контактная информация</h3>
            
            <input
              type="text"
              placeholder="ФИО *"
              value={formData.clientName}
              onChange={(e) => setFormData({...formData, clientName: e.target.value})}
              className="contact-input"
              required
            />

            <input
              type="tel"
              placeholder="Номер телефона *"
              value={formData.clientPhone}
              onChange={(e) => setFormData({...formData, clientPhone: e.target.value})}
              className="contact-input"
              required
            />

            <input
              type="email"
              placeholder="E-mail *"
              value={formData.clientEmail}
              onChange={(e) => setFormData({...formData, clientEmail: e.target.value})}
              className="contact-input"
              required
            />

            <input
              type="text"
              placeholder="Telegram"
              value={formData.clientTelegram}
              onChange={(e) => setFormData({...formData, clientTelegram: e.target.value})}
              className="contact-input"
            />

            <input
              type="text"
              placeholder={`Адрес кошелька для получения ${getCurrencyInfo(getCurrentDirection().to)?.symbol} *`}
              value={formData.recipientWallet}
              onChange={(e) => setFormData({...formData, recipientWallet: e.target.value})}
              className="contact-input"
              required
            />

            <textarea
              placeholder="Дополнительная информация"
              value={formData.recipientDetails}
              onChange={(e) => setFormData({...formData, recipientDetails: e.target.value})}
              className="contact-textarea"
              rows={3}
            />
          </div>

          {/* Согласие с условиями */}
          <div className="agreement">
            <label className="agreement-checkbox">
              <input type="checkbox" required />
              <span className="checkmark"></span>
              Я соглашаюсь с <button type="button" className="agreement-link">условиями сервиса</button>
            </label>
          </div>

          {/* Ошибки */}
          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          {/* Кнопка отправки */}
          <button
            type="submit"
            disabled={loading || !calculation || !fromAmount}
            className="submit-button"
          >
            {loading ? 'Создание заявки...' : 'Обменять'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default ExchangeForm;