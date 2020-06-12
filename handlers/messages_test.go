package handlers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/repositories"
	"github.com/dennis/hello_go/services"
)

var message1 models.Message = models.Message{ID: "1", Author: "foo", Topic: "Topic1", Body: "Body1"}
var message2 models.Message = models.Message{ID: "2", Author: "bar", Topic: "Topic2", Body: "Body2"}
var fooUser models.User = models.User{Username: "foo"}
var barUser models.User = models.User{Username: "bar"}

func setupContext() (*context.Context, *context.Session) {
	userRepository := repositories.UserRepository{}
	messageRepository := repositories.MessageRepository{}
	messageRepository.Insert(message1)
	messageRepository.Insert(message2)

	return &context.Context{
		AuthenticationService: services.AuthenticationService{UserRepository: &userRepository},
		MessageService:        services.MessageService{MessageRepository: &messageRepository},
	}, &context.Session { CurrentUser: fooUser }
}

func setupRequestWithContent(content io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	// Any information provided via the URL is handled in App. IDs
	// and similar stuff is provided via the vars argument to the handler.
	// So our handler should care about the URL or Verb
	r := httptest.NewRequest("GET", "/this/doesnt/matter", content)
	w := httptest.NewRecorder()

	return r, w
}

func setupRequest() (*http.Request, *httptest.ResponseRecorder) {
	return setupRequestWithContent(nil)
}

func assertStatusCode(t *testing.T, resp *http.Response, expected int) {
	actual := resp.StatusCode

	if expected != actual {
		t.Errorf("Status code mismatch: expected=%v, actual=%v", expected, actual)
	}
}

func assertContentType(t *testing.T, resp *http.Response, expected string) {
	actual := resp.Header.Get("Content-Type")

	if actual != expected {
		t.Errorf("Content-Type mismatch: expected=%v, actual=%v", expected, actual)
	}
}

func assertEmptyBody(t *testing.T, resp *http.Response) {
	body, _ := ioutil.ReadAll(resp.Body)

	if len(body) > 0 {
		t.Errorf("Unexpected body in response")
	}
}

func assertEqual(t *testing.T, actual, expected, message string) {
	if actual != expected {
		t.Errorf("Assertion '%s' not met: actual=%v, expected=%v", message, actual, expected)
	}
}

func assertMessageJSON(t *testing.T, resp *http.Response) *models.Message {
	var message models.Message

	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		t.Errorf("Error decoding json-response: %v", err)
		return nil
	}

	return &message
}

func assertMessage(t *testing.T, message *models.Message, expectedAuthor, expectedTopic, expectedBody, expectedID string) {
	assertEqual(t, message.Author, expectedAuthor, "Author is correct")
	assertEqual(t, message.Topic, expectedTopic, "Topic is correct")
	assertEqual(t, message.Body, expectedBody, "Body is correct")
	assertEqual(t, message.ID, expectedID, "Id is correct")
}

func includesString(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

func assertArrayContains(t *testing.T, haystack []string, needle, message string) {
	if !includesString(haystack, needle) {
		t.Errorf("Assertion '%s' not met", message)
	}
}

var noVars = map[string]string{}

func TestGetMessages(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequest()

	GetMessages(ctx, session, w, r, noVars)

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertContentType(t, resp, "application/json")

	// check content

	var messages []models.Message

	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		t.Errorf("Error decoding json-response: %v", err)
	}

	if len(messages) != 2 {
		t.Errorf("Got %v messages expected %v", len(messages), 2)
	} else {
		expected_messages :=
			(messages[0] == message1 || messages[0] == message2) &&
				(messages[1] == message1 || messages[1] == message2)

		if !expected_messages {
			t.Errorf("Unexpect JSON returned")
		}
	}
}

func TestGetMessage_WhenMessageExists(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequest()

	GetMessage(ctx, session, w, r, map[string]string{
		"id": "1",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertContentType(t, resp, "application/json")

	var message models.Message

	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		t.Errorf("Error decoding json-response: %v", err)
	}
}

func TestGetMessage_WhenMessageDoesNotExists(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequest()

	GetMessage(ctx, session, w, r, map[string]string{
		"id": "28",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 404)
}

func TestCreateMessage_WithCorrectData(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader("{\"id\":\"42\",\"author\":\"phony\", \"topic\":\"topic\", \"body\":\"body\"}"))

	CreateMessage(ctx, session, w, r, noVars)

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertContentType(t, resp, "application/json")

	// check response
	message := assertMessageJSON(t, resp)
	if message != nil {
		assertMessage(t, message, session.CurrentUser.Username, "topic", "body", message.ID)

		// Check repository
		storedMessage, err := ctx.MessageService.GetMessage(message.ID)
		if storedMessage != nil && err == nil {
			assertMessage(t, storedMessage, session.CurrentUser.Username, "topic", "body", message.ID)
		} else {
			t.Error("Message not found in repository")
		}
	}
}

func TestCreateMessage_WithMissingData(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader("{}"))

	CreateMessage(ctx, session, w, r, noVars)

	resp := w.Result()

	assertStatusCode(t, resp, 422)

	// content

	var errors []string

	if err := json.NewDecoder(resp.Body).Decode(&errors); err != nil {
		t.Errorf("Error decoding json-response: %v", err)
	} else {
		assertArrayContains(t, errors, "Topic is mandatory", "Errors contains 'Topic is mandatory'")
		assertArrayContains(t, errors, "Body is mandatory", "Errors contains 'Body is mandatory'")
	}
}

