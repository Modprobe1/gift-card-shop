package controllers

import (
	"auth-backend/config"
	"auth-backend/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetProfile возвращает профиль текущего пользователя
func GetProfile(c *gin.Context) {
	// Получение ID пользователя из контекста (установлен middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Преобразование ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получение коллекции пользователей
	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Поиск пользователя
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         user.ID.Hex(),
			"username":   user.Username,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}

// UpdateProfile обновляет профиль текущего пользователя
func UpdateProfile(c *gin.Context) {
	// Получение ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	// Структура для обновления профиля
	var updateInput struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username" binding:"omitempty,min=3,max=30"`
	}

	// Валидация входных данных
	if err := c.ShouldBindJSON(&updateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразование ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получение коллекции пользователей
	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подготовка обновления
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Добавление полей для обновления
	if updateInput.FirstName != "" {
		update["$set"].(bson.M)["first_name"] = updateInput.FirstName
	}
	if updateInput.LastName != "" {
		update["$set"].(bson.M)["last_name"] = updateInput.LastName
	}
	if updateInput.Username != "" {
		// Проверка уникальности username
		var existingUser models.User
		err := collection.FindOne(ctx, bson.M{
			"username": updateInput.Username,
			"_id":      bson.M{"$ne": objectID},
		}).Decode(&existingUser)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}
		update["$set"].(bson.M)["username"] = updateInput.Username
	}

	// Выполнение обновления
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Получение обновленного пользователя
	var updatedUser models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":         updatedUser.ID.Hex(),
			"username":   updatedUser.Username,
			"email":      updatedUser.Email,
			"first_name": updatedUser.FirstName,
			"last_name":  updatedUser.LastName,
			"updated_at": updatedUser.UpdatedAt,
		},
	})
}