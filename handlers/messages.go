package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
)

func GetMessages(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	messages := ctx.MessageRepository.GetAll()

	json.NewEncoder(w).Encode(messages)
}

func GetMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	if message := ctx.MessageRepository.FindByID(vars["id"]); message != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func CreateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	var message models.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err == nil {
		message.Author = ctx.CurrentUser.Username

		id := ctx.MessageRepository.Insert(message)

		storedMessage := ctx.MessageRepository.FindByID(id)

		json.NewEncoder(w).Encode(storedMessage)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func UpdateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) {
	id := vars["id"]

	if message := ctx.MessageRepository.FindByID(id); message != nil {
		if message.Author == ctx.CurrentUser.Username {
			var message models.Message
			if err := json.NewDecoder(r.Body).Decode(&message); err == nil {
				message.ID = id
				message.Author = ctx.CurrentUser.Username

				ctx.MessageRepository.Update(message)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(message)
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
