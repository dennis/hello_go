package services

import (
	"strings"

	"github.com/dennis/hello_go/models"
	"github.com/dennis/hello_go/repositories"
)

type NotFoundError struct{}

func (e *NotFoundError) Error() string { return "Not found" }

type NotValidError struct {
	Errors []string
}

func (e *NotValidError) Error() string { return strings.Join(e.Errors, ". ") }

type NotOwnerError struct{}

func (e *NotOwnerError) Error() string { return "Not owner" }

type MessageService struct {
	MessageRepository *repositories.MessageRepository
}

func (s *MessageService) GetMessages() ([]models.Message, error) {
	return s.MessageRepository.GetAll(), nil
}

func (s *MessageService) GetMessage(id string) (*models.Message, error) {
	if message := s.MessageRepository.FindByID(id); message != nil {
		return message, nil
	} else {
		return nil, &NotFoundError{}
	}
}

func (s *MessageService) CreateMessage(message models.Message, user models.User) (*models.Message, error) {
	if errors := message.Validate(); len(errors) == 0 {
		message.Author = user.Username

		id := s.MessageRepository.Insert(message)

		return s.MessageRepository.FindByID(id), nil
	} else {
		return nil, &NotValidError{Errors: errors}
	}
}

func (s *MessageService) UpdateMessage(message models.Message, user models.User) (*models.Message, error) {
	if storedMessage := s.MessageRepository.FindByID(message.ID); storedMessage != nil {
		if storedMessage.Author == user.Username {
			message.Author = user.Username

			if errors := message.Validate(); len(errors) == 0 {
				s.MessageRepository.Update(message)

				return s.MessageRepository.FindByID(message.ID), nil
			} else {
				return nil, &NotValidError{Errors: errors}
			}
		} else {
			return nil, &NotOwnerError{}
		}
	} else {
		return nil, &NotFoundError{}
	}
}

func (s *MessageService) DeleteMessage(id string, user models.User) error {
	if message := s.MessageRepository.FindByID(id); message != nil {
		if message.Author == user.Username {
			s.MessageRepository.DeleteByID(id)

			return nil
		} else {
			return &NotOwnerError{}
		}

	} else {
		return &NotFoundError{}
	}
}
