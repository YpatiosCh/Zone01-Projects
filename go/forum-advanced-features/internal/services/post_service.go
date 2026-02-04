package services

import (
	"errors"
	"fmt"
	"forum/internal/config"
	"forum/internal/models"
	"forum/internal/repository"
	"forum/internal/utils/image"
	"math"
	"net/http"
	"time"
)

type PostService struct {
	repo     *repository.Manager
	reaction *ReactionService
	category *CategoryService
	user     *UserService
	comment  *CommentService
	config   *config.AppConfig
}

func NewPostService(repo *repository.Manager, reaction *ReactionService, category *CategoryService, user *UserService, comment *CommentService, config *config.AppConfig) *PostService {

	return &PostService{
		repo:     repo,
		reaction: reaction,
		category: category,
		user:     user,
		comment:  comment,
		config:   config,
	}
}

// GetSinglePost gets a single post by ID, including user, categories, reactions, and comments
// This function is used for displaying a single post in detail
func (s *PostService) GetSinglePost(postID string) (models.Post, error) {
	var result models.Post

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
	categories, err := s.category.GetCategoriesForPost(post["id"].(string))
	if err != nil {
		return result, err
	}

	// get reactions
	likes, dislikes, err := s.reaction.GetReactionsForPost(post["id"].(string))
	if err != nil {
		return result, err
	}

	// get comments
	commentsmap, err := s.repo.Get("comment", "post_id", post["id"].(string))
	if err != nil {
		return result, err
	}

	comments, err := s.comment.ConvertToCommentStruct(commentsmap)
	if err != nil {
		return result, err
	}

	// Check if post has an image
	hasImage := false
	var filepath string

	if hasImageField, ok := post["has_image"].(bool); ok && hasImageField {
		hasImage = true
		// Get image filename from database
		if post["image"] != nil {
			filepath = post["image"].(string)
		}
	} else if hasImageField, ok := post["has_image"].(int64); ok && hasImageField == 1 {
		hasImage = true
		// Get image filename from database
		if post["image"] != nil {
			filepath = post["image"].(string)
		}
	}

	result = models.Post{
		ID:         post["id"].(string),
		Username:   username,
		Categories: categories,
		Title:      post["title"].(string),
		Content:    post["content"].(string),
		CreatedAt:  post["created_at"].(time.Time),
		Likes:      likes,
		Dislikes:   dislikes,
		Comments:   comments,
		HasImage:   hasImage,
		Image:      filepath,
	}

	return result, nil
}

// PAGINATED FUNCTIONS
// =============================

// GetPaginatedAllPostsByEngagement gets all posts ordered by engagement (likes, comments, dislikes)
func (s *PostService) GetPaginatedAllPostsByEngagement(page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count total posts
	var totalCount int
	err := s.repo.Db.QueryRow(repository.CountAllPosts).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.AllByEngagement, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts by engagement: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image *string // Use pointer to handle NULL values
		var likes, dislikes, comments int

		// Fixed scan order to match SELECT order
		err := rows.Scan(&id, &userID, &title, &content, &createdAt, &hasImage, &image, &likes, &dislikes, &comments)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}

		// Handle NULL image field
		var imageFilename string
		if image != nil {
			imageFilename = *image
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      imageFilename,
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
	var totalCount int
	err := s.repo.Db.QueryRow(repository.CountAllPosts).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.AllNewest, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get newest posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image string

		err := rows.Scan(&id, &userID, &title, &content, &createdAt, &hasImage, &image)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      image,
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
func (s *PostService) GetPaginatedUserPosts(username string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	userID, err := s.repo.GetUserIDbyUsername(username)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user ID by username: %w", err)
	}
	if userID == "" {
		return nil, 0, fmt.Errorf("user not found")
	}

	// Count user's posts
	var totalCount int
	err = s.repo.Db.QueryRow(repository.CountUserPosts, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.UserPosts, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, postUserID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image string

		err := rows.Scan(&id, &postUserID, &title, &content, &createdAt, &hasImage, &image)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user post: %w", err)
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    postUserID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      image,
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
func (s *PostService) GetPaginatedUserLikedPosts(username string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	userID, err := s.repo.GetUserIDbyUsername(username)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user ID by username: %w", err)
	}
	if userID == "" {
		return nil, 0, fmt.Errorf("user not found")
	}

	// Count user's liked posts
	var totalCount int
	err = s.repo.Db.QueryRow(repository.CountUserLikedPosts, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count liked posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.UserLikedPosts, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get liked posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, postUserID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image *string

		err := rows.Scan(&id, &postUserID, &title, &content, &createdAt, &hasImage, &image)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan liked post: %w", err)
		}

		// Handle NULL image field
		var imageFilename string
		if image != nil {
			imageFilename = *image
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    postUserID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      imageFilename,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedUserDisikedPosts - for user's disliked posts
func (s *PostService) GetPaginatedUserDislikedPosts(username string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	userID, err := s.repo.GetUserIDbyUsername(username)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user ID by username: %w", err)
	}
	if userID == "" {
		return nil, 0, fmt.Errorf("user not found")
	}

	// Count user's liked posts
	var totalCount int
	err = s.repo.Db.QueryRow(repository.CountUserDislikedPosts, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count liked posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.UserDislikedPosts, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get liked posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, postUserID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image *string

		err := rows.Scan(&id, &postUserID, &title, &content, &createdAt, &hasImage, &image)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan liked post: %w", err)
		}

		// Handle NULL image field
		var imageFilename string
		if image != nil {
			imageFilename = *image
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    postUserID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      imageFilename,
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
func (s *PostService) GetPaginatedUserCommentedPosts(username string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	userID, err := s.repo.GetUserIDbyUsername(username)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user ID by username: %w", err)
	}
	if userID == "" {
		return nil, 0, fmt.Errorf("user not found")
	}

	// Count user's commented posts
	var totalCount int
	err = s.repo.Db.QueryRow(repository.CountUserCommentedPosts, userID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count commented posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.UserCommentedPosts, userID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get commented posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, postUserID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image *string

		err := rows.Scan(&id, &postUserID, &title, &content, &createdAt, &hasImage, &image)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan commented post: %w", err)
		}

		// Handle NULL image field
		var imageFilename string
		if image != nil {
			imageFilename = *image
		}
		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    postUserID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      imageFilename,
		})
	}

	result, err := s.ConvertToPostStruct(posts)
	if err != nil {
		return nil, 0, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))
	return result, totalPages, nil
}

