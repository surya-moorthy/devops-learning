package auth

import (
	"github.com/devops-learning/ecommerce/internals/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
  DB *gorm.DB	
}

func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *AuthRepository) GetUserByEmail(email *string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", *email).First(&user).Error
	if err != nil {
		return nil, err
	}
	// Optionally verify password hash here
	return &user, nil
}

func (r *AuthRepository) FindUserByEmail(email *string) (*models.User, error) {
	var user models.User
	result := r.DB.Where("email = ?", *email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}
