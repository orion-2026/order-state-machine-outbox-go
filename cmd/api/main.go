package main

import (
	"log"
	"order-state-machine-outbox-go/internal/config"
	httpapi "order-state-machine-outbox-go/internal/http"
	"order-state-machine-outbox-go/internal/repository"
	"order-state-machine-outbox-go/internal/service"
)

func main() {
	cfg := config.Load()
	db, err := repository.Open(cfg)
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}

	stateMachine := service.NewStateMachine()
	orderService := service.NewOrderService(db, stateMachine)
	router := httpapi.NewRouter(orderService)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