// GetPaginatedPostsByCategoryByEngagement
func (s *PostService) GetPaginatedPostsByCategoryByEngagement(categoryID string, page, perPage int) ([]models.Post, int, error) {
	offset := (page - 1) * perPage

	// Count posts in category
	var totalCount int
	err := s.repo.Db.QueryRow(repository.CountCategoryPosts, categoryID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count category posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.CategoryByEngagement, categoryID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get category posts by engagement: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image *string
		var likes, dislikes, comments int

		// Fixed scan order to match SELECT order
		err := rows.Scan(&id, &userID, &title, &content, &createdAt, &hasImage, &image, &likes, &dislikes, &comments)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category post: %w", err)
		}

		// Handle NULL image field
		var imageFilename string
		if image != nil {
			imageFilename = *image
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      imageFilename,
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
	var totalCount int
	err := s.repo.Db.QueryRow(repository.CountCategoryPosts, categoryID).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count category posts: %w", err)
	}

	rows, err := s.repo.Db.Query(repository.CategoryNewest, categoryID, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get newest category posts: %w", err)
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var id, userID, title, content string
		var createdAt time.Time
		var hasImage bool
		var image *string

		// Scan order matches SELECT order
		err := rows.Scan(&id, &userID, &title, &content, &createdAt, &hasImage, &image)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan newest category post: %w", err)
		}

		// Handle NULL image field
		var imageFilename string
		if image != nil {
			imageFilename = *image
		}

		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"title":      title,
			"content":    content,
			"created_at": createdAt,
			"has_image":  hasImage,
			"image":      imageFilename,
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
		categories, err := s.category.GetCategoriesForPost(post["id"].(string))
		if err != nil {
			return nil, err
		}

		// Get reactions
		likes, dislikes, err := s.reaction.GetReactionsForPost(post["id"].(string))
		if err != nil {
			return nil, err
		}

		// Get comment count
		comments, err := s.repo.Get("comment", "post_id", post["id"].(string))
		if err != nil {
			return nil, err
		}

		commentsStruct, err := s.comment.ConvertToCommentStruct(comments)
		if err != nil {
			return nil, err
		}

		username := user[0]["username"].(string)
		// Check if post has an image
		hasImage := false
		var filepath string

		if hasImageField, ok := post["has_image"].(bool); ok && hasImageField {
			hasImage = true
			// Get image filename from database
			if post["image"] != nil {
				filepath = post["image"].(string)
			}
		} else if hasImageField, ok := post["has_image"].(int64); ok && hasImageField == 1 {
			hasImage = true
			// Get image filename from database
			if post["image"] != nil {
				filepath = post["image"].(string)
			}
		}

		result = append(result, models.Post{
			ID:         post["id"].(string),
			Username:   username,
			Categories: categories,
			Title:      post["title"].(string),
			Content:    post["content"].(string),
			CreatedAt:  post["created_at"].(time.Time),
			Likes:      likes,
			Dislikes:   dislikes,
			Comments:   commentsStruct,
			HasImage:   hasImage,
			Image:      filepath,
		})
	}

	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// CreatePost creates a new post with associated categories
