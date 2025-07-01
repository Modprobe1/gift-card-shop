import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { exchangeApi, Currency, CalculateResponse, CreateOrderRequest } from '../services/api';
import './ExchangeForm.css';

const ExchangeForm: React.FC = () => {
  const navigate = useNavigate();
  
  // Состояния
  const [currencies, setCurrencies] = useState<Currency[]>([]);
  const [fromCurrency, setFromCurrency] = useState<string>('USDT_TRC20');
  const [toCurrency, setToCurrency] = useState<string>('RUB_TBANK');
  const [fromAmount, setFromAmount] = useState<string>('100');
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

  // Пересчет при изменении валют или суммы
  useEffect(() => {
    if (fromAmount && parseFloat(fromAmount) > 0) {
      calculateExchange();
    }
  }, [fromCurrency, toCurrency, fromAmount]);

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
    
    try {
      setLoading(true);
      const result = await exchangeApi.calculate({
        from_currency: fromCurrency,
        to_currency: toCurrency,
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

  const handleSwapCurrencies = () => {
    const temp = fromCurrency;
    setFromCurrency(toCurrency);
    setToCurrency(temp);
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
      
      const orderRequest: CreateOrderRequest = {
        from_currency: fromCurrency,
        to_currency: toCurrency,
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

  const getCurrencyOptions = () => {
    return currencies.map(currency => (
      <option key={currency.code} value={currency.code}>
        {currency.name} ({currency.symbol})
      </option>
    ));
  };

  const getMinAmount = () => {
    const currency = currencies.find(c => c.code === fromCurrency);
    return currency?.min_amount || 0;
  };

  return (
    <div className="exchange-container">
      {/* Заголовок с курсом */}
      <div className="exchange-header">
        <div className="rate-info">
          {calculation && (
            <div className="current-rate">
              1.00 {currencies.find(c => c.code === fromCurrency)?.symbol} = {calculation.rate.toFixed(4)} {currencies.find(c => c.code === toCurrency)?.symbol}
            </div>
          )}
        </div>
      </div>

      {/* Основная форма обмена */}
      <div className="exchange-form-container">
        <form onSubmit={handleSubmit} className="exchange-form">
          
          {/* Блок "Вы отдаете" */}
          <div className="exchange-block from-block">
            <label className="exchange-label">Вы отдаете</label>
            <div className="exchange-input-container">
              <input
                type="number"
                value={fromAmount}
                onChange={(e) => setFromAmount(e.target.value)}
                placeholder="0"
                className="amount-input"
                min={getMinAmount()}
                step="0.01"
                required
              />
              <select 
                value={fromCurrency}
                onChange={(e) => setFromCurrency(e.target.value)}
                className="currency-select"
              >
                {getCurrencyOptions()}
              </select>
            </div>
            {getMinAmount() > 0 && (
              <div className="min-amount">Min: {getMinAmount()}</div>
            )}
          </div>

          {/* Кнопка обмена */}
          <div className="swap-button-container">
            <button
              type="button"
              onClick={handleSwapCurrencies}
              className="swap-button"
              title="Поменять валюты местами"
            >
              ⇅
            </button>
          </div>

          {/* Блок "Вы получаете" */}
          <div className="exchange-block to-block">
            <label className="exchange-label">Вы получаете</label>
            <div className="exchange-input-container">
              <div className="amount-display">
                {calculation ? calculation.to_amount.toFixed(2) : '0'}
              </div>
              <select 
                value={toCurrency}
                onChange={(e) => setToCurrency(e.target.value)}
                className="currency-select"
              >
                {getCurrencyOptions()}
              </select>
            </div>
          </div>

          {/* Информация о комиссии */}
          {calculation && (
            <div className="commission-info">
              <div className="commission-row">
                <span>Комиссия ({calculation.commission_rate}%):</span>
                <span>{calculation.commission.toFixed(2)} {currencies.find(c => c.code === toCurrency)?.symbol}</span>
              </div>
            </div>
          )}

          {/* Поля контактной информации */}
          <div className="contact-fields">
            <div className="contact-section">
              <h3>Контактная информация</h3>
              
              <input
                type="tel"
                placeholder="Ваш номер телефона для СБП *"
                value={formData.clientPhone}
                onChange={(e) => setFormData({...formData, clientPhone: e.target.value})}
                className="contact-input"
                required
              />

              <input
                type="text"
                placeholder="Ваш ФИО для удобства перевода *"
                value={formData.clientName}
                onChange={(e) => setFormData({...formData, clientName: e.target.value})}
                className="contact-input"
                required
              />

              <input
                type="email"
                placeholder="Ваш E-mail *"
                value={formData.clientEmail}
                onChange={(e) => setFormData({...formData, clientEmail: e.target.value})}
                className="contact-input"
                required
              />

              <div className="contact-row">
                <select className="contact-method">
                  <option value="telegram">Telegram</option>
                  <option value="whatsapp">WhatsApp</option>
                  <option value="viber">Viber</option>
                </select>
                <input
                  type="text"
                  placeholder="Никнейм"
                  value={formData.clientTelegram}
                  onChange={(e) => setFormData({...formData, clientTelegram: e.target.value})}
                  className="contact-input nickname-input"
                />
              </div>

              <input
                type="text"
                placeholder="Адрес кошелька для получения *"
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
          </div>

          {/* Согласие с условиями */}
          <div className="agreement">
            <label className="agreement-checkbox">
              <input type="checkbox" required />
              <span className="checkmark"></span>
              Я соглашаюсь с <a href="#" className="agreement-link">условиями сервиса</a> и <a href="#" className="agreement-link">базовой проверкой аккаунта AML и KYC</a>
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
            disabled={loading || !calculation}
            className="submit-button"
          >
            {loading ? 'Создание заявки...' : 'Обменять сейчас'}
          </button>
        </form>
      </div>
    </div>
  );
};

export default ExchangeForm;