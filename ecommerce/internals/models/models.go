package models

import (
	"time"
	"gorm.io/gorm"
)

	// User represents a user in the system
	type User struct {
		ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
		Email     string         `gorm:"not null;unique;size:255" json:"email"`
		Username  string         `gorm:"size:50" json:"username,omitempty"`
		Password  string         `gorm:"not null;size:255" json:"password"` // Never expose password in JSON
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
		
		// Relationships
		Cart   *Cart   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cart,omitempty"`
		Orders []Order `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"orders,omitempty"`
	}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// Product represents a product in the catalog
type Product struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `gorm:"not null" json:"price"`
	Stock       uint    `gorm:"not null" json:"stock"`
	ImageURL    string  `json:"image_url,omitempty"`
}

// Cart represents a user's shopping cart
type Cart struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint           `gorm:"not null;unique;index" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	User      User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CartItems []CartItem `gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cart_items,omitempty"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CartID    uint           `gorm:"not null;index" json:"cart_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	Quantity  uint           `gorm:"not null;check:quantity > 0" json:"quantity"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	Cart    Cart    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Product Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product"`
}

// Ensure unique cart item per product per cart
func (CartItem) TableName() string {
	return "cart_items"
}

// OrderStatus represents the possible order statuses
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
)

// Order represents a customer order
type Order struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	OrderNumber string         `gorm:"not null;unique;size:50" json:"order_number"`
	Status      OrderStatus    `gorm:"not null;default:'pending';size:20;index" json:"status"`
	TotalAmount float64        `gorm:"not null;type:decimal(10,2);check:total_amount >= 0" json:"total_amount"`
	ShippingAddress *string    `gorm:"type:text" json:"shipping_address,omitempty"`
	Notes       *string        `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	User       User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order_items,omitempty"`
}

// OrderItem represents an item within an order
type OrderItem struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	Quantity  uint           `gorm:"not null;check:quantity > 0" json:"quantity"`
	Price     float64        `gorm:"not null;type:decimal(10,2);check:price >= 0" json:"price"` // Price at time of order
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	Order   Order   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Product Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password,omitempty" validate:"min=8"`
}

// MigrateAll runs AutoMigrate on all models
func MigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Product{},
		&Cart{},
		&CartItem{},
		&Order{},
		&OrderItem{},
	)
}
