package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/pkg/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) RegisterUser(username, password, role string, outletID *uint) (*models.User, error) {
	// Check if username already exists
	var existingUser models.User
	if s.DB.Where("username = ?", username).First(&existingUser).Error == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	user := models.User{
		Username: username,
		Password: hashedPassword,
		Role:     role,
		OutletID: outletID,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}

func (s *AuthService) LoginUser(username, password string) (string, *models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		log.Printf("Error finding user: %v", err)
		return "", nil, errors.New("database error")
	}

	if user.IsBlocked {
		return "", nil, errors.New("user is blocked")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.Username, user.Role, user.OutletID, user.Uuid)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", nil, errors.New("failed to generate token")
	}

	return token, &user, nil
}

func (s *AuthService) BlockUser(userUuid uuid.UUID) error {
	var user models.User
	if err := s.DB.Where("uuid = ?", userUuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		log.Printf("Error finding user to block: %v", err)
		return errors.New("failed to retrieve user")
	}

	user.IsBlocked = true
	if err := s.DB.Save(&user).Error; err != nil {
		log.Printf("Error blocking user: %v", err)
		return errors.New("failed to block user")
	}
	return nil
}

func (s *AuthService) UnblockUser(useruuid uuid.UUID) error {
	var user models.User
	if err := s.DB.Where("uuid = ?", useruuid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		log.Printf("Error finding user to unblock: %v", err)
		return errors.New("failed to retrieve user")
	}

	user.IsBlocked = false
	if err := s.DB.Save(&user).Error; err != nil {
		log.Printf("Error unblocking user: %v", err)
		return errors.New("failed to unblock user")
	}
	return nil
}

func (s *AuthService) GetUserByuuid(uuid uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
