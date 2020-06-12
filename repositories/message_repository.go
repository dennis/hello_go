package repositories

import (
	"github.com/dennis/hello_go/models"
	"strconv"
	"sync"
)

type MessageRepository struct {
	messages []models.Message
	sequence uint64
	sync.Mutex
}

func (r *MessageRepository) nextID() string {
	r.sequence += 1
	return strconv.FormatUint(r.sequence, 10)
}

func (r *MessageRepository) Insert(message models.Message) string {
	r.Lock()
	defer r.Unlock()
	message.ID = r.nextID()
	r.messages = append(r.messages, message)

	return message.ID
}

func (r *MessageRepository) GetAll() []models.Message {
	r.Lock()
	defer r.Unlock()

	messages := []models.Message{}

	for _, m := range r.messages {
		messages = append(messages, m)
	}

	return messages
}

func (r *MessageRepository) FindByID(id string) *models.Message {
	r.Lock()
	defer r.Unlock()
	for _, message := range r.messages {
		if message.ID == id {
			// return a copy of message
			d := message
			return &d
		}
	}

	return nil
}

func (r *MessageRepository) Update(message models.Message) {
	r.Lock()
	defer r.Unlock()
	r.deleteByIDWithoutLock(message.ID)
	r.messages = append(r.messages, message)
}

func (r *MessageRepository) deleteByIDWithoutLock(id string) {
	for index, message := range r.messages {
		if message.ID == id {
			r.messages = append(r.messages[:index], r.messages[index+1:]...)
		}
	}
}

func (r *MessageRepository) DeleteByID(id string) {
	r.Lock()
	defer r.Unlock()
	r.deleteByIDWithoutLock(id)
}
