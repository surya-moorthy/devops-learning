package cart

import (
	"errors"
	"strconv"
	"time"

	"github.com/devops-learning/ecommerce/internals/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DTO structures for request/response
type CreateCartItemRequest struct {
	CartID    uint `json:"cart_id" validate:"required"`
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  uint `json:"quantity" validate:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity uint `json:"quantity" validate:"required,min=1"`
}

type CartItemResponse struct {
	ID        uint                 `json:"id"`
	CartID    uint                 `json:"cart_id"`
	ProductID uint                 `json:"product_id"`
	Quantity  uint                 `json:"quantity"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
	Product   *models.Product      `json:"product,omitempty"`
}

type CartItemHandler struct {
	cartItemRepo cartItemRepository
}

func NewCartItemHandler(cartItemRepo cartItemRepository) *CartItemHandler {
	return &CartItemHandler{
		cartItemRepo: cartItemRepo,
	}
}

// CreateCartItem creates a new cart item or updates existing one
func (h *CartItemHandler) CreateCartItem(c *fiber.Ctx) error {
	var req CreateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if item already exists in cart
	existingItem, err := h.cartItemRepo.GetByCartAndProduct(req.CartID, req.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check existing cart item",
		})
	}

	// If item exists, update quantity
	if existingItem != nil {
		existingItem.Quantity += req.Quantity
		if err := h.cartItemRepo.Update(existingItem); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update cart item quantity",
			})
		}
		
		// Get updated item with product details
		updatedItem, err := h.cartItemRepo.GetByID(existingItem.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve updated cart item",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Cart item quantity updated successfully",
			"data":    h.toCartItemResponse(updatedItem),
		})
	}

	// Create new cart item
	cartItem := &models.CartItem{
		CartID:    req.CartID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := h.cartItemRepo.Create(cartItem); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create cart item",
		})
	}

	// Get created item with product details
	createdItem, err := h.cartItemRepo.GetByID(cartItem.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve created cart item",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Cart item created successfully",
		"data":    h.toCartItemResponse(createdItem),
	})
}

// GetCartItem retrieves a cart item by ID
func (h *CartItemHandler) GetCartItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cart item ID",
		})
	}

	cartItem, err := h.cartItemRepo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cart item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart item",
		})
	}

	return c.JSON(fiber.Map{
		"data": h.toCartItemResponse(cartItem),
	})
}

// GetCartItems retrieves all items for a specific cart
func (h *CartItemHandler) GetCartItems(c *fiber.Ctx) error {
	cartID, err := strconv.ParseUint(c.Params("cartId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cart ID",
		})
	}

	cartItems, err := h.cartItemRepo.GetByCartID(uint(cartID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart items",
		})
	}

	var response []CartItemResponse
	for _, item := range cartItems {
		response = append(response, h.toCartItemResponse(&item))
	}

	return c.JSON(fiber.Map{
		"data": response,
	})
}

// UpdateCartItem updates a cart item
func (h *CartItemHandler) UpdateCartItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cart item ID",
		})
	}

	var req UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if cart item exists
	_, err = h.cartItemRepo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cart item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart item",
		})
	}

	// Update quantity
	if err := h.cartItemRepo.UpdateQuantity(uint(id), req.Quantity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart item",
		})
	}

	// Get updated item
	updatedItem, err := h.cartItemRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve updated cart item",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cart item updated successfully",
		"data":    h.toCartItemResponse(updatedItem),
	})
}

// DeleteCartItem removes a cart item
func (h *CartItemHandler) DeleteCartItem(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid cart item ID",
		})
	}

	// Check if cart item exists
	_, err = h.cartItemRepo.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cart item not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart item",
		})
	}

	if err := h.cartItemRepo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete cart item",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cart item deleted successfully",
	})
}

// Helper function to convert model to response
func (h *CartItemHandler) toCartItemResponse(cartItem *models.CartItem) CartItemResponse {
	return CartItemResponse{
		ID:        cartItem.ID,
		CartID:    cartItem.CartID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
		CreatedAt: cartItem.CreatedAt.Format(time.RFC3339),
		UpdatedAt: cartItem.UpdatedAt.Format(time.RFC3339),
		Product:   &cartItem.Product,
	}
}