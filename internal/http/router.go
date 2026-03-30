package http

import (
	"errors"
	"net/http"
	"order-state-machine-outbox-go/internal/domain"
	"order-state-machine-outbox-go/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(orderService *service.OrderService) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "Order State Machine + Outbox Demo (Go)",
			"pattern": "State machine + service + outbox",
		})
	})

	r.GET("/api/orders", func(c *gin.Context) {
		orders, err := orderService.ListOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	})

	r.GET("/api/orders/:id", func(c *gin.Context) {
		order, err := orderService.Get(c.Param("id"))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"message": "Order not found."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.POST("/api/orders", func(c *gin.Context) {
		var req domain.CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
			return
		}
		if strings.TrimSpace(req.CustomerID) == "" || strings.TrimSpace(req.ProductSKU) == "" || req.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "CustomerId, ProductSku, and positive Quantity are required."})
			return
		}
		order, err := orderService.Create(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, order)
	})

	r.GET("/api/orders/:id/actions", func(c *gin.Context) {
		actions, err := orderService.AllowedActions(c.Param("id"))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"message": "Order not found."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"actions": actions})
	})

	r.POST("/api/orders/:id/transitions", func(c *gin.Context) {
		var req domain.ChangeStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Action) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Action is required."})
			return
		}
		order, err := orderService.ChangeStatus(c.Param("id"), req)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"message": "Order not found."})
				return
			}
			if strings.Contains(strings.ToLower(err.Error()), "invalid transition") {
				c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.GET("/api/outbox", func(c *gin.Context) {
		events, err := orderService.ListOutbox()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, events)
	})

	r.POST("/api/outbox/publish", func(c *gin.Context) {
		events, err := orderService.PublishPendingOutbox()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"publishedCount": len(events), "events": events})
	})

	return r
}
