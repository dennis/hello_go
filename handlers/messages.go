package handlers

import (
	"net/http"
	"encoding/json"

	"github.com/dennis/hello_go/context"
	"github.com/dennis/hello_go/models"
	"github.com/gorilla/mux"
)

func GetMessages(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ctx.Messages)
}

func GetMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, message := range ctx.Messages {
		if message.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func CreateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message models.Message

	if err := json.NewDecoder(r.Body).Decode(&message); err == nil {
		ctx.Messages = append(ctx.Messages, message)

		json.NewEncoder(w).Encode(message)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func UpdateMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for index, message := range ctx.Messages {
		if message.ID == params["id"] {
			// Remove the old message
			ctx.Messages = append(ctx.Messages[:index], ctx.Messages[index+1:]...)

			// Add new
			var message models.Message
			_ = json.NewDecoder(r.Body).Decode(&message)

			ctx.Messages = append(ctx.Messages, message)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func DeleteMessage(ctx *context.Context, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for index, message := range ctx.Messages {
		if message.ID == params["id"] {
			ctx.Messages = append(ctx.Messages[:index], ctx.Messages[index+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

