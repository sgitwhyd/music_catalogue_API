package repositorys

import (
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"gorm.io/gorm"
)


type UserRepository interface{
	Upsert(model models.User) error
	Find(email, username string, id uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Upsert(model models.User) error {
	err := r.db.Save(&model).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Find(email, username string, id uint) (*models.User, error) {
	user := models.User{}
	err := r.db.Where("email = ?", email).Or("username = ?", username).Or("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}