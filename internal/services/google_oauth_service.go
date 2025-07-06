package services

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/pkg/utils"
)

type GoogleOAuthService struct {
	DB          *gorm.DB
	AuthService *AuthService
	GoogleOauthConfig *oauth2.Config
}

func NewGoogleOAuthService(db *gorm.DB, authService *AuthService) *GoogleOAuthService {
	return &GoogleOAuthService{
		DB:          db,
		AuthService: authService,
		GoogleOauthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (s *GoogleOAuthService) GetGoogleLoginURL() string {
	return s.GoogleOauthConfig.AuthCodeURL("randomstate") // "randomstate" is a CSRF token
}

func (s *GoogleOAuthService) HandleGoogleCallback(code string) (string, *models.User, error) {
	token, err := s.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange failed: %v", err)
		return "", nil, errors.New("failed to exchange code for token")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		return "", nil, errors.New("failed to get user info from Google")
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return "", nil, errors.New("failed to read Google response")
	}

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		VerifiedEmail bool `json:"verified_email"`
	}
	if err := json.Unmarshal(contents, &userInfo); err != nil {
		log.Printf("Failed to unmarshal user info: %v", err)
		return "", nil, errors.New("failed to parse Google user info")
	}

	// Check if user exists in your database
	var user models.User
	err = s.DB.Where("email = ?", userInfo.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User does not exist, auto-register
			log.Printf("User not found, auto-registering: %s", userInfo.Email)
			// Generate a random password for the new user
			randomPassword := utils.GenerateRandomString(16) // You need to implement this utility
			// Use the AuthService to register the user
			registeredUser, regErr := s.AuthService.RegisterUser(randomPassword, "owner", nil, nil, &userInfo.Email, nil)
			if regErr != nil {
				log.Printf("Failed to auto-register user: %v", regErr)
				return "", nil, errors.New("failed to auto-register user")
			}
			user = *registeredUser
			// If Google verified the email, set email_verified_at
			if userInfo.VerifiedEmail {
				now := time.Now()
				user.EmailVerifiedAt = &now
				if err := s.DB.Save(&user).Error; err != nil {
					log.Printf("Error updating email_verified_at for new user: %v", err)
				}
			}
		} else {
			log.Printf("Database error finding user: %v", err)
			return "", nil, errors.New("database error")
		}
	} else {
		// User exists, check if email needs to be marked as verified
		if userInfo.VerifiedEmail && user.EmailVerifiedAt == nil {
			now := time.Now()
			user.EmailVerifiedAt = &now
			if err := s.DB.Save(&user).Error; err != nil {
				log.Printf("Error updating email_verified_at for existing user: %v", err)
			}
		}
	}

	// If user exists or was just registered, generate JWT token
	tokenString, err := utils.GenerateToken(user.Username, user.Role, user.ID)
	if err != nil {
		log.Printf("Failed to generate token for user %s: %v", user.Username, err)
		return "", nil, errors.New("failed to generate authentication token")
	}

	return tokenString, &user, nil
}
