package domain

import "time"

type OrderStatus string

const (
	PendingPayment OrderStatus = "PendingPayment"
	Paid           OrderStatus = "Paid"
	Fulfilling     OrderStatus = "Fulfilling"
	Shipped        OrderStatus = "Shipped"
	Completed      OrderStatus = "Completed"
	Cancelled      OrderStatus = "Cancelled"
	Refunded       OrderStatus = "Refunded"
)

type Order struct {
	ID           string      `json:"id" gorm:"primaryKey;size:36"`
	CustomerID   string      `json:"customerId" gorm:"size:100;not null"`
	ProductSKU   string      `json:"productSku" gorm:"size:100;not null"`
	Quantity     int         `json:"quantity"`
	Status       OrderStatus `json:"status" gorm:"size:50;not null"`
	Version      int         `json:"version"`
	CreatedAtUTC time.Time   `json:"createdAtUtc"`
	UpdatedAtUTC time.Time   `json:"updatedAtUtc"`
}

type OutboxEvent struct {
	ID             string     `json:"id" gorm:"primaryKey;size:36"`
	OrderID        string     `json:"orderId" gorm:"size:36;not null"`
	EventType      string     `json:"eventType" gorm:"size:100;not null"`
	PayloadJSON    string     `json:"payloadJson" gorm:"type:text;not null"`
	OccurredAtUTC  time.Time  `json:"occurredAtUtc"`
	Published      bool       `json:"published"`
	PublishedAtUTC *time.Time `json:"publishedAtUtc"`
}

type CreateOrderRequest struct {
	CustomerID string `json:"customerId"`
	ProductSKU string `json:"productSku"`
	Quantity   int    `json:"quantity"`
}

type ChangeStatusRequest struct {
	Action string  `json:"action"`
	Reason *string `json:"reason"`
}