func (s *PostService) CreatePost(data models.PostData) (string, error) {

	var err error
	var filepath string
	// if we have an image we will save it localy
	if data.HasImage {
		filepath, err = image.SaveImageLocally(s.config.UploadDir, data.ImageHeader, data.ImageFile)
		if err != nil {
			return "", err
		}
	}

	// Insert post
	postData := map[string]interface{}{
		"user_id":   data.UserID,
		"title":     data.Title,
		"content":   data.Content,
		"has_image": data.HasImage,
		"image":     filepath,
	}

	// Create the post record
	postID, err := s.repo.CreateRecord("post", postData)
	if err != nil {
		return "", err
	}

	// Associate categories with post
	for _, categoryID := range data.CategoryIDs {
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

// UpdatePost updates an existing post with new title, content, categories, and image status
func (s *PostService) UpdatePost(postID, userID, title, content string, categoryIDs []string, hasImage bool, image string) error {
	// First, verify that the user owns this post
	posts, err := s.repo.Get("post", "id", postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}
	if len(posts) == 0 {
		return fmt.Errorf("post not found")
	}

	if posts[0]["user_id"].(string) != userID {
		return fmt.Errorf("user does not own this post")
	}
	var postData map[string]interface{}
	if image == "" {
		postData = map[string]interface{}{
			"title":     title,
			"content":   content,
			"has_image": false,
			"image":     "",
		}
	} else {
		// Update the post record
		postData = map[string]interface{}{
			"title":     title,
			"content":   content,
			"has_image": true,
			"image":     image,
		}
	}

	err = s.repo.UpdateRecord("post", postID, postData)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	// Delete existing post-category relationships
	err = s.repo.DeleteJunctionByColumn("post_category", "post_id", postID)
	if err != nil {
		return fmt.Errorf("failed to delete existing post categories: %w", err)
	}

	// Insert new post-category relationships
	for _, categoryID := range categoryIDs {
		junctionData := map[string]interface{}{
			"post_id":     postID,
			"category_id": categoryID,
		}

		err = s.repo.CreateJunction("post_category", junctionData)
		if err != nil {
			return fmt.Errorf("failed to insert post category: %w", err)
		}
	}

	return nil
}

// DeletePost deletes a post
func (s *PostService) DeletePost(postID string, userID string) (error, int) {
	// check if user owns this post
	postInfo, err := s.repo.Get("post", "id", postID)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	post := postInfo[0]

	if post["user_id"].(string) != userID {
		return errors.New("This is not your post to delete :)"), 401
	}

	//delete the reactions of the post
	postReactions, err := s.repo.Get("reaction", "post_id", postID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	for _, postReaction := range postReactions {
		err := s.repo.DeleteRecord("reaction", postReaction["id"].(string))
		if err != nil {
			return err, http.StatusInternalServerError
		}
	}

	//delete comment reactions
	//and then the comments it self
	comments, err := s.repo.Get("comment", "post_id", postID)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	for _, comment := range comments {
		commentReactions, err := s.repo.Get("reaction", "comment_id", comment["id"].(string))
		if err != nil {
			return err, http.StatusInternalServerError
		}
		//comment reaction
		for _, commentReaction := range commentReactions {
			err := s.repo.DeleteRecord("reaction", commentReaction["id"].(string))
			if err != nil {
				return err, http.StatusInternalServerError
			}
		}
		//comment
		err = s.repo.DeleteRecord("comment", comment["id"].(string))
		if err != nil {
			return err, http.StatusInternalServerError
		}
	}

	// remove post row from post table
	err = s.repo.DeleteRecord("post", post["id"].(string))
	if err != nil {
		return err, http.StatusInternalServerError
	}

	// delete junction rows associated to this post
	err = s.repo.DeleteJunctionByColumn("post_category", "post_id", postID)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	if post["has_image"].(bool) == true {
		// remove locally
		err = image.RemoveImageLocaly("./static/files/uploads", post["image"].(string))
		if err != nil {
			return err, http.StatusInternalServerError
		}
	}

	return nil, http.StatusOK
}
