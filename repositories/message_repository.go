package repositories

import (
	"github.com/dennis/hello_go/models"
	"strconv"
)

type MessageRepository struct {
	messages []models.Message
	sequence uint64
}

func (r *MessageRepository) nextID() string {
	r.sequence += 1
	return strconv.FormatUint(r.sequence, 10)
}

func (r *MessageRepository) Insert(message models.Message) string {
	message.ID = r.nextID()
	r.messages = append(r.messages, message)

	return message.ID
}

func (r *MessageRepository) GetAll() []models.Message {
	return r.messages
}

func (r *MessageRepository) FindByID(id string) *models.Message {
	for _, message := range r.messages {
		if message.ID == id {
			return &message
		}
	}

	return nil
}

func (r *MessageRepository) Update(message models.Message) {
	r.DeleteByID(message.ID)
	r.messages = append(r.messages, message)
}

func (r *MessageRepository) DeleteByID(id string) {
	for index, message := range r.messages {
		if message.ID == id {
			r.messages = append(r.messages[:index], r.messages[index+1:]...)
		}
	}
}
