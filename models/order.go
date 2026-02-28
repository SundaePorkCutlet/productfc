package models

import "time"

type ProductStockUpdatedEvent struct {
	OrderID   int64         `json:"order_id"`
	Products  []ProductItem `json:"products"`
	EventTime time.Time     `json:"event_time"`
}

type ProductItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type ProductStockRollbackEvent struct {
	OrderID   int64         `json:"order_id"`
	Products  []ProductItem `json:"products"`
	EventTime time.Time     `json:"event_time"`
}
