package app

import (
	"encoding/json"
	"log"
	"os"

	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/repositories"
)

func PopulateMessages(r *repositories.MessageRepository) {
	file, err := os.Open("messages.json")

	if err != nil {
		if os.IsNotExist(err) {
			// Silently ignore this
			log.Println("no messages.json found - no data prepopulated")
			return
		} else {
			panic(err)
		}
		return
	}

	log.Println("Loading messages.json")

	defer file.Close()

	var messages []models.Message

	json_err := json.NewDecoder(file).Decode(&messages)

	if json_err != nil {
		log.Printf("Error parsing messages.json: %v", json_err)
		return
	}

	for _, m := range messages {
		r.Insert(m)
	}
}

func PopulateUsers(r *repositories.UserRepository) {
	file, err := os.Open("users.json")

	if err != nil {
		if os.IsNotExist(err) {
			// Silently ignore this
			log.Println("no users.json found - no data prepopulated")
			return
		} else {
			panic(err)
		}
		return
	}

	log.Println("Loading users.json")

	defer file.Close()

	var users []models.User

	json_err := json.NewDecoder(file).Decode(&users)

	if json_err != nil {
		log.Printf("Error parsing users.json: %v", json_err)
		return
	}

	for _, u := range users {
		r.Insert(u)
	}
}
