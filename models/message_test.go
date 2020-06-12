package models

import (
	"testing"
)

func TestValidMessage(t *testing.T) {
	m := Message{
		ID: "ID",
		Topic: "Topic",
		Body: "Body",
		Author: "Author",
	}

	if err := m.Validate(); len(err) > 0 {
		t.Errorf("Expected message to be valid, but got errors: %v", err)
	}
}

func TestMissingTopic(t *testing.T) {
	m := Message{
		ID: "ID",
		Topic: "",
		Body: "Body",
		Author: "Author",
	}

	err := m.Validate()

	if len(err) != 0 && err[0] != "Topic is mandatory" {
		t.Errorf("Expected validation to fail with 'Topic is mandatory', but got: %v", err)
	}
}

func TestMissingBody(t *testing.T) {
	m := Message{
		ID: "ID",
		Topic: "Topic",
		Body: "",
		Author: "Author",
	}

	err := m.Validate()

	if len(err) != 0 && err[0] != "Body is mandatory" {
		t.Errorf("Expected validation to fail with 'Body is mandatory', but got: %v", err)
	}
}

func TestValidMessageWithMissingAuthor(t *testing.T) {
	m := Message{
		ID: "ID",
		Topic: "Topic",
		Body: "Body",
		Author: "",
	}

	if err := m.Validate(); len(err) > 0 {
		t.Errorf("Expected message to be valid, but got errors: %v", err)
	}
}

func TestValidMessageWithMissingID(t *testing.T) {
	m := Message{
		ID: "",
		Topic: "Topic",
		Body: "Body",
		Author: "Author",
	}

	if err := m.Validate(); len(err) > 0 {
		t.Errorf("Expected message to be valid, but got errors: %v", err)
	}
}

