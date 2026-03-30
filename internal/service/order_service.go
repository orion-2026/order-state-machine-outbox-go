package service

import (
	"encoding/json"
	"fmt"
	"order-state-machine-outbox-go/internal/domain"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderService struct {
	db           *gorm.DB
	stateMachine *StateMachine
}

func NewOrderService(db *gorm.DB, stateMachine *StateMachine) *OrderService {
	return &OrderService{db: db, stateMachine: stateMachine}
}

func (s *OrderService) Create(req domain.CreateOrderRequest) (*domain.Order, error) {
	now := time.Now().UTC()
	order := &domain.Order{
		ID:           uuid.NewString(),
		CustomerID:   req.CustomerID,
		ProductSKU:   req.ProductSKU,
		Quantity:     req.Quantity,
		Status:       domain.PendingPayment,
		Version:      1,
		CreatedAtUTC: now,
		UpdatedAtUTC: now,
	}

	return order, s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return appendOutbox(tx, order.ID, "OrderCreated", map[string]any{
			"id":         order.ID,
			"customerId": order.CustomerID,
			"productSku": order.ProductSKU,
			"quantity":   order.Quantity,
			"status":     order.Status,
			"version":    order.Version,
		})
	})
}

func (s *OrderService) ListOrders() ([]domain.Order, error) {
	var orders []domain.Order
	if err := s.db.Find(&orders).Error; err != nil {
		return nil, err
	}
	sort.Slice(orders, func(i, j int) bool { return orders[i].CreatedAtUTC.Before(orders[j].CreatedAtUTC) })
	return orders, nil
}

func (s *OrderService) Get(id string) (*domain.Order, error) {
	var order domain.Order
	if err := s.db.First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) AllowedActions(id string) ([]string, error) {
	order, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	actions := s.stateMachine.AllowedActions(order.Status)
	sort.Strings(actions)
	return actions, nil
}

func (s *OrderService) ChangeStatus(id string, req domain.ChangeStatusRequest) (*domain.Order, error) {
	var order domain.Order
	return &order, s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&order, "id = ?", id).Error; err != nil {
			return err
		}
		action := strings.TrimSpace(strings.ToLower(req.Action))
		nextStatus, ok := s.stateMachine.TryTransition(order.Status, action)
		if !ok {
			allowed := strings.Join(s.stateMachine.AllowedActions(order.Status), ", ")
			return fmt.Errorf("invalid transition from %s using action '%s'. allowed actions: %s", order.Status, req.Action, allowed)
		}
		previous := order.Status
		order.Status = nextStatus
		order.Version++
		order.UpdatedAtUTC = time.Now().UTC()
		if err := tx.Save(&order).Error; err != nil {
			return err
		}
		return appendOutbox(tx, order.ID, "OrderStatusChanged", map[string]any{
			"id":             order.ID,
			"previousStatus": previous,
			"newStatus":      nextStatus,
			"action":         action,
			"reason":         req.Reason,
			"version":        order.Version,
		})
	})
}

func (s *OrderService) ListOutbox() ([]domain.OutboxEvent, error) {
	var events []domain.OutboxEvent
	if err := s.db.Find(&events).Error; err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].OccurredAtUTC.Before(events[j].OccurredAtUTC) })
	return events, nil
}

func (s *OrderService) PublishPendingOutbox() ([]domain.OutboxEvent, error) {
	var events []domain.OutboxEvent
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("published = ?", false).Find(&events).Error; err != nil {
			return err
		}
		now := time.Now().UTC()
		for i := range events {
			events[i].Published = true
			events[i].PublishedAtUTC = &now
			if err := tx.Save(&events[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(events, func(i, j int) bool { return events[i].OccurredAtUTC.Before(events[j].OccurredAtUTC) })
	return events, nil
}

func appendOutbox(tx *gorm.DB, orderID, eventType string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	event := domain.OutboxEvent{
		ID:            uuid.NewString(),
		OrderID:       orderID,
		EventType:     eventType,
		PayloadJSON:   string(body),
		OccurredAtUTC: time.Now().UTC(),
		Published:     false,
	}
	return tx.Create(&event).Error
}
