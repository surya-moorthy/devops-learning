package cart

import (
	"github.com/devops-learning/ecommerce/internals/models"
	"gorm.io/gorm"
)

type cartItemRepository struct {
	db *gorm.DB
}

func (r *cartItemRepository) Create(cartItem *models.CartItem) error {
	return r.db.Create(cartItem).Error
}

func (r *cartItemRepository) GetByID(id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Preload("Product").First(&cartItem, id).Error
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *cartItemRepository) GetByCartID(cartID uint) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	err := r.db.Preload("Product").Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}

func (r *cartItemRepository) GetByCartAndProduct(cartID, productID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *cartItemRepository) Update(cartItem *models.CartItem) error {
	return r.db.Save(cartItem).Error
}

func (r *cartItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}

func (r *cartItemRepository) UpdateQuantity(id uint, quantity uint) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", id).Update("quantity", quantity).Error
}

