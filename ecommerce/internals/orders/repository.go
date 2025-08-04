package orders

import (
	"time"

	"github.com/devops-learning/ecommerce/internals/models"
	"gorm.io/gorm"
)

// ============== REPOSITORY ==============

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder creates a new order with items
func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

// GetOrderByID retrieves an order by ID with related data
func (r *OrderRepository) GetOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems.Product").Preload("User").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderByOrderNumber retrieves an order by order number
func (r *OrderRepository) GetOrderByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems.Product").Preload("User").Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (r *OrderRepository) GetOrdersByUserID(userID uint, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// Count total orders
	if err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get orders with pagination
	err := r.db.Preload("OrderItems.Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

// GetAllOrders retrieves all orders with pagination
func (r *OrderRepository) GetAllOrders(limit, offset int, status models.OrderStatus) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total orders
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get orders with pagination
	query = r.db.Preload("OrderItems.Product").Preload("User")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

// UpdateOrderStatus updates only the order status
func (r *OrderRepository) UpdateOrderStatus(id uint, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateOrder updates order information (excluding items)
func (r *OrderRepository) UpdateOrder(order *models.Order) error {
	return r.db.Model(order).Updates(map[string]interface{}{
		"shipping_address": order.ShippingAddress,
		"notes":           order.Notes,
		"total_amount":    order.TotalAmount,
		"updated_at":      time.Now(),
	}).Error
}

// DeleteOrder soft deletes an order
func (r *OrderRepository) DeleteOrder(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

// AddOrderItem adds an item to an existing order
func (r *OrderRepository) AddOrderItem(orderItem *models.OrderItem) error {
	return r.db.Create(orderItem).Error
}

// UpdateOrderItem updates an order item
func (r *OrderRepository) UpdateOrderItem(orderItem *models.OrderItem) error {
	return r.db.Save(orderItem).Error
}

// RemoveOrderItem removes an item from an order
func (r *OrderRepository) RemoveOrderItem(orderID, productID uint) error {
	return r.db.Where("order_id = ? AND product_id = ?", orderID, productID).Delete(&models.OrderItem{}).Error
}

// GetOrderItemsByOrderID retrieves all items for a specific order
func (r *OrderRepository) GetOrderItemsByOrderID(orderID uint) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.Preload("Product").Where("order_id = ?", orderID).Find(&items).Error
	return items, err
}
