package services

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"math"
	"time"
)

type PostService struct {
	repo *repository.Manager
}

func NewPostService(repo *repository.Manager) *PostService {
	return &PostService{
		repo: repo,
	}
}

// GetSinglePost - keep this as it's not about pagination
func (s *PostService) GetSinglePost(postID string) (models.SinglePost, error) {
	var result models.SinglePost

	posts, err := s.repo.Get("post", "id", postID)
	if err != nil {
		return result, fmt.Errorf("failed to fetch post: %w", err)
	}
	if len(posts) == 0 {
		return result, fmt.Errorf("post not found")
	}

	post := posts[0]

	// get the user who created the post
	user, err := s.repo.Get("user", "id", post["user_id"].(string))
	if err != nil {
		return result, err
	}
	username := user[0]["username"].(string)

	// get the categories for the post
	categories, err := s.GetCategoriesForPost(post["id"].(string))
	if err != nil {
		return result, err
	}

	// get reactions
	likes, dislikes, err := s.GetReactionsForPost(post["id"].(string))
	if err != nil {
		return result, err
	}

	// get comments
	commentsmap, err := s.repo.Get("comment", "post_id", post["id"].(string))
	if err != nil {
		return result, err
	}

	comments, err := s.ConvertToCommentStruct(commentsmap)
	if err != nil {
		return result, err
	}

	result = models.SinglePost{
		ID:         post["id"].(string),
		Username:   username,
		Categories: categories,
		Title:      post["title"].(string),
		Content:    post["content"].(string),
		CreatedAt:  post["created_at"].(time.Time),
		Likes:      likes,
		Dislikes:   dislikes,
		Comments:   comments,
	}

	return result, nil
}

// PAGINATED FUNCTIONS
// =============================

