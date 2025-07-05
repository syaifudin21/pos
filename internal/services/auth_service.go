package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"github.com/msyaifudin/pos/pkg/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) RegisterUser(username, password, role string, outletID *uint, creatorID *uint, email *string, phoneNumber *string) (*models.User, error) {
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
		Username:  username,
		Password:  hashedPassword,
		Role:      role,
		CreatorID: creatorID,
	}

	// Assign email if not nil
	if email != nil {
		user.Email = *email
	}
	// Assign phoneNumber if not nil
	if phoneNumber != nil {
		user.PhoneNumber = *phoneNumber
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

	token, err := utils.GenerateToken(user.Username, user.Role, user.ID)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", nil, errors.New("failed to generate token")
	}

	return token, &user, nil
}

func (s *AuthService) BlockUser(userID uint) error {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
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

func (s *AuthService) UnblockUser(userID uint) error {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
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

func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) GetUserByuuid(userUuid uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("uuid = ?", userUuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) GetAllUsers(adminID uint) ([]models.User, error) {
	var users []models.User
	if err := s.DB.Where("creator_id = ?", adminID).Find(&users).Error; err != nil {
		log.Printf("Error getting all users: %v", err)
		return nil, errors.New("failed to retrieve users")
	}
	return users, nil
}

func (s *AuthService) UpdateUser(userID uint, updates *dtos.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		log.Printf("Error finding user for update: %v", err)
		return nil, errors.New("failed to retrieve user for update")
	}

	// Update fields if provided
	if updates.Username != nil {
		user.Username = *updates.Username
	}
	if updates.Password != nil {
		hashedPassword, err := utils.HashPassword(*updates.Password)
		if err != nil {
			log.Printf("Error hashing new password: %v", err)
			return nil, errors.New("failed to hash new password")
		}
		user.Password = hashedPassword
	}
	if updates.Role != nil {
		if !isValidRole(*updates.Role) {
			return nil, errors.New("invalid role specified")
		}
		user.Role = *updates.Role
	}

	if err := s.DB.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, errors.New("failed to update user")
	}

	return &user, nil
}

func (s *AuthService) DeleteUser(userID uint) error {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		log.Printf("Error finding user for deletion: %v", err)
		return errors.New("failed to retrieve user for deletion")
	}

	if err := s.DB.Delete(&user).Error; err != nil {
		log.Printf("Error deleting user: %v", err)
		return errors.New("failed to delete user")
	}

	return nil
}

func isValidRole(role string) bool {
	for _, r := range models.AllowedUserRoles {
		if r == role {
			return true
		}
	}
	return false
}
