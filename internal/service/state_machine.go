package service

import "order-state-machine-outbox-go/internal/domain"

type StateMachine struct {
	transitions map[domain.OrderStatus]map[string]domain.OrderStatus
}

func NewStateMachine() *StateMachine {
	return &StateMachine{transitions: map[domain.OrderStatus]map[string]domain.OrderStatus{
		domain.PendingPayment: {"pay": domain.Paid, "cancel": domain.Cancelled},
		domain.Paid:           {"start-fulfillment": domain.Fulfilling, "refund": domain.Refunded},
		domain.Fulfilling:     {"ship": domain.Shipped, "cancel": domain.Cancelled},
		domain.Shipped:        {"complete": domain.Completed, "refund": domain.Refunded},
	}}
}

func (s *StateMachine) TryTransition(current domain.OrderStatus, action string) (domain.OrderStatus, bool) {
	next, ok := s.transitions[current][action]
	return next, ok
}

func (s *StateMachine) AllowedActions(current domain.OrderStatus) []string {
	actions := make([]string, 0)
	for action := range s.transitions[current] {
		actions = append(actions, action)
	}
	return actions
}
