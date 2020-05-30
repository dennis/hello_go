package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Message struct {
	ID    string `json:"id"`
	Topic string `json:"topic"`
	Body  string `json:"body"`
}

var messages []Message

func getMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(messages)
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, message := range messages {
		if message.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		messages = append(messages, message)

		json.NewEncoder(w).Encode(message)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for index, message := range messages {
		if message.ID == params["id"] {
			// Remove the old message
			messages = append(messages[:index], messages[index+1:]...)

			// Add new
			var message Message
			_ = json.NewDecoder(r.Body).Decode(&message)

			messages = append(messages, message)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(message)

			break
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteMessage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for index, message := range messages {
		if message.ID == params["id"] {
			messages = append(messages[:index], messages[index+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func main() {
	messages = append(messages, Message{ID: "1", Topic: "Hello World", Body: "Lorem lipsum"})
	messages = append(messages, Message{ID: "2", Topic: "re: Hello World", Body: "Really?"})

	// Router
	r := mux.NewRouter()

	r.HandleFunc("/api/messages", getMessages).Methods("GET")
	r.HandleFunc("/api/messages/{id}", getMessage).Methods("GET")
	r.HandleFunc("/api/messages", createMessage).Methods("POST")
	r.HandleFunc("/api/messages/{id}", updateMessage).Methods("PUT")
	r.HandleFunc("/api/messages/{id}", deleteMessage).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
