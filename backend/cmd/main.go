package main

import (
  "context"
  "exchange/handlers"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "log"
  "net/http"
  "github.com/gin-gonic/gin"
)

func main() {
  client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo:27017"))
  handlers.OrderCollection = client.Database("exchange").Collection("orders")

  r := gin.Default()
  r.POST("/api/orders", handlers.CreateOrder)
  r.GET("/tx_:id", handlers.GetOrderByID)
  r.GET("/api/admin/orders", handlers.ListOrders)
  r.GET("/api/rate/usdt_rub", handlers.GetUSDTtoRUBRate)

  log.Println("http://localhost:8080")
  http.ListenAndServe(":8080", r)
}
