import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { exchangeApi, ExchangeOrder } from '../services/api';
import './OrderPage.css';

const OrderPage: React.FC = () => {
  const { orderNumber, orderId } = useParams<{ orderNumber: string; orderId: string }>();
  const [order, setOrder] = useState<ExchangeOrder | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    if (orderNumber) {
      loadOrder();
    }
  }, [orderNumber]);

  const loadOrder = async () => {
    if (!orderNumber) return;
    
    try {
      setLoading(true);
      const orderData = await exchangeApi.getOrder(orderNumber);
      setOrder(orderData);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Заявка не найдена');
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return '#ffc107';
      case 'confirmed': return '#17a2b8';
      case 'processing': return '#007bff';
      case 'completed': return '#28a745';
      case 'cancelled': return '#dc3545';
      case 'expired': return '#6c757d';
      default: return '#6c757d';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'pending': return 'Ожидает обработки';
      case 'confirmed': return 'Подтверждена';
      case 'processing': return 'В обработке';
      case 'completed': return 'Завершена';
      case 'cancelled': return 'Отменена';
      case 'expired': return 'Истекла';
      default: return status;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  if (loading) {
    return (
      <div className="order-page">
        <div className="order-container">
          <div className="loading-spinner">
            <div className="spinner"></div>
            <p>Загрузка заявки...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="order-page">
        <div className="order-container">
          <div className="error-container">
            <h2>❌ Ошибка</h2>
            <p>{error}</p>
            <button onClick={() => window.history.back()} className="back-button">
              Вернуться назад
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="order-page">
        <div className="order-container">
          <div className="error-container">
            <h2>Заявка не найдена</h2>
            <button onClick={() => window.history.back()} className="back-button">
              Вернуться назад
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="order-page">
      <div className="order-container">
        {/* Заголовок */}
        <div className="order-header">
          <h1>🎉 Заявка создана!</h1>
          <p className="order-subtitle">
            Ваша заявка успешно создана и отправлена на обработку
          </p>
        </div>

        {/* Основная информация о заявке */}
        <div className="order-info-card">
          <div className="order-number">
            <span className="label">Номер заявки:</span>
            <span className="value">{order.order_number}</span>
          </div>
          
          <div className="order-status">
            <span className="label">Статус:</span>
            <span 
              className="status-badge"
              style={{ backgroundColor: getStatusColor(order.status.name) }}
            >
              {getStatusText(order.status.name)}
            </span>
          </div>
          
          <div className="order-date">
            <span className="label">Создана:</span>
            <span className="value">{formatDate(order.created_at)}</span>
          </div>
        </div>

        {/* Детали обмена */}
        <div className="exchange-details">
          <h3>Детали обмена</h3>
          
          <div className="exchange-summary">
            <div className="exchange-row">
              <div className="exchange-from">
                <div className="currency-info">
                  <span className="currency-name">{order.from_currency.name}</span>
                  <span className="currency-network">({order.from_currency.network})</span>
                </div>
                <div className="amount">{order.from_amount} {order.from_currency.symbol}</div>
              </div>
              
              <div className="exchange-arrow">→</div>
              
              <div className="exchange-to">
                <div className="currency-info">
                  <span className="currency-name">{order.to_currency.name}</span>
                  <span className="currency-network">({order.to_currency.network})</span>
                </div>
                <div className="amount">{order.to_amount.toFixed(2)} {order.to_currency.symbol}</div>
              </div>
            </div>
            
            <div className="exchange-rate">
              <span>Курс: 1 {order.from_currency.symbol} = {order.rate.toFixed(4)} {order.to_currency.symbol}</span>
            </div>
          </div>
        </div>

        {/* Контактная информация */}
        <div className="contact-info">
          <h3>Контактная информация</h3>
          <div className="contact-grid">
            <div className="contact-item">
              <span className="contact-label">Имя:</span>
              <span className="contact-value">{order.client_name}</span>
            </div>
            <div className="contact-item">
              <span className="contact-label">Телефон:</span>
              <span className="contact-value">{order.client_phone}</span>
            </div>
            <div className="contact-item">
              <span className="contact-label">Email:</span>
              <span className="contact-value">{order.client_email}</span>
            </div>
            {order.client_telegram && (
              <div className="contact-item">
                <span className="contact-label">Telegram:</span>
                <span className="contact-value">{order.client_telegram}</span>
              </div>
            )}
          </div>
        </div>

        {/* Реквизиты */}
        <div className="wallet-info">
          <h3>Адрес для получения</h3>
          <div className="wallet-address">
            {order.recipient_wallet}
          </div>
          {order.recipient_details && (
            <div className="wallet-details">
              <strong>Дополнительная информация:</strong>
              <p>{order.recipient_details}</p>
            </div>
          )}
        </div>

        {/* Следующие шаги */}
        <div className="next-steps">
          <h3>Что дальше?</h3>
          <ol>
            <li>Мы получили вашу заявку и скоро свяжемся с вами</li>
            <li>После подтверждения данных заявка перейдет в обработку</li>
            <li>Средства будут переведены на указанный адрес</li>
            <li>Вы получите уведомление о завершении операции</li>
          </ol>
        </div>

        {/* Действия */}
        <div className="order-actions">
          <button 
            onClick={() => window.location.href = '/'}
            className="new-order-button"
          >
            Создать новую заявку
          </button>
          
          <button 
            onClick={() => navigator.clipboard.writeText(window.location.href)}
            className="copy-link-button"
          >
            Скопировать ссылку
          </button>
        </div>

        {/* Важная информация */}
        <div className="important-info">
          <h4>⚠️ Важно</h4>
          <ul>
            <li>Сохраните эту страницу или номер заявки для отслеживания</li>
            <li>Время обработки заявки: до 30 минут</li>
            <li>При возникновении вопросов обращайтесь в поддержку</li>
          </ul>
        </div>
      </div>
    </div>
  );
};

export default OrderPage;