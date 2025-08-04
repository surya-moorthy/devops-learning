package orders

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/devops-learning/ecommerce/internals/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderHandler struct {
	repo *OrderRepository
}

func NewOrderHandler(repo *OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// CreateOrderRequest represents the request payload for creating an order
type CreateOrderRequest struct {
	UserID          uint                      `json:"user_id" validate:"required"`
	ShippingAddress *string                   `json:"shipping_address,omitempty"`
	Notes           *string                   `json:"notes,omitempty"`
	OrderItems      []CreateOrderItemRequest  `json:"order_items" validate:"required,min=1"`
}

type CreateOrderItemRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  uint    `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=0"`
}

// UpdateOrderRequest represents the request payload for updating an order
type UpdateOrderRequest struct {
	ShippingAddress *string `json:"shipping_address,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

// UpdateOrderStatusRequest represents the request payload for updating order status
type UpdateOrderStatusRequest struct {
	Status models.OrderStatus `json:"status" validate:"required"`
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate order number (you might want to use a more sophisticated approach)
	orderNumber := fmt.Sprintf("ORD-%d-%d", req.UserID, time.Now().Unix())

	// Calculate total amount
	var totalAmount float64
	for _, item := range req.OrderItems {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Create order
	order := models.Order{
		UserID:          req.UserID,
		OrderNumber:     orderNumber,
		Status:          models.OrderStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: req.ShippingAddress,
		Notes:           req.Notes,
		OrderItems:      make([]models.OrderItem, len(req.OrderItems)),
	}

	// Add order items
	for i, item := range req.OrderItems {
		order.OrderItems[i] = models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}

	if err := h.repo.CreateOrder(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	// Fetch the created order with relations
	createdOrder, err := h.repo.GetOrderByID(order.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch created order",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Order created successfully",
		"order":   createdOrder,
	})
}

// GetOrder retrieves an order by ID
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := h.repo.GetOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch order",
		})
	}

	return c.JSON(fiber.Map{
		"order": order,
	})
}

// GetOrderByNumber retrieves an order by order number
func (h *OrderHandler) GetOrderByNumber(c *fiber.Ctx) error {
	orderNumber := c.Params("orderNumber")
	if orderNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order number is required",
		})
	}

	order, err := h.repo.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch order",
		})
	}

	return c.JSON(fiber.Map{
		"order": order,
	})
}

// GetUserOrders retrieves all orders for a specific user
func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userID"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	orders, total, err := h.repo.GetOrdersByUserID(uint(userID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAllOrders retrieves all orders with optional filtering
func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	// Parse status filter
	status := models.OrderStatus(c.Query("status", ""))

	orders, total, err := h.repo.GetAllOrders(limit, offset, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateOrder updates order information
func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req UpdateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if order exists
	existingOrder, err := h.repo.GetOrderByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch order",
		})
	}

	// Update order fields
	existingOrder.ShippingAddress = req.ShippingAddress
	existingOrder.Notes = req.Notes

	if err := h.repo.UpdateOrder(existingOrder); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order",
		})
	}

	// Fetch updated order
	updatedOrder, err := h.repo.GetOrderByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch updated order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order updated successfully",
		"order":   updatedOrder,
	})
}

// UpdateOrderStatus updates order status
func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	validStatuses := []models.OrderStatus{
		models.OrderStatusPending, models.OrderStatusConfirmed, models.OrderStatusShipped,
		models.OrderStatusDelivered, models.OrderStatusCancelled, models.OrderStatusRefunded,
	}
	
	isValid := false
	for _, status := range validStatuses {
		if req.Status == status {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order status",
		})
	}

	if err := h.repo.UpdateOrderStatus(uint(id), req.Status); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	// Fetch updated order
	updatedOrder, err := h.repo.GetOrderByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch updated order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
		"order":   updatedOrder,
	})
}

// DeleteOrder soft deletes an order
func (h *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	if err := h.repo.DeleteOrder(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete order",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order deleted successfully",
	})
}
