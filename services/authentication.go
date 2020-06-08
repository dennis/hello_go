package services

import (
	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/repositories"
)

type AuthenticationService struct {
	UserRepository *repositories.UserRepository
}

func (s *AuthenticationService) Authenticate(token string) *models.User {
	return s.UserRepository.FindByToken(token)
}
