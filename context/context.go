package context

import (
	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/repositories"
)

type Context struct {
	Users             []models.User
	CurrentUser       models.User
	MessageRepository repositories.MessageRepository
}
