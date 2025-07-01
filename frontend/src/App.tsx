import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ExchangeForm from './components/ExchangeForm';
import OrderPage from './components/OrderPage';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <div className="header-content">
            <h1 className="site-title">🔄 CryptoExchange</h1>
            <p className="site-subtitle">Быстрый и безопасный обмен криптовалют</p>
          </div>
        </header>

        <main className="App-main">
          <Routes>
            <Route path="/" element={<ExchangeForm />} />
            <Route path="/order/:orderNumber/:orderId" element={<OrderPage />} />
          </Routes>
        </main>

        <footer className="App-footer">
          <div className="footer-content">
            <p>&copy; 2024 CryptoExchange. Все права защищены.</p>
            <div className="footer-links">
              <a href="#support">Поддержка</a>
              <a href="#terms">Условия</a>
              <a href="#privacy">Конфиденциальность</a>
            </div>
          </div>
        </footer>
      </div>
    </Router>
  );
}

export default App;
