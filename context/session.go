package context

import (
	"github.com/dennis/hello_go/models"
)

type Session struct {
	CurrentUser           models.User
}

