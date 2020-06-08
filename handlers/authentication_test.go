package handlers

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/repositories"
	"github.com/dennis/hello_go/services"
)

var dennis models.User = models.User{Username: "foo", AuthToken: "authtokendennis"}
var marianne models.User = models.User{Username: "bar", AuthToken: "authtokenmarianne"}

func setup(authHeader string) (*context.Context, *http.Request) {
	userRepository := repositories.UserRepository{}

	userRepository.Insert(dennis)
	userRepository.Insert(marianne)

	messageRepository := repositories.MessageRepository{}

	ctx := &context.Context{
		AuthenticationService: services.AuthenticationService{UserRepository: &userRepository},
		MessageService:        services.MessageService{MessageRepository: &messageRepository},
	}

	r := httptest.NewRequest("GET", "/this/doesnt/matter", nil)

	if len(authHeader) > 0 {
		r.Header.Add("Authorization", authHeader)
	}

	return ctx, r
}

func base64Encode(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func TestAutenticate_ValidAuthentication(t *testing.T) {
	user := Authenticate(setup("Basic " + base64Encode("authtokendennis:")))

	if user == nil || *user != dennis {
		t.Errorf("Authentication expected to be successful for 'dennis'. Got %v", user)
	}
}

func TestAutenticate_InvalidScheme(t *testing.T) {
	user := Authenticate(setup("rot13 " + base64Encode("authtokendennis:")))

	if user != nil {
		t.Errorf("Authentication expected to fail, but got %v", user)
	}
}

func TestAutenticate_BadEncoding(t *testing.T) {
	user := Authenticate(setup("Basic " + base64Encode("authtokendennis:") + "NOPE"))

	if user != nil {
		t.Errorf("Authentication expected to fail, but got %v", user)
	}
}

func TestAutenticate_InvalidString(t *testing.T) {
	user := Authenticate(setup("Basic " + base64Encode("this-is-not-valid")))

	if user != nil {
		t.Errorf("Authentication expected to fail, but got %v", user)
	}
}

func TestAutenticate_IncorrectToken(t *testing.T) {
	user := Authenticate(setup("Basic " + base64Encode("badtoken:")))

	if user != nil {
		t.Errorf("Authentication expected to fail, but got %v", user)
	}
}
