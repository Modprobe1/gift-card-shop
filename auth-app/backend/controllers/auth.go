package controllers

import (
	"auth-backend/config"
	"auth-backend/models"
	"auth-backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Register обрабатывает регистрацию нового пользователя
func Register(c *gin.Context) {
	var input models.RegisterInput

	// Валидация входных данных
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение коллекции пользователей
	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверка, существует ли пользователь с таким email
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Проверка, существует ли пользователь с таким username
	err = collection.FindOne(ctx, bson.M{"username": input.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this username already exists"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Создание нового пользователя
	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Username:  input.Username,
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохранение пользователя в базе данных
	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Генерация JWT токена
	token, err := utils.GenerateToken(newUser.ID.Hex(), newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user": gin.H{
			"id":         newUser.ID.Hex(),
			"username":   newUser.Username,
			"email":      newUser.Email,
			"first_name": newUser.FirstName,
			"last_name":  newUser.LastName,
		},
	})
}

// Login обрабатывает вход пользователя
func Login(c *gin.Context) {
	var input models.LoginInput

	// Валидация входных данных
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение коллекции пользователей
	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Поиск пользователя по email
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	// Проверка пароля
	if !utils.CheckPassword(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Генерация JWT токена
	token, err := utils.GenerateToken(user.ID.Hex(), user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":         user.ID.Hex(),
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}