package models

type Message struct {
	ID     string `json:"id"`
	Topic  string `json:"topic"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

func (m *Message) Validate() []string {
	errors := make([]string, 0)

	if len(m.Topic) == 0 {
		errors = append(errors, "Topic is mandatory")
	}
	if len(m.Body) == 0 {
		errors = append(errors, "Body is mandatory")
	}

	return errors
}
