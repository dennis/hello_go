package context

import (
	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/services"
)

// Struct are provided to handlers for each request containing the
// authenticated users and some additional state for the handlers to use
type Context struct {
	CurrentUser           models.User
	MessageService        services.MessageService
	AuthenticationService services.AuthenticationService
}
