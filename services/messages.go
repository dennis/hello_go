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
	message := s.MessageRepository.FindByID(id)

	if message == nil {
		return nil, &NotFoundError{}
	}

	return message, nil
}

func (s *MessageService) CreateMessage(message models.Message, user models.User) (*models.Message, error) {
	errors := message.Validate()

	if len(errors) > 0 {
		return nil, &NotValidError{Errors: errors}
	}

	message.Author = user.Username

	id := s.MessageRepository.Insert(message)

	return s.MessageRepository.FindByID(id), nil
}

func (s *MessageService) UpdateMessage(message models.Message, user models.User) (*models.Message, error) {
	storedMessage := s.MessageRepository.FindByID(message.ID)

	if storedMessage == nil {
		return nil, &NotFoundError{}
	}

	if storedMessage.Author != user.Username {
		return nil, &NotOwnerError{}
	}

	message.Author = user.Username

	if errors := message.Validate(); len(errors) == 0 {
		s.MessageRepository.Update(message)

		return s.MessageRepository.FindByID(message.ID), nil
	} else {
		return nil, &NotValidError{Errors: errors}
	}
}

func (s *MessageService) DeleteMessage(id string, user models.User) error {
	message := s.MessageRepository.FindByID(id)

	if message == nil {
		return &NotFoundError{}
	}

	if message.Author != user.Username {
		return &NotOwnerError{}
	}

	s.MessageRepository.DeleteByID(id)

	return nil
}
