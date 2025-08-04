package products

import (
	"strconv"

	"github.com/devops-learning/ecommerce/internals/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductHandler struct {
	repo ProductRepositoryImpl
}

func NewProductHandler(repo ProductRepositoryImpl) *ProductHandler {
	return &ProductHandler{repo: repo}
}


func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	if product.Name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Name is required")
	}
	if product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Price must be greater than 0")
	}

	if err := h.repo.Create(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create product")
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	product, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product")
	}

	return c.JSON(product)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	name := c.Query("name")
	var (
		products []models.Product
		err      error
	)

	if name != "" {
		products, err = h.repo.GetByName(name)
	} else {
		products, err = h.repo.GetAll()
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get products")
	}

	return c.JSON(products)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	existingProduct, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product")
	}

	var updatedProduct models.Product
	if err := c.BodyParser(&updatedProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	updatedProduct.ID = existingProduct.ID
	updatedProduct.CreatedAt = existingProduct.CreatedAt

	if updatedProduct.Name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Name is required")
	}
	if updatedProduct.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).SendString("Price must be greater than 0")
	}

	if err := h.repo.Update(&updatedProduct); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update product")
	}

	return c.JSON(updatedProduct)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	_, err = h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product")
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete product")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ProductHandler) UpdateProductStock(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	var stockUpdate struct {
		Stock uint `json:"stock"`
	}

	if err := c.BodyParser(&stockUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid JSON")
	}

	_, err = h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).SendString("Product not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product")
	}

	if err := h.repo.UpdateStock(uint(id), stockUpdate.Stock); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update stock")
	}

	updatedProduct, _ := h.repo.GetByID(uint(id))
	return c.JSON(updatedProduct)
}