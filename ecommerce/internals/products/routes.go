	package products

	import (
		"github.com/devops-learning/ecommerce/internals/middleware"
		"github.com/gofiber/fiber/v2"
	)

	func (h *ProductHandler) ProductRoutes(app *fiber.App) {
	// All paths below start with /products
	user  := app.Group("/products", middleware.UserMiddleware)
	admin := app.Group("/products", middleware.AdminMiddleware)

	// -------- Public (authenticated-user) endpoints --------
	user.Get("/",      h.GetAllProducts) // GET /products
	user.Get("/:id",   h.GetProduct)     // GET /products/:id

	// -------- Admin-only endpoints --------
	admin.Post("/",            h.CreateProduct)       // POST   /products
	admin.Put("/:id",          h.UpdateProduct)       // PUT    /products/:id
	admin.Delete("/:id",       h.DeleteProduct)       // DELETE /products/:id
	admin.Patch("/:id/stock",  h.UpdateProductStock)  // PATCH  /products/:id/stock
}