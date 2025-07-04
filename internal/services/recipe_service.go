package services

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/msyaifudin/pos/internal/models"
	"gorm.io/gorm"
)

type RecipeService struct {
	DB *gorm.DB
}

func NewRecipeService(db *gorm.DB) *RecipeService {
	return &RecipeService{DB: db}
}

// GetRecipeByUuid retrieves a recipe by its Uuid.
func (s *RecipeService) GetRecipeByUuid(uuid uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := s.DB.Preload("MainProduct").Preload("Component").Where("uuid = ?", uuid).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recipe not found")
		}
		log.Printf("Error getting recipe by uuid: %v", err)
		return nil, errors.New("failed to retrieve recipe")
	}
	return &recipe, nil
}

// GetRecipesByMainProduct retrieves all recipes for a given main product.
func (s *RecipeService) GetRecipesByMainProduct(mainProductUuid uuid.UUID) ([]models.Recipe, error) {
	var mainProduct models.Product
	if err := s.DB.Where("uuid = ?", mainProductUuid).First(&mainProduct).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var recipes []models.Recipe
	if err := s.DB.Preload("MainProduct").Preload("Component").Where("main_product_id = ?", mainProduct.ID).Find(&recipes).Error; err != nil {
		log.Printf("Error getting recipes by main product: %v", err)
		return nil, errors.New("failed to retrieve recipes")
	}
	return recipes, nil
}

// CreateRecipe creates a new recipe.
func (s *RecipeService) CreateRecipe(mainProductUuid, componentUuid uuid.UUID, quantity float64) (*models.Recipe, error) {
	var mainProduct models.Product
	if err := s.DB.Where("uuid = ?", mainProductUuid).First(&mainProduct).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var component models.Product

	if err := s.DB.Where("uuid = ?", componentUuid).First(&component).Error; err != nil {
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
	}

	if err := s.DB.Create(&recipe).Error; err != nil {
		log.Printf("Error creating recipe: %v", err)
		return nil, errors.New("failed to create recipe")
	}
	return &recipe, nil
}

// UpdateRecipe updates an existing recipe.
func (s *RecipeService) UpdateRecipe(uuid uuid.UUID, mainProductUuid, componentUuid uuid.UUID, quantity float64) (*models.Recipe, error) {
	var recipe models.Recipe
	if err := s.DB.Where("uuid = ?", uuid).First(&recipe).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recipe not found")
		}
		log.Printf("Error finding recipe for update: %v", err)
		return nil, errors.New("failed to retrieve recipe for update")
	}

	var mainProduct models.Product
	if err := s.DB.Where("uuid = ?", mainProductUuid).First(&mainProduct).Error; err != nil {
		return nil, errors.New("main product not found")
	}

	var component models.Product
	if err := s.DB.Where("uuid = ?", componentUuid).First(&component).Error; err != nil {
		return nil, errors.New("component product not found")
	}

	// Check if main product is of type fnb_main_product and component is fnb_component
	if mainProduct.Type != "fnb_main_product" || component.Type != "fnb_component" {
		return nil, errors.New("invalid product types for recipe: main product must be 'fnb_main_product' and component must be 'fnb_component'")
	}

	recipe.MainProductID = mainProduct.ID
	recipe.ComponentID = component.ID
	recipe.Quantity = quantity

	if err := s.DB.Save(&recipe).Error; err != nil {
		log.Printf("Error updating recipe: %v", err)
		return nil, errors.New("failed to update recipe")
	}
	return &recipe, nil
}

// DeleteRecipe deletes a recipe by its Uuid.
func (s *RecipeService) DeleteRecipe(uuid uuid.UUID) error {
	if err := s.DB.Where("uuid = ?", uuid).Delete(&models.Recipe{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("recipe not found")
		}
		log.Printf("Error deleting recipe: %v", err)
		return errors.New("failed to delete recipe")
	}
	return nil
}
