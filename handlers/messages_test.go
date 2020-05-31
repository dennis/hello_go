package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/handlers"
	"github.com/dennis/hello_go/models"
)

var message1 models.Message = models.Message{ID: "1", Author: "foo", Topic: "Topic1", Body: "Body1"}
var message2 models.Message = models.Message{ID: "2", Author: "bar", Topic: "Topic2", Body: "Body2"}
var fooUser models.User = models.User{Username: "foo"}
var barUser models.User = models.User{Username: "bar"}

func setupContext() *context.Context {
	ctx := context.Context{}

	ctx.CurrentUser = fooUser

	ctx.MessageRepository.Insert(message1)
	ctx.MessageRepository.Insert(message2)

	return &ctx
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

func TestMessages(t *testing.T) {
	ctx := setupContext()

	r := httptest.NewRequest("GET", "/api/messages", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{}

	handlers.GetMessages(ctx, w, r, vars)

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
	ctx := setupContext()

	r := httptest.NewRequest("GET", "/api/messages/1", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}

	handlers.GetMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertContentType(t, resp, "application/json")

	var message models.Message

	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		t.Errorf("Error decoding json-response: %v", err)
	}
}

func TestGetMessage_WhenMessageDoesNotExists(t *testing.T) {
	ctx := setupContext()

	r := httptest.NewRequest("GET", "/api/messages/28", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "28",
	}

	handlers.GetMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 404)
}

func assertMessage(t *testing.T, message *models.Message, expectedAuthor, expectedTopic, expectedBody, expectedID string) {
	assertEqual(t, message.Author, expectedAuthor, "Author is correct")
	assertEqual(t, message.Topic, expectedTopic, "Topic is correct")
	assertEqual(t, message.Body, expectedBody, "Body is correct")
	assertEqual(t, message.ID, expectedID, "Id is correct")
}

func TestCreateMessage_WithValidJson(t *testing.T) {
	ctx := setupContext()
	body := strings.NewReader("{\"id\":\"42\",\"author\":\"phony\", \"topic\":\"topic\", \"body\":\"body\"}")

	r := httptest.NewRequest("POST", "/api/messages", body)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "28",
	}

	handlers.CreateMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertContentType(t, resp, "application/json")

	// check response
	message := assertMessageJSON(t, resp)
	assertMessage(t, message, ctx.CurrentUser.Username, "topic", "body", "42")

	// Check repository
	storedMessage := ctx.MessageRepository.FindByID("42")
	assertMessage(t, storedMessage, ctx.CurrentUser.Username, "topic", "body", "42")
}

func TestCreateMessage_WithInvalidJson(t *testing.T) {
	ctx := setupContext()

	r := httptest.NewRequest("POST", "/api/messages", strings.NewReader(""))
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "28",
	}

	handlers.CreateMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 400)
	assertEmptyBody(t, resp)
}

func TestUpdateMessage_OwnerUpdatesMessage(t *testing.T) {
	ctx := setupContext()

	body := strings.NewReader("{\"id\":\"1\",\"author\":\"phony\", \"topic\":\"modified topic\", \"body\":\"modified body\"}")

	r := httptest.NewRequest("PUT", "/api/messages/1", body)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}

	handlers.UpdateMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertContentType(t, resp, "application/json")

	// check response
	message := assertMessageJSON(t, resp)
	assertMessage(t, message, ctx.CurrentUser.Username, "modified topic", "modified body", "1")

	// Check repository
	storedMessage := ctx.MessageRepository.FindByID("1")
	assertMessage(t, storedMessage, ctx.CurrentUser.Username, "modified topic", "modified body", "1")
}

func TestUpdateMessage_NonexistantMessage(t *testing.T) {
	ctx := setupContext()

	body := strings.NewReader("{\"id\":\"666\",\"author\":\"phony\", \"topic\":\"modified topic\", \"body\":\"modified body\"}")

	r := httptest.NewRequest("PUT", "/api/messages/666", body)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "666",
	}

	handlers.UpdateMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 404)
	assertEmptyBody(t, resp)
}

func TestUpdateMessage_OtherUserUpdatesMessage(t *testing.T) {
	ctx := setupContext()

	body := strings.NewReader("{\"id\":\"2\",\"author\":\"phony\", \"topic\":\"modified topic\", \"body\":\"modified body\"}")

	r := httptest.NewRequest("PUT", "/api/messages/2", body)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "2",
	}

	handlers.UpdateMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 401)
	assertEmptyBody(t, resp)

	// Verify that it isnt modified!
	storedMessage := ctx.MessageRepository.FindByID("2")
	assertMessage(t, storedMessage, "bar", "Topic2", "Body2", "2")
}

func TestDeleteMessage_OwnerDeletesMessage(t *testing.T) {
	ctx := setupContext()

	r := httptest.NewRequest("DELETE", "/api/messages/1", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "1",
	}

	handlers.DeleteMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 200)
	assertEmptyBody(t, resp)

	// Check if it was removed from repository
	if ctx.MessageRepository.FindByID("1") != nil {
		t.Errorf("Deleted Message still exists in Repository!")
	}
}

func TestDeleteMessage_NonexistantMessage(t *testing.T) {
	ctx := setupContext()

	r := httptest.NewRequest("DELETE", "/api/messages/666", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "666",
	}

	handlers.DeleteMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 404)
	assertEmptyBody(t, resp)
}

func TestDeleteMessage_OtherUserDeletesMessage(t *testing.T) {
	ctx := setupContext()

	r := httptest.NewRequest("DELETE", "/api/messages/2", nil)
	w := httptest.NewRecorder()
	vars := map[string]string{
		"id": "2",
	}

	handlers.DeleteMessage(ctx, w, r, vars)

	resp := w.Result()

	assertStatusCode(t, resp, 401)
	assertEmptyBody(t, resp)

	// Check if it was removed from repository
	if ctx.MessageRepository.FindByID("1") == nil {
		t.Errorf("Message was unexpectedly removed from Repository!")
	}
}
