package services

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/database"
	"github.com/msyaifudin/pos/internal/models"
	"github.com/msyaifudin/pos/internal/models/dtos"
	"gorm.io/gorm"
)

type RecipeService struct {
	DB *gorm.DB
}

func NewRecipeService(db *gorm.DB) *RecipeService {
	return &RecipeService{DB: db}
}

// GetOwnerID retrieves the owner's ID for a given user.
// If the user is a manager or cashier, it returns their creator's ID.
// Otherwise, it returns the user's own ID.
func (s *RecipeService) GetOwnerID(userID uint) (uint, error) {
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

// GetRecipeByUuid retrieves a recipe by its Uuid.
func (s *RecipeService) GetRecipeByUuid(uuid uuid.UUID, userID uint) (*dtos.RecipeResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var recipe models.Recipe
	if err := s.DB.Preload("MainProduct").Preload("Component").Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recipe not found")
		}
		log.Printf("Error getting recipe by uuid: %v", err)
		return nil, errors.New("failed to retrieve recipe")
	}
	return &dtos.RecipeResponse{
		ID:              recipe.ID,
		Uuid:            recipe.Uuid,
		MainProductID:   recipe.MainProductID,
		MainProductUuid: recipe.MainProduct.Uuid,
		ComponentID:     recipe.ComponentID,
		ComponentUuid:   recipe.Component.Uuid,
		Quantity:        recipe.Quantity,
	}, nil
}

// GetRecipesByMainProduct retrieves all recipes for a given main product.
func (s *RecipeService) GetRecipesByMainProduct(mainProductUuid uuid.UUID, userID uint) ([]dtos.RecipeResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var mainProduct models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", mainProductUuid, ownerID).First(&mainProduct).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var recipes []models.Recipe
	if err := s.DB.Preload("MainProduct").Preload("Component").Where("main_product_id = ? AND user_id = ?", mainProduct.ID, ownerID).Find(&recipes).Error; err != nil {
		log.Printf("Error getting recipes by main product: %v", err)
		return nil, errors.New("failed to retrieve recipes")
	}
	var recipeResponses []dtos.RecipeResponse
	for _, recipe := range recipes {
		recipeResponses = append(recipeResponses, dtos.RecipeResponse{
			ID:              recipe.ID,
			Uuid:            recipe.Uuid,
			MainProductID:   recipe.MainProductID,
			MainProductUuid: recipe.MainProduct.Uuid,
			ComponentID:     recipe.ComponentID,
			ComponentUuid:   recipe.Component.Uuid,
			Quantity:        recipe.Quantity,
		})
	}
	return recipeResponses, nil
}

// CreateRecipe creates a new recipe.
func (s *RecipeService) CreateRecipe(mainProductUuid, componentUuid uuid.UUID, quantity float64, userID uint) (*dtos.RecipeResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var mainProduct models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", mainProductUuid, ownerID).First(&mainProduct).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var component models.Product

	if err := s.DB.Where("uuid = ? AND user_id = ?", componentUuid, ownerID).First(&component).Error; err != nil {
		return nil, errors.New("component product not found")
	}

	// Check if main product is of type fnb_main_product and component is fnb_component
	if mainProduct.Type != "fnb_main_product" || component.Type != "fnb_component" {
		return nil, errors.New("invalid product types for recipe: main product must be 'fnb_main_product' and component must be 'fnb_component'")
	}

	recipe := models.Recipe{
		MainProductID: mainProduct.ID,
		ComponentID:   component.ID,
		Quantity:      quantity,
		UserID:        ownerID,
	}

	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Create(&recipe).Error; err != nil {
		log.Printf("Error creating recipe: %v", err)
		return nil, errors.New("failed to create recipe")
	}
	return &dtos.RecipeResponse{
		ID:              recipe.ID,
		Uuid:            recipe.Uuid,
		MainProductID:   recipe.MainProductID,
		MainProductUuid: mainProduct.Uuid,
		ComponentID:     recipe.ComponentID,
		ComponentUuid:   component.Uuid,
		Quantity:        recipe.Quantity,
	}, nil
}

// UpdateRecipe updates an existing recipe.
func (s *RecipeService) UpdateRecipe(uuid uuid.UUID, mainProductUuid, componentUuid uuid.UUID, quantity float64, userID uint) (*dtos.RecipeResponse, error) {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return nil, err
	}
	var recipe models.Recipe
	if err := s.DB.Where("uuid = ? AND user_id = ?", uuid, ownerID).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recipe not found")
		}
		log.Printf("Error finding recipe for update: %v", err)
		return nil, errors.New("failed to retrieve recipe for update")
	}

	var mainProduct models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", mainProductUuid, ownerID).First(&mainProduct).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var component models.Product
	if err := s.DB.Where("uuid = ? AND user_id = ?", componentUuid, ownerID).First(&component).Error; err != nil {
		return nil, errors.New("component product not found")
	}

	// Check if main product is of type fnb_main_product and component is fnb_component
	if mainProduct.Type != "fnb_main_product" || component.Type != "fnb_component" {
		return nil, errors.New("invalid product types for recipe: main product must be 'fnb_main_product' and component must be 'fnb_component'")
	}

	recipe.MainProductID = mainProduct.ID
	recipe.ComponentID = component.ID
	recipe.Quantity = quantity

	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Save(&recipe).Error; err != nil {
		log.Printf("Error updating recipe: %v", err)
		return nil, errors.New("failed to update recipe")
	}
	return &dtos.RecipeResponse{
		ID:              recipe.ID,
		Uuid:            recipe.Uuid,
		MainProductID:   recipe.MainProductID,
		MainProductUuid: mainProduct.Uuid,
		ComponentID:     recipe.ComponentID,
		ComponentUuid:   component.Uuid,
		Quantity:        recipe.Quantity,
	}, nil
}

// DeleteRecipe deletes a recipe by its Uuid.
func (s *RecipeService) DeleteRecipe(uuid uuid.UUID, userID uint) error {
	ownerID, err := s.GetOwnerID(userID)
	if err != nil {
		return err
	}
	if err := s.DB.WithContext(context.WithValue(context.Background(), database.UserIDContextKey, userID)).Where("uuid = ? AND user_id = ?", uuid, ownerID).Delete(&models.Recipe{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("recipe not found")
		}
		log.Printf("Error deleting recipe: %v", err)
		return errors.New("failed to delete recipe")
	}
	return nil
}
