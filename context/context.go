package context

import (
	"github.com/dennis/hello_go/models"
)

type Context struct {
	Messages []models.Message
	Users []models.User
}
