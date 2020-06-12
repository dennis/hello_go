package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/services"
)

func handleError(w http.ResponseWriter, err error) {
	if _, ok := err.(*services.NotFoundError); ok {
		w.WriteHeader(http.StatusNotFound)
	} else if serviceErr, ok := err.(*services.NotValidError); ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(serviceErr.Errors)
	} else if _, ok := err.(*services.NotOwnerError); ok {
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		// Catch all
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Returns a JSON array with all the available Messages
// returns:
//   200 success: if successful
func GetMessages(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	messages, err := ctx.MessageService.GetMessages()

	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// Returns a specific message as json
// returns:
//   200 success: if successful
//   404 bad request: if message was not found
func GetMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	message, err := ctx.MessageService.GetMessage(vars["id"])

	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// Creates a new Message. Will force Author to be CurrentUser. ID is assigned by
// service.  The response will contain the message as JSON
// returns:
//   200 success: if message was successful created
//   400 bad request: in case of errors (reading the json)
//   422 unprocessable entity: if provided JSON isn't valid
func CreateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	var message models.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		handleError(w, err)
		return
	}

	storedMessage, serviceError := ctx.MessageService.CreateMessage(message, ctx.CurrentUser)

	if serviceError != nil {
		handleError(w, serviceError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(storedMessage)
}

// Updates the Message.
// returns:
//   200 success: if message was successful updated
//   400 bad request: in case of errors (reading the json)
//   401 unauthorized: if CurrentUser isn't the owner of the message
//   404 not found: if message wasn't found
//   422 unprocessable entity: if provided JSON isn't valid
func UpdateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	id := vars["id"]

	var message models.Message

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {
		handleError(w, err)
		return
	}

	message.ID = id

	storedMessage, serviceError := ctx.MessageService.UpdateMessage(message, ctx.CurrentUser)
	
	if serviceError != nil {
		handleError(w, serviceError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storedMessage)
}

// Deletes a Message.
// returns:
//   200 success: if message was successful updated
//   401 unauthorized: if CurrentUser isn't the owner of the message
//   404 not found: if message wasn't found
func DeleteMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	err := ctx.MessageService.DeleteMessage(vars["id"], ctx.CurrentUser)
	
	if err != nil {
		handleError(w, err)
		return
	}
}
