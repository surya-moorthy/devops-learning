package auth

import (
	"github.com/devops-learning/ecommerce/internals/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
  DB *gorm.DB	
}

func(r *AuthRepository) CreateUser(user *models.UserModel) error {
     return nil
}

func(r *AuthRepository) GetUserByEmail(email *string,password *string) error {
return nil
}

func(r *AuthRepository) FindUserByEmail(email *string) error {
return nil
}


func(r *AuthRepository) UpdateUser(user *models.UserModel) error {
return nil
}