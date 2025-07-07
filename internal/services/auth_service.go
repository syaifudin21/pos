package services

import (
	"errors"
	"log"
	"time"

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

func (s *AuthService) RegisterUser(name, email, password, role string, outletID *uint, creatorID *uint, phoneNumber *string, isGoogleAuth bool) (*models.User, error) {
	var hashedPassword string
	if password != "" {
		hashed, err := utils.HashPassword(password)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return nil, errors.New("failed to hash password")
		}
		hashedPassword = hashed
	}

	user := models.User{
		Name:            name,
		Email:           email,
		Password:        hashedPassword,
		Role:            role,
		CreatorID:       creatorID,
		EmailVerifiedAt: nil, // User is not verified until OTP is confirmed
	}

	// Assign phoneNumber if not nil
	if phoneNumber != nil {
		user.PhoneNumber = *phoneNumber
	}

	if err := s.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, errors.New("failed to create user")
	}

	if !isGoogleAuth {
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

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) LoginUser(email, password string) (string, *models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
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

	if user.Password != "" && !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.Email, user.Role, user.ID)
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

func (s *AuthService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	log.Printf("AuthService.UpdatePassword: UserID: %d, OldPassword: %s, NewPassword: %s", userID, oldPassword, newPassword)

	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("AuthService.UpdatePassword: User not found: %v", err)
			return errors.New("user not found")
		}
		log.Printf("AuthService.UpdatePassword: Error finding user for password update: %v", err)
		return errors.New("failed to retrieve user for password update")
	}

	// If user has a password set, validate old password
	if user.Password != "" {
		log.Printf("AuthService.UpdatePassword: User has existing password. Checking old password.")
		if !utils.CheckPasswordHash(oldPassword, user.Password) {
			log.Printf("AuthService.UpdatePassword: Invalid old password for user %d", userID)
			return errors.New("invalid old password")
		}
	} else {
		log.Printf("AuthService.UpdatePassword: User has no existing password. Skipping old password check.")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		log.Printf("AuthService.UpdatePassword: Error hashing new password: %v", err)
		return errors.New("failed to hash new password")
	}

	user.Password = hashedPassword
	if err := s.DB.Save(&user).Error; err != nil {
		log.Printf("AuthService.UpdatePassword: Error updating password in DB: %v", err)
		return errors.New("failed to update password")
	}

	log.Printf("AuthService.UpdatePassword: Password updated successfully for user %d", userID)
	return nil
}

func (s *AuthService) SendOTPForEmailUpdate(userID uint, newEmail string) error {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		log.Printf("Error finding user for email update OTP: %v", err)
		return errors.New("failed to retrieve user")
	}

	// Check if new email is already in use by another user
	var existingUser models.User
	if err := s.DB.Where("email = ? AND id != ?", newEmail, userID).First(&existingUser).Error; err == nil {
		return errors.New("email already in use by another account")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error checking new email existence: %v", err)
		return errors.New("database error checking email")
	}

	// Generate and store OTP
	otpCode, err := GenerateOTP()
	if err != nil {
		log.Printf("Error generating OTP: %v", err)
		return errors.New("failed to generate OTP")
	}

	hashedOTP, err := utils.HashOTP(otpCode)
	if err != nil {
		log.Printf("Error hashing OTP: %v", err)
		return errors.New("failed to hash OTP")
	}

	// Delete any existing OTPs for email update for this user
	s.DB.Where("user_id = ? AND purpose = ?", userID, "email_update").Delete(&models.OTP{})

	otpRecord := models.OTP{
		UserID:    userID,
		OTP:       hashedOTP,
		Purpose:   "email_update",
		Target:    newEmail,
		ExpiresAt: time.Now().Add(10 * time.Minute), // OTP valid for 10 minutes
	}

	if err := s.DB.Create(&otpRecord).Error; err != nil {
		log.Printf("Error saving OTP for email update: %v", err)
		return errors.New("failed to save OTP")
	}

	// Send verification email
	if err := SendVerificationEmail(newEmail, otpCode); err != nil {
		log.Printf("Error sending verification email for email update: %v", err)
		return errors.New("failed to send verification email")
	}

	return nil
}

func (s *AuthService) UpdateEmail(userID uint, newEmail, otp string) error {
	var user models.User
	if err := s.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		log.Printf("Error finding user for email update: %v", err)
		return errors.New("failed to retrieve user")
	}

	var otpRecord models.OTP
	if err := s.DB.Where("user_id = ? AND purpose = ? AND target = ?", userID, "email_update", newEmail).First(&otpRecord).Error; err != nil {
		return errors.New("OTP not found or invalid for this email")
	}

	// Check if OTP is expired
	if time.Now().After(otpRecord.ExpiresAt) {
		s.DB.Delete(&otpRecord) // Delete expired OTP
		return errors.New("OTP expired")
	}

	// Check if OTP matches
	if !utils.CheckOTPHash(otp, otpRecord.OTP) {
		return errors.New("invalid OTP")
	}

	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	user.Email = newEmail
	user.EmailVerifiedAt = func() *time.Time { t := time.Now(); return &t }() // Mark new email as verified
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating user email: %v", err)
		return errors.New("failed to update email")
	}

	// Delete OTP after successful verification
	if err := tx.Delete(&otpRecord).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting OTP record after email update: %v", err)
		return errors.New("failed to delete OTP record")
	}

	if err := tx.Commit().Error; err != nil {
		return tx.Error
	}

	return nil
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
