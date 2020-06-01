package context

import (
	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/repositories"
)

// Struct are provided to handlers for each request containing the
// authenticated users and some additional state for the handlers to use
type Context struct {
	CurrentUser       models.User
	MessageRepository repositories.MessageRepository
	UserRepository    repositories.UserRepository
}
