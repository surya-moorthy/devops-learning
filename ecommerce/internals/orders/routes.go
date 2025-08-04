package orders

import "github.com/gofiber/fiber/v2"

func SetupOrderRoutes(app *fiber.App, handler *OrderHandler) {
	api := app.Group("/api/v1")
	
	// Order routes
	orders := api.Group("/orders")
	
	// Create a new order
	orders.Post("/", handler.CreateOrder)
	
	// Get all orders (with pagination and filtering)
	orders.Get("/", handler.GetAllOrders)
	
	// Get order by ID
	orders.Get("/:id", handler.GetOrder)
	
	// Get order by order number
	orders.Get("/number/:orderNumber", handler.GetOrderByNumber)
	
	// Update order
	orders.Put("/:id", handler.UpdateOrder)
	
	// Update order status
	orders.Patch("/:id/status", handler.UpdateOrderStatus)
	
	// Delete order
	orders.Delete("/:id", handler.DeleteOrder)
	
	// User-specific order routes
	users := api.Group("/users")
	
	// Get all orders for a specific user
	users.Get("/:userID/orders", handler.GetUserOrders)
}
