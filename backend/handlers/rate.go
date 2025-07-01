package handlers

import (
  "encoding/json"
  "fmt"
  "net/http"
  "sync"
  "time"

  "github.com/gin-gonic/gin"
)

var (
  cachedRate float64
  rateMutex  sync.RWMutex
)

func init() {
  go updateRatePeriodically()
  go updateRateOnce()
}

func updateRatePeriodically() {
  for {
    fetchAndUpdateRate()
    time.Sleep(60 * time.Second)
  }
}

func updateRateOnce() {
  fetchAndUpdateRate()
}

func fetchAndUpdateRate() {
  resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?symbol=USDTRUB")
  if err != nil {
    fmt.Println("Ошибка запроса Binance:", err)
    return
  }
  defer resp.Body.Close()

  var data map[string]string
  if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
    fmt.Println("Ошибка разбора ответа:", err)
    return
  }

  priceStr := data["price"]
  var price float64
  if _, err := fmt.Sscanf(priceStr, "%f", &price); err != nil {
    fmt.Println("Ошибка конвертации курса:", err)
    return
  }

  rateMutex.Lock()
  cachedRate = price * 0.97 // 3% наценка
  rateMutex.Unlock()
}

func GetUSDTtoRUBRate(c *gin.Context) {
  rateMutex.RLock()
  rate := cachedRate
  rateMutex.RUnlock()
  if rate == 0 {
    c.JSON(http.StatusServiceUnavailable, gin.H{"error": "rate not available yet"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"rate": rate})
}