// GetPaginatedAllPostsByEngagement - FIXED VERSION
func (s *PostService) GetPaginatedAllPostsByEngagement(page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count total posts
	countQuery := "SELECT COUNT(*) FROM post"
	var totalCount int
	err := s.repo.Db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// FIXED: Simplified engagement query using CASE statements instead of subqueries
	postsQuery := `
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.created_at,
			(
				SELECT COUNT(*) FROM reaction r 
				WHERE r.post_id = p.id AND r.type = 'like'
			) as likes,
			(
				SELECT COUNT(*) FROM reaction r 
				WHERE r.post_id = p.id AND r.type = 'dislike'
			) as dislikes,
			(
				SELECT COUNT(*) FROM comment c 
				WHERE c.post_id = p.id
			) as comments
		FROM post p
		ORDER BY (
			(SELECT COUNT(*) FROM reaction r WHERE r.post_id = p.id AND r.type = 'like') +
			(SELECT COUNT(*) FROM comment c WHERE c.post_id = p.id) -
			(SELECT COUNT(*) FROM reaction r WHERE r.post_id = p.id AND r.type = 'dislike')
		) DESC, p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.Db.Query(postsQuery, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts by engagement: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time
		var likes, dislikes, comments int

		err := rows.Scan(&id, &userID, &title, &content, &createdAt, &likes, &dislikes, &comments)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedAllNewestPosts - gets newest posts first, then paginates
func (s *PostService) GetPaginatedAllNewestPosts(page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count total posts
	countQuery := "SELECT COUNT(*) FROM post"
	var totalCount int
	err := s.repo.Db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	// Get posts ordered by creation date (newest first)
	postsQuery := `
		SELECT id, user_id, title, content, created_at
		FROM post
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.Db.Query(postsQuery, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get newest posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time

		err := rows.Scan(&id, &userID, &title, &content, &createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedUserPosts - for user's own posts
func (s *PostService) GetPaginatedUserPosts(userID string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count user's posts
	countQuery := "SELECT COUNT(*) FROM post WHERE user_id = ?"
	var totalCount int
	err := s.repo.Db.QueryRow(countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user posts: %w", err)
	}

	// Get user's posts
	postsQuery := `
		SELECT id, user_id, title, content, created_at
		FROM post
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.Db.Query(postsQuery, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, postUserID, title, content string
		var createdAt time.Time

		err := rows.Scan(&id, &postUserID, &title, &content, &createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    postUserID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedUserLikedPosts - for user's liked posts
func (s *PostService) GetPaginatedUserLikedPosts(userID string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count user's liked posts
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM post p
		JOIN reaction r ON p.id = r.post_id
		WHERE r.user_id = ? AND r.type = 'like' AND r.post_id IS NOT NULL
	`
	var totalCount int
	err := s.repo.Db.QueryRow(countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count liked posts: %w", err)
	}

	// Get user's liked posts
	postsQuery := `
		SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.created_at
		FROM post p
		JOIN reaction r ON p.id = r.post_id
		WHERE r.user_id = ? AND r.type = 'like' AND r.post_id IS NOT NULL
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.Db.Query(postsQuery, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get liked posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, postUserID, title, content string
		var createdAt time.Time

		err := rows.Scan(&id, &postUserID, &title, &content, &createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan liked post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    postUserID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedUserCommentedPosts - for user's commented posts
func (s *PostService) GetPaginatedUserCommentedPosts(userID string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count user's commented posts
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM post p
		JOIN comment c ON p.id = c.post_id
		WHERE c.user_id = ?
	`
	var totalCount int
	err := s.repo.Db.QueryRow(countQuery, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count commented posts: %w", err)
	}

	// Get user's commented posts
	postsQuery := `
		SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.created_at
		FROM post p
		JOIN comment c ON p.id = c.post_id
		WHERE c.user_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.Db.Query(postsQuery, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get commented posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, postUserID, title, content string
		var createdAt time.Time

		err := rows.Scan(&id, &postUserID, &title, &content, &createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan commented post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    postUserID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedPostsByCategoryByEngagement - FIXED VERSION
func (s *PostService) GetPaginatedPostsByCategoryByEngagement(categoryID string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count posts in category
	countQuery := `
		SELECT COUNT(*)
		FROM post p
		JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
	`
	var totalCount int
	err := s.repo.Db.QueryRow(countQuery, categoryID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count category posts: %w", err)
	}

	// Simplified engagement query for category posts
	postsQuery := `
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.created_at,
			(
				SELECT COUNT(*) FROM reaction r 
				WHERE r.post_id = p.id AND r.type = 'like'
			) as likes,
			(
				SELECT COUNT(*) FROM reaction r 
				WHERE r.post_id = p.id AND r.type = 'dislike'
			) as dislikes,
			(
				SELECT COUNT(*) FROM comment c 
				WHERE c.post_id = p.id
			) as comments
		FROM post p
		JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY (
			(SELECT COUNT(*) FROM reaction r WHERE r.post_id = p.id AND r.type = 'like') +
			(SELECT COUNT(*) FROM comment c WHERE c.post_id = p.id) -
			(SELECT COUNT(*) FROM reaction r WHERE r.post_id = p.id AND r.type = 'dislike')
		) DESC, p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.Db.Query(postsQuery, categoryID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get category posts by engagement: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time
		var likes, dislikes, comments int

		err := rows.Scan(&id, &userID, &title, &content, &createdAt, &likes, &dislikes, &comments)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedPostsByCategoryNewest - for category pages sorted by newest
func (s *PostService) GetPaginatedPostsByCategoryNewest(categoryID string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count posts in category
	countQuery := `
		SELECT COUNT(*)
		FROM post p
		JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
	`
	var totalCount int
	err := s.repo.Db.QueryRow(countQuery, categoryID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count category posts: %w", err)
	}

	// Get posts in category sorted by newest
	postsQuery := `
		SELECT p.id, p.user_id, p.title, p.content, p.created_at
		FROM post p
		JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.repo.Db.Query(postsQuery, categoryID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get newest category posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time

		err := rows.Scan(&id, &userID, &title, &content, &createdAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan newest category post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// HELPER FUNCTIONS
// ===============

// ConvertToPostStruct - keep as helper
func (s *PostService) ConvertToPostStruct(posts []map[string]interface{}) ([]models.Post, error) {
	var result []models.Post
	for _, post := range posts {
		// Get the user who created the post
		user, err := s.repo.Get("user", "id", post["user_id"].(string))
		if err != nil {
			return nil, err
		}

		// Get categories
		categories, err := s.GetCategoriesForPost(post["id"].(string))
		if err != nil {
			return nil, err
		}

		// Get reactions
		likes, dislikes, err := s.GetReactionsForPost(post["id"].(string))
		if err != nil {
			return nil, err
		}

		// Get comment count
		comments, err := s.repo.Get("comment", "post_id", post["id"].(string))
		if err != nil {
			return nil, err
		}

		username := user[0]["username"].(string)
		result = append(result, models.Post{
			ID:         post["id"].(string),
			Username:   username,
			Title:      post["title"].(string),
			Content:    post["content"].(string),
			CreatedAt:  post["created_at"].(time.Time),
			Categories: categories,
			Likes:      likes,
			Dislikes:   dislikes,
			Comments:   len(comments),
		})
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// ConvertToCommentStruct - keep as helper
func (s *PostService) ConvertToCommentStruct(comments []map[string]interface{}) ([]models.Comment, error) {
	var result []models.Comment
	for _, comment := range comments {
		// Get user
		user, err := s.repo.Get("user", "id", comment["user_id"].(string))
		if err != nil {
			return nil, err
		}
		username := user[0]["username"].(string)

		// Get reactions
		likes, dislikes, err := s.GetReactionsForComment(comment["id"].(string))
		if err != nil {
			return nil, err
		}

		result = append(result, models.Comment{
			ID:        comment["id"].(string),
			Username:  username,
			Content:   comment["content"].(string),
			CreatedAt: comment["created_at"].(time.Time),
			Likes:     likes,
			Dislikes:  dislikes,
		})
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// GetCategoriesForPost - keep as helper
func (s *PostService) GetCategoriesForPost(postID string) ([]models.Category, error) {
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

// GetReactionsForPost - keep as helper
func (s *PostService) GetReactionsForPost(postID string) (int, int, error) {
	reactions, err := s.repo.Get("reaction", "post_id", postID)
	if err != nil {
		return 0, 0, err
	}

	if reactions == nil {
		return 0, 0, nil
	}

	var likes, dislikes int
	for _, reaction := range reactions {
		if reactionType, ok := reaction["type"].(string); ok {
			if reactionType == "like" {
				likes++
			} else if reactionType == "dislike" {
				dislikes++
			}
		}
	}

	return likes, dislikes, nil
}

// GetReactionsForComment - keep as helper
func (s *PostService) GetReactionsForComment(commentID string) (int, int, error) {
	reactions, err := s.repo.Get("reaction", "comment_id", commentID)
	if err != nil {
		return 0, 0, err
	}

	if reactions == nil {
		return 0, 0, nil
	}

	var likes, dislikes int
	for _, reaction := range reactions {
		if reactionType, ok := reaction["type"].(string); ok {
			if reactionType == "like" {
				likes++
			} else if reactionType == "dislike" {
				dislikes++
			}
		}
	}

	return likes, dislikes, nil
}

// CreatePost creates a new post with associated categories
func (s *PostService) CreatePost(userID, title, content string, categoryIDs []string) (string, error) {
	// Insert post
	postData := map[string]interface{}{
		"user_id": userID,
		"title":   title,
		"content": content,
	}

	// Create the post record
	postID, err := s.repo.CreateRecord("post", postData)
	if err != nil {
		return "", err
	}

	// Associate categories with post
	for _, categoryID := range categoryIDs {
		junctionData := map[string]interface{}{
			"post_id":     postID,
			"category_id": categoryID,
		}

		err = s.repo.CreateJunction("post_category", junctionData)
		if err != nil {
			return "", err
		}
	}

	return postID, nil
}