func TestCreateMessage_WithInvalidJson(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader(""))

	CreateMessage(ctx, session, w, r, noVars)

	resp := w.Result()

	assertStatusCode(t, resp, 400)
	assertEmptyBody(t, resp)
}

func TestUpdateMessage_OwnerUpdatesMessage(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader("{\"id\":\"1\",\"author\":\"phony\", \"topic\":\"modified topic\", \"body\":\"modified body\"}"))

	UpdateMessage(ctx, session, w, r, map[string]string{
		"id": "1",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertContentType(t, resp, "application/json")

	// check response
	message := assertMessageJSON(t, resp)
	if message != nil {
		assertMessage(t, message, session.CurrentUser.Username, "modified topic", "modified body", "1")

		// Check repository
		storedMessage, err := ctx.MessageService.GetMessage("1")
		if storedMessage != nil && err == nil {
			assertMessage(t, storedMessage, session.CurrentUser.Username, "modified topic", "modified body", "1")
		} else {
			t.Error("Message not found in repository")
		}
	}
}

func TestUpdateMessage_NonexistantMessage(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader("{\"id\":\"666\",\"author\":\"phony\", \"topic\":\"modified topic\", \"body\":\"modified body\"}"))

	UpdateMessage(ctx, session, w, r, map[string]string{
		"id": "666",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 404)
	assertEmptyBody(t, resp)
}

func TestUpdateMessage_WithMissingData(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader("{}"))

	UpdateMessage(ctx, session, w, r, map[string]string{
		"id": "1",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 422)

	// content

	var errors []string

	if err := json.NewDecoder(resp.Body).Decode(&errors); err != nil {
		t.Errorf("Error decoding json-response: %v", err)
	} else {
		assertArrayContains(t, errors, "Topic is mandatory", "Errors contains 'Topic is mandatory'")
		assertArrayContains(t, errors, "Body is mandatory", "Errors contains 'Body is mandatory'")
	}
}

func TestUpdateMessage_OtherUserUpdatesMessage(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader("{\"id\":\"2\",\"author\":\"phony\", \"topic\":\"modified topic\", \"body\":\"modified body\"}"))

	UpdateMessage(ctx, session, w, r, map[string]string{
		"id": "2",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 401)
	assertEmptyBody(t, resp)

	// Verify that it isnt modified!
	storedMessage, err := ctx.MessageService.GetMessage("2")
	if storedMessage != nil && err == nil {
		assertMessage(t, storedMessage, "bar", "Topic2", "Body2", "2")
	} else {
		t.Error("Message not found in repository")
	}
}

func TestUpdateMessage_WithInvalidJson(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequestWithContent(strings.NewReader(""))

	UpdateMessage(ctx, session, w, r, map[string]string{
		"id": "1",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 400)
	assertEmptyBody(t, resp)
}

func TestDeleteMessage_OwnerDeletesMessage(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequest()

	DeleteMessage(ctx, session, w, r, map[string]string{
		"id": "1",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertEmptyBody(t, resp)

	// Check if it was removed from repository
	if msg, _ := ctx.MessageService.GetMessage("1"); msg != nil {
		t.Errorf("Deleted Message still exists in Repository!")
	}
}

func TestDeleteMessage_NonexistantMessage(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequest()

	DeleteMessage(ctx, session, w, r, map[string]string{
		"id": "666",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 404)
	assertEmptyBody(t, resp)
}

func TestDeleteMessage_OtherUserDeletesMessage(t *testing.T) {
	ctx, session := setupContext()

	r, w := setupRequest()

	DeleteMessage(ctx, session, w, r, map[string]string{
		"id": "2",
	})

	resp := w.Result()

	assertStatusCode(t, resp, 401)
	assertEmptyBody(t, resp)

	// Check if it was removed from repository
	if msg, err := ctx.MessageService.GetMessage("1"); msg == nil && err == nil {
		t.Errorf("Message was unexpectedly removed from Repository!")
	}
}
