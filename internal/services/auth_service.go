package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

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

func (s *AuthService) RegisterUser(password, role string, outletID *uint, creatorID *uint, email *string, phoneNumber *string) (*models.User, error) {
	var generatedUsername string
	if email != nil {
		parts := strings.Split(*email, "@")
		if len(parts) > 0 {
			generatedUsername = parts[0]
		}
	}

	// Ensure username is unique
	finalUsername := generatedUsername
	for i := 0; ; i++ {
		var existingUser models.User
		err := s.DB.Where("username = ?", finalUsername).First(&existingUser).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break // Username is unique
		}
		if err != nil {
			log.Printf("Error checking username uniqueness: %v", err)
			return nil, errors.New("failed to check username uniqueness")
		}
		// Username exists, append random digits
		finalUsername = fmt.Sprintf("%s%05d", generatedUsername, rand.Intn(100000))
	}
	username := finalUsername

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, errors.New("failed to hash password")
	}

	user := models.User{
		Username:        username,
		Password:        hashedPassword,
		Role:            role,
		CreatorID:       creatorID,
		EmailVerifiedAt: nil, // User is not verified until OTP is confirmed
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

	// Generate and store OTP
	otpCode, err := GenerateOTP()
	if err != nil {
		log.Printf("Error generating OTP: %v", err)
		return nil, errors.New("failed to generate OTP")
	}

	hashedOTP, err := utils.HashOTP(otpCode)
	if err != nil {
		log.Printf("Error hashing OTP: %v", err)
		return nil, errors.New("failed to hash OTP")
	}

	otpRecord := models.OTP{
		UserID:    user.ID,
		OTP:       hashedOTP,
		Purpose:   "email_verification",
		Target:    user.Email,
		ExpiresAt: time.Now().Add(10 * time.Minute), // OTP valid for 10 minutes
	}

	if err := s.DB.Create(&otpRecord).Error; err != nil {
		log.Printf("Error saving OTP: %v", err)
		return nil, errors.New("failed to save OTP")
	}

	// Send verification email
	if err := SendVerificationEmail(user.Email, otpCode); err != nil {
		log.Printf("Error sending verification email: %v", err)
		// For now, we'll just log the error and not fail the registration
	}

	return &user, nil
}

func (s *AuthService) VerifyOTP(email, otp string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email")
	}

	var otpRecord models.OTP
	if err := s.DB.Where("user_id = ? AND purpose = ? AND target = ?", user.ID, "email_verification", email).First(&otpRecord).Error; err != nil {
		return nil, errors.New("OTP not found or already used")
	}

	// Check if OTP is expired
	if time.Now().After(otpRecord.ExpiresAt) {
		// Delete expired OTP
		s.DB.Delete(&otpRecord)
		return nil, errors.New("OTP expired")
	}

	// Check if OTP matches
	if !utils.CheckOTPHash(otp, otpRecord.OTP) {
		return nil, errors.New("invalid OTP")
	}

	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete OTP after successful verification
	if err := tx.Delete(&otpRecord).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting OTP record: %v", err)
		return nil, errors.New("failed to delete OTP record")
	}

	// If the user is an admin, create a default cash payment method
	if user.Role == "owner" {
		paymentMethod := models.PaymentMethod{
			Name:      "Cash",
			Type:      "cash",
			IsActive:  true,
			CreatorID: &user.ID,
		}
		if err := tx.Create(&paymentMethod).Error; err != nil {
			tx.Rollback()
			log.Printf("Error creating payment method: %v", err)
			return nil, errors.New("failed to create payment method")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
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

	if user.EmailVerifiedAt == nil {
		return "", nil, errors.New("user not verified")
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
