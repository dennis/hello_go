package handlers

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
)


func findUser(ctx *context.Context, token string) *models.User {
	for _, user := range ctx.Users {
		if user.AuthToken == token {
			return &user
		}
	}

	return nil
}

func Authenticate(ctx *context.Context, r *http.Request) *models.User {
	const basicScheme string = "Basic "

	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, basicScheme) {
		return nil
	}

	str, err := base64.StdEncoding.DecodeString(auth[len(basicScheme):])
	if err != nil {
		return nil
	}

	username_password := bytes.SplitN(str, []byte(":"), 2)

	if len(username_password) != 2 {
		return nil
	}

	if len(username_password[1]) > 0 {
		return nil
	}

	username := string(username_password[0])

	return findUser(ctx, username)
}

