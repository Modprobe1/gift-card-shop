#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080/api/v1"

echo -e "${YELLOW}🧪 Тестирование API криптообменника${NC}"
echo "=================================="

# Функция для проверки ответа
check_response() {
    if [[ $1 -eq 200 || $1 -eq 201 ]]; then
        echo -e "${GREEN}✅ SUCCESS (HTTP $1)${NC}"
    else
        echo -e "${RED}❌ FAILED (HTTP $1)${NC}"
    fi
}

# Health check
echo -e "\n${YELLOW}1. Health Check${NC}"
echo "GET /health"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json http://localhost:8080/health)
check_response $response
cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json

# Получить валюты
echo -e "\n${YELLOW}2. Получить все валюты${NC}"
echo "GET $API_BASE/currencies"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json $API_BASE/currencies)
check_response $response
cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json

# Получить конкретную валюту
echo -e "\n${YELLOW}3. Получить валюту USDT_TRC20${NC}"
echo "GET $API_BASE/currencies/USDT_TRC20"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json $API_BASE/currencies/USDT_TRC20)
check_response $response
cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json

# Получить курсы
echo -e "\n${YELLOW}4. Получить все курсы обмена${NC}"
echo "GET $API_BASE/rates"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json $API_BASE/rates)
check_response $response
cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json

# Калькулятор
echo -e "\n${YELLOW}5. Калькулятор обмена${NC}"
echo "POST $API_BASE/calculate"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
  -X POST $API_BASE/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USDT_TRC20",
    "to_currency": "RUB_TBANK",
    "from_amount": 100
  }')
check_response $response
cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json

# Создать заявку
echo -e "\n${YELLOW}6. Создать заявку на обмен${NC}"
echo "POST $API_BASE/orders"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
  -X POST $API_BASE/orders \
  -H "Content-Type: application/json" \
  -d '{
    "from_currency": "USDT_TRC20",
    "to_currency": "RUB_TBANK",
    "from_amount": 100,
    "client_phone": "+79123456789",
    "client_name": "Тестовый Пользователь",
    "client_email": "test@example.com",
    "client_telegram": "@test_user",
    "recipient_wallet": "TRX7n4VkyYF7PgAwVUGdNdJg5dgrEpqhPH",
    "recipient_details": "Тестовая заявка"
  }')
check_response $response
ORDER_RESPONSE=$(cat /tmp/response.json)
echo "$ORDER_RESPONSE" | jq 2>/dev/null || echo "$ORDER_RESPONSE"

# Извлекаем номер заявки из ответа
ORDER_NUMBER=$(echo "$ORDER_RESPONSE" | jq -r '.data.order_number' 2>/dev/null)

if [[ "$ORDER_NUMBER" != "null" && "$ORDER_NUMBER" != "" ]]; then
    # Получить заявку по номеру
    echo -e "\n${YELLOW}7. Получить заявку по номеру${NC}"
    echo "GET $API_BASE/orders/$ORDER_NUMBER"
    response=$(curl -s -w "%{http_code}" -o /tmp/response.json $API_BASE/orders/$ORDER_NUMBER)
    check_response $response
    cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json
fi

# Админские эндпоинты
echo -e "\n${YELLOW}8. Получить статистику (админ)${NC}"
echo "GET $API_BASE/admin/statistics"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json $API_BASE/admin/statistics)
check_response $response
cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json

echo -e "\n${YELLOW}9. Получить список заявок (админ)${NC}"
echo "GET $API_BASE/admin/orders"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json $API_BASE/admin/orders)
check_response $response
cat /tmp/response.json | jq 2>/dev/null || cat /tmp/response.json

echo -e "\n${YELLOW}🎉 Тестирование завершено!${NC}"
echo "=================================="

# Очистка
rm -f /tmp/response.json