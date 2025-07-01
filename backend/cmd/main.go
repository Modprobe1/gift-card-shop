package main

import (
	"crypto-exchange-backend/internal/config"
	"crypto-exchange-backend/internal/database"
	"crypto-exchange-backend/internal/handlers"
	"crypto-exchange-backend/internal/services"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Настраиваем логирование
	setupLogging(cfg.Server.LogLevel)
	
	logrus.Info("Starting Crypto Exchange Backend...")
	logrus.Infof("Server Mode: %s", cfg.Server.Mode)
	logrus.Infof("Database Host: %s:%s", cfg.Database.Host, cfg.Database.Port)

	// Устанавливаем режим Gin
	gin.SetMode(cfg.Server.Mode)

	// Ждем немного, чтобы база данных успела запуститься
	time.Sleep(5 * time.Second)

	// Подключаемся к базе данных с повторными попытками
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = database.Connect(cfg)
		if err == nil {
			logrus.Info("Successfully connected to database")
			break
		}
		logrus.Warnf("Failed to connect to database (attempt %d/10): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database after 10 attempts:", err)
	}

	// Выполняем миграции
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Инициализируем сервисы
	currencyService := services.NewCurrencyService(db)
	rateService := services.NewRateService(db, cfg)
	orderService := services.NewOrderService(db, currencyService, rateService)

	// Запускаем обновление курсов в фоне
	go rateService.StartRateUpdater()

	// Инициализируем обработчики
	handler := handlers.NewHandler(currencyService, rateService, orderService)

	// Создаем роутер
	router := setupRouter(handler, cfg)

	// Запускаем сервер
	logrus.Infof("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupLogging(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func setupRouter(h *handlers.Handler, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS настройки
	corsOrigins := []string{"http://localhost:3000", "http://localhost:3001"}
	if cfg.Server.Mode == gin.ReleaseMode {
		// В продакшене разрешаем все origins (или можно указать конкретные)
		corsOrigins = []string{"*"}
	}
	
	router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false, // Устанавливаем false при использовании "*"
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})
	
	// Health check для API
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// API группа
	api := router.Group("/api/v1")
	{
		// Валюты
		api.GET("/currencies", h.GetCurrencies)
		api.GET("/currencies/:code", h.GetCurrency)

		// Курсы обмена
		api.GET("/rates", h.GetExchangeRates)
		api.GET("/rates/:from/:to", h.GetExchangeRate)

		// Калькулятор
		api.POST("/calculate", h.Calculate)

		// Заявки
		api.POST("/orders", h.CreateOrder)
		api.GET("/orders/:number", h.GetOrder)
		
		// Админские эндпоинты (в будущем добавим авторизацию)
		admin := api.Group("/admin")
		{
			admin.GET("/orders", h.GetOrders)
			admin.PUT("/orders/:id/status", h.UpdateOrderStatus)
			admin.GET("/statistics", h.GetStatistics)
		}
	}
	
	// Добавляем поддержку старых путей /api без версии
	apiOld := router.Group("/api")
	{
		apiOld.GET("/currencies", h.GetCurrencies)
		apiOld.GET("/rates", h.GetExchangeRates)
		apiOld.POST("/orders", h.CreateOrder)
		apiOld.GET("/orders/:number", h.GetOrder)
	}

	return router
}