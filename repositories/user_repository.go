package repositories

import (
	"github.com/dennis/hello_go/models"
)

type UserRepository struct {
	users []models.User
}

func (r *UserRepository) Insert(user models.User) {
	r.users = append(r.users, user)
}

func (r *UserRepository) FindByToken(token string) *models.User {
	for _, user := range r.users {
		if user.AuthToken == token {
			return &user
		}
	}

	return nil
}
