package services

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/pkg/utils"
	"gorm.io/gorm"
)

type UserContextService struct {
	DB *gorm.DB
}

func NewUserContextService(db *gorm.DB) *UserContextService {
	return &UserContextService{DB: db}
}

// GetUserIDFromEchoContext extracts the user ID (uint) from the Echo context.
// This assumes that the authentication middleware has set the user information in the context
// as a *jwt.Token with Claims of type *utils.Claims.
func (s *UserContextService) GetUserIDFromEchoContext(c echo.Context) (uint, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok || user == nil {
		return 0, fmt.Errorf("user token not found in context")
	}

	claims, ok := user.Claims.(*utils.Claims)
	if !ok || claims == nil {
		return 0, fmt.Errorf("user claims not found in token")
	}

	return claims.ID, nil
}

// GetOwnerID retrieves the owner's ID for a given user.
// If the user is a manager or cashier, it returns their creator's ID.
// Otherwise, it returns the user's own ID.
func (s *UserContextService) GetOwnerID(userID uint) (uint, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		log.Printf("Error finding user: %v", err)
		return 0, errors.New("user not found")
	}

	if (user.Role == "manager" || user.Role == "cashier") && user.CreatorID != nil {
		return *user.CreatorID, nil
	}

	return userID, nil
}
