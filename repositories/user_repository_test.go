package repositories

import (
	"testing"

	"github.com/dennis/hello_go/repositories"
	"github.com/dennis/hello_go/models"
)

func TestFindingAUserThroughToken(t *testing.T) {
	repo := repositories.UserRepository {}

	u := models.User { Username: "username", AuthToken: "token" }

	repo.Insert(u)

	if f := repo.FindByToken("token"); f == nil {
		t.Error("Expected to find user by authtoken, but didnt")
	}
}

func TestAnInvalidToken(t *testing.T) {
	repo := repositories.UserRepository {}

	if f := repo.FindByToken("token"); f != nil {
		t.Errorf("Expected to find no user, but got: %v", f)
	}
}
