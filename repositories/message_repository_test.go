package repositories

import (
	"github.com/dennis/hello_go/models"
	"testing"
)

func TestInsertingMessagesAssignsIDToThem(t *testing.T) {
	repo := MessageRepository{}
	m1 := models.Message{}
	m2 := models.Message{}

	if n := repo.Insert(m1); n != "1" {
		t.Errorf("Message `m1` got unexpected ID: %v", n)
	} else if n := repo.Insert(m2); n != "2" {
		t.Errorf("Message `m2` got unexpected ID: %v", n)
	}
}

func TestGetAllReturnsNothingIfNothingIsAdded(t *testing.T) {
	repo := MessageRepository{}

	if r := repo.GetAll(); len(r) > 0 {
		t.Errorf("Expected GetAll() not to return any messages, but got %v", r)
	}
}

func TestGetAllReturnsAddedMessages(t *testing.T) {
	repo := MessageRepository{}
	m := models.Message{}
	repo.Insert(m)

	if r := repo.GetAll(); len(r) != 1 {
		t.Errorf("Expected GetAll() to return one message, but got %v", r)
	}
}

func TestModifyingAMessage(t *testing.T) {
	repo := MessageRepository{}

	m := models.Message{Body: "original"}
	id := repo.Insert(m)

	u := models.Message{ID: id, Body: "updated"}
	repo.Update(u)

	n := repo.FindByID(id)

	if n.Body != "updated" {
		t.Errorf("Updated message got unexpected content: %v", n.Body)
	}
}

func TestRemovingAMessage(t *testing.T) {
	repo := MessageRepository{}

	m := models.Message{}
	id := repo.Insert(m)

	repo.DeleteByID(id)

	if r := repo.GetAll(); len(r) > 0 {
		t.Errorf("Expected message to be deleted. Got: %v", r)
	}

}
