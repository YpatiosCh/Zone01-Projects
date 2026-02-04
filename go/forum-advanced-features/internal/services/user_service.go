package services

import "forum/internal/repository"

type UserService struct {
	repo *repository.Manager
}

func NewUserService(repo *repository.Manager) *UserService {
	return &UserService{
		repo: repo,
	}
}

// UserExist checks if a user exists by username
func (u *UserService) UserExist(username string) bool {
	user, err := u.repo.Get("user", "username", username)
	if err != nil {
		return false
	}
	return len(user) > 0
}
