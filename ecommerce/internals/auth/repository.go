package auth

import (
	"github.com/devops-learning/ecommerce/internals/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
  DB *gorm.DB	
}

func (r *AuthRepository) CreateUser(user *models.UserModel) error {
	return r.DB.Create(user).Error
}

func (r *AuthRepository) GetUserByEmail(email *string, password *string) (*models.UserModel, error) {
	var user models.UserModel
	err := r.DB.Where("email = ?", *email).First(&user).Error
	if err != nil {
		return nil, err
	}
	// Optionally verify password hash here
	return &user, nil
}


func (r *AuthRepository) FindUserByEmail(email *string) error {
	var user models.UserModel
	return r.DB.Where("email = ?", *email).First(&user).Error
}

func (r *AuthRepository) UpdateUser(user *models.UserModel) error {
	return r.DB.Save(user).Error
}
