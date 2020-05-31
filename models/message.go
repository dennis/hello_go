package models

type Message struct {
	ID     string `json:"id"`
	Topic  string `json:"topic"`
	Body   string `json:"body"`
	Author string `json:"author"`
}
