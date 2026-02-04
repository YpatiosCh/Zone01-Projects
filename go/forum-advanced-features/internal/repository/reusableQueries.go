package repository

const (
	CountAllPosts = "SELECT COUNT(*) FROM post"

	CountUserPosts = "SELECT COUNT(*) FROM post WHERE user_id = ?"

	CountUserLikedPosts = `
		SELECT COUNT(DISTINCT p.id)
		FROM post p
		JOIN reaction r ON p.id = r.post_id
		WHERE r.user_id = ? AND r.type = 'like' AND r.post_id IS NOT NULL
	`

	CountUserDislikedPosts = `
		SELECT COUNT(DISTINCT p.id)
		FROM post p
		JOIN reaction r ON p.id = r.post_id
		WHERE r.user_id = ? AND r.type = 'dislike' AND r.post_id IS NOT NULL
	`

	CountUserCommentedPosts = `
		SELECT COUNT(DISTINCT p.id)
		FROM post p
		JOIN comment c ON p.id = c.post_id
		WHERE c.user_id = ?
	`

	CountCategoryPosts = `
		SELECT COUNT(*)
		FROM post p
		JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
	`

	AllByEngagement = `
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.created_at,
			p.has_image,
			p.image,
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

	CategoryByEngagement = `
		SELECT 
			p.id, 
			p.user_id, 
			p.title, 
			p.content, 
			p.created_at,
			p.has_image,
			p.image,
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

	AllNewest = `
		SELECT id, user_id, title, content, created_at, has_image, image
		FROM post
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	UserPosts = `
		SELECT id, user_id, title, content, created_at, has_image, image
		FROM post
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	UserLikedPosts = `
		SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.created_at, p.has_image, p.image
		FROM post p
		JOIN reaction r ON p.id = r.post_id
		WHERE r.user_id = ? AND r.type = 'like' AND r.post_id IS NOT NULL
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	UserDislikedPosts = `
		SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.created_at, p.has_image, p.image
		FROM post p
		JOIN reaction r ON p.id = r.post_id
		WHERE r.user_id = ? AND r.type = 'dislike' AND r.post_id IS NOT NULL
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	UserCommentedPosts = `
		SELECT DISTINCT p.id, p.user_id, p.title, p.content, p.created_at, p.has_image, p.image
		FROM post p
		JOIN comment c ON p.id = c.post_id
		WHERE c.user_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	CategoryNewest = `
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.has_image, p.image
		FROM post p
		JOIN post_category pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`
)
