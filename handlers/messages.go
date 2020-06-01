package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
)

// Returns a JSON array with all the available Messages
// returns:
//   200 success: if successful
func GetMessages(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	messages := ctx.MessageRepository.GetAll()

	json.NewEncoder(w).Encode(messages)
}

// Returns a specific message as json
// returns:
//   200 success: if successful
//   404 bad request: if message was not found
func GetMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	if message := ctx.MessageRepository.FindByID(vars["id"]); message != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Creates a new Message. Will force Author to be CurrentUser. ID is assigned by
// service.  The response will contain the message as JSON
// returns:
//   200 success: if message was successful created
//   400 bad request: in case of errors (reading the json)
//   422 unprocessable entity: if provided JSON isn't valid
func CreateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var message models.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err == nil {
		message.Author = ctx.CurrentUser.Username

		if errors := message.Validate(); len(errors) == 0 {
			id := ctx.MessageRepository.Insert(message)

			storedMessage := ctx.MessageRepository.FindByID(id)

			json.NewEncoder(w).Encode(storedMessage)
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(errors)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
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

	if message := ctx.MessageRepository.FindByID(id); message != nil {
		if message.Author == ctx.CurrentUser.Username {
			var message models.Message
			if err := json.NewDecoder(r.Body).Decode(&message); err == nil {
				message.ID = id
				message.Author = ctx.CurrentUser.Username

				if errors := message.Validate(); len(errors) == 0 {
					ctx.MessageRepository.Update(message)

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(message)
				} else {
					w.WriteHeader(http.StatusUnprocessableEntity)
					json.NewEncoder(w).Encode(errors)
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Deletes a Message. 
// returns:
//   200 success: if message was successful updated
//   401 unauthorized: if CurrentUser isn't the owner of the message
//   404 not found: if message wasn't found
func DeleteMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	id := vars["id"]

	if message := ctx.MessageRepository.FindByID(id); message != nil {
		if message.Author == ctx.CurrentUser.Username {
			ctx.MessageRepository.DeleteByID(id)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
