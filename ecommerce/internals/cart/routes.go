package cart

import (
	"github.com/devops-learning/ecommerce/internals/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupCartItemRoutes(app *fiber.App, cartItemHandler *CartItemHandler) {
	
	// Cart item routes
	cartItems := app.Group("/cart-items",middleware.UserMiddleware)
	cartItems.Post("/", cartItemHandler.CreateCartItem)       // POST /api/v1/cart-items
	cartItems.Get("/:id", cartItemHandler.GetCartItem)       // GET /api/v1/cart-items/:id
	cartItems.Put("/:id", cartItemHandler.UpdateCartItem)    // PUT /api/v1/cart-items/:id
	cartItems.Delete("/:id", cartItemHandler.DeleteCartItem) // DELETE /api/v1/cart-items/:id
	
	// Cart-specific routes
	carts := app.Group("/carts",middleware.UserMiddleware)
	carts.Get("/:cartId/items", cartItemHandler.GetCartItems) // GET /api/v1/carts/:cartId/items
}