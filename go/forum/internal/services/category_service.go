package services

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type CategoryService struct {
	repo *repository.Manager
}

// NewCategoryService creates a new instance of CategoryService
func NewCategoryService(repo *repository.Manager) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// GetAllPostCategories retrieves all post categories
func (c *CategoryService) GetAllCategories() ([]models.Category, error) {
	categories, err := c.repo.GetAll("category")
	if err != nil {
		return nil, err
	}

	var result []models.Category
	for _, category := range categories {
		result = append(result, models.Category{
			ID:   category["id"].(string),
			Name: category["name"].(string),
		})
	}

	return result, nil
}
