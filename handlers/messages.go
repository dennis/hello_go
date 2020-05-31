package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
	"github.com/gorilla/mux"
)

func GetMessages(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages := ctx.MessageRepository.GetAll()

	json.NewEncoder(w).Encode(messages)
}

func GetMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if message := ctx.MessageRepository.FindByID(params["id"]); message != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(message)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func CreateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message models.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err == nil {
		message.Author = ctx.CurrentUser.Username

		ctx.MessageRepository.Insert(message)

		json.NewEncoder(w).Encode(message)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func UpdateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

	if message := ctx.MessageRepository.FindByID(id); message != nil {
		if message.Author == ctx.CurrentUser.Username {
			ctx.MessageRepository.DeleteByID(id)

			var message models.Message
			_ = json.NewDecoder(r.Body).Decode(&message)

			message.Author = ctx.CurrentUser.Username

			ctx.MessageRepository.Update(message)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func DeleteMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id := params["id"]

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
