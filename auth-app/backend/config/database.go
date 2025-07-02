package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	// Получение URI MongoDB из переменной окружения или использование значения по умолчанию
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Настройка опций клиента
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	// Получение базы данных
	DB = client.Database("auth_app")
}

// GetCollection возвращает коллекцию MongoDB
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}