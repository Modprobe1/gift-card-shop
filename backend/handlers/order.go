package handlers

import (
  "context"
  "exchange/models"
  "github.com/gin-gonic/gin"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "net/http"
  "time"
  "github.com/google/uuid"
)

var OrderCollection *mongo.Collection

func CreateOrder(c *gin.Context) {
  var input models.OrderInput
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
    return
  }
  profit := input.Amount * input.Rate * 0.03
  id := uuid.New().String()
  order := models.Order{
    ID: id, From: input.From, To: input.To,
    Amount: input.Amount, Rate: input.Rate,
    ReceiveAmount: input.ReceiveAmount, Profit: profit,
    FIO: input.FIO, Card: input.Card,
    IP: c.ClientIP(), UserAgent: c.Request.UserAgent(),
    CreatedAt: time.Now(),
  }
  OrderCollection.InsertOne(context.TODO(), order)
  c.JSON(http.StatusOK, gin.H{"order_id": id})
}

func GetOrderByID(c *gin.Context) {
  id := c.Param("id")
  var order models.Order
  err := OrderCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&order)
  if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
    return
  }
  c.JSON(http.StatusOK, order)
}

func ListOrders(c *gin.Context) {
  cursor, _ := OrderCollection.Find(context.TODO(), bson.M{})
  var orders []models.Order
  cursor.All(context.TODO(), &orders)
  c.JSON(http.StatusOK, orders)
}
