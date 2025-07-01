package models

import "time"

type Order struct {
  ID string `json:"id" bson:"_id"`
  From string `json:"from"`
  To string `json:"to"`
  Amount float64 `json:"amount"`
  Rate float64 `json:"rate"`
  ReceiveAmount float64 `json:"receive_amount"`
  Profit float64 `json:"profit"`
  FIO string `json:"fio"`
  Card string `json:"card"`
  IP string `json:"ip"`
  UserAgent string `json:"user_agent"`
  CreatedAt time.Time `json:"created_at"`
}

type OrderInput struct {
  From string `json:"from"`
  To string `json:"to"`
  Amount float64 `json:"amount"`
  Rate float64 `json:"rate"`
  ReceiveAmount float64 `json:"receive_amount"`
  FIO string `json:"fio"`
  Card string `json:"card"`
}
