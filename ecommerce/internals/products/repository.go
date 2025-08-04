package products

import (
	"github.com/devops-learning/ecommerce/internals/models"
	"gorm.io/gorm"
)


type ProductRepositoryImpl struct {
	db *gorm.DB
}

// Create creates a new product
func (r *ProductRepositoryImpl) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetByID retrieves a product by ID
func (r *ProductRepositoryImpl) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetAll retrieves all products
func (r *ProductRepositoryImpl) GetAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

// Update updates a product
func (r *ProductRepositoryImpl) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product by ID
func (r *ProductRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// GetByName searches products by name (case-insensitive)
func (r *ProductRepositoryImpl) GetByName(name string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%").Find(&products).Error
	return products, err
}

// UpdateStock updates only the stock of a product
func (r *ProductRepositoryImpl) UpdateStock(id uint, stock uint) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", stock).Error
}
