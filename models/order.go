package models

import "time"

type ProductStockUpdatedEvent struct {
	SchemaVersion int           `json:"schema_version"`
	OrderID       int64         `json:"order_id"`
	UserID        int64         `json:"user_id"`
	Products      []ProductItem `json:"products"`
	EventTime     time.Time     `json:"event_time"`
}

type OrderCreatedEvent struct {
	OrderID         int64         `json:"order_id"`
	UserID          int64         `json:"user_id"`
	TotalAmount     float64       `json:"total_amount"`
	PaymentMethod   string        `json:"payment_method"`
	ShippingAddress string        `json:"shipping_address"`
	Products        []ProductItem `json:"products"`
}

type ProductItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type ProductStockRollbackEvent struct {
	SchemaVersion int           `json:"schema_version"`
	OrderID       int64         `json:"order_id"`
	UserID        int64         `json:"user_id"`
	Products      []ProductItem `json:"products"`
	EventTime     time.Time     `json:"event_time"`
}

type StockReservationEvent struct {
	SchemaVersion int           `json:"schema_version"`
	OrderID       int64         `json:"order_id"`
	UserID        int64         `json:"user_id"`
	TotalAmount   float64       `json:"total_amount"`
	Products      []ProductItem `json:"products"`
	Reason        string        `json:"reason,omitempty"`
	EventTime     time.Time     `json:"event_time"`
}
