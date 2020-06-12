package repositories

import (
	"github.com/dennis/hello_go/models"
	"sync"
)

type UserRepository struct {
	users []models.User
	sync.Mutex
}

func (r *UserRepository) Insert(user models.User) {
	r.Lock()
	defer r.Unlock()
	r.users = append(r.users, user)
}

func (r *UserRepository) FindByToken(token string) *models.User {
	r.Lock()
	defer r.Unlock()
	for _, user := range r.users {
		if user.AuthToken == token {
			return &user
		}
	}

	return nil
}
