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

// GetCategoriesForPost - keep as helper
func (s *CategoryService) GetCategoriesForPost(postID string) ([]models.Category, error) {
	postcategories, err := s.repo.Get("post_category", "post_id", postID)
	if err != nil {
		return nil, err
	}

	if postcategories == nil || len(postcategories) == 0 {
		return make([]models.Category, 0), nil
	}

	var result []models.Category
	for _, pc := range postcategories {
		if categoryID, ok := pc["category_id"].(string); ok {
			categoryData, err := s.repo.Get("category", "id", categoryID)
			if err != nil {
				continue
			}
			if categoryData != nil && len(categoryData) > 0 {
				category := categoryData[0]
				result = append(result, models.Category{
					ID:   category["id"].(string),
					Name: category["name"].(string),
				})
			}
		}
	}

	return result, nil
}
