package services

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/sgitwhyd/music-catalogue/internal/models"
	"github.com/sgitwhyd/music-catalogue/internal/repositorys"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:generate mockgen -source=user_service.go -destination=user_service_mock_test.go -package=services

type UserRepo interface {
	repositorys.UserRepository
} 

type UserService interface{
	Register(request models.SignUpRequest) error
}

type userService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService)  Register(request models.SignUpRequest) error {
	// check the user already registered
	_, err := s.userRepo.Find(request.Email, request.Username, 0)
	if err != gorm.ErrRecordNotFound {
		return errors.New("email or username already registered")
	}

	// bind with user model
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	body := models.User{
		Username: request.Username,
		Email: request.Email,
		Password: string(hashedPassword),
	}

	// create user
	err = s.userRepo.Upsert(body)
	if err != nil {
		log.Error().Err(err).Msgf("service create: error create with request email: %s. username: %s, id: %d", request.Email, request.Username, 0)
		return err
	}

	return nil
}