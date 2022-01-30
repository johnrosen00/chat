package messages

import (
	"chat/server/gateway/models/users"
	"encoding/json"
	"errors"
	"time"
)

//Message : A struct that contains a message for export to JSON.
type Message struct {
	ID        int64       `json:"id"`
	ChannelID int64       `json:"channelid"`
	Body      string      `json:"body"`
	CreatedAt time.Time   `json:"createdat"`
	Creator   *users.User `json:"creator"`
	EditedAt  time.Time   `json:"editedat"`
}

//MessageEdit : a struct that contains info necessary to update message
type MessageEdit struct {
	Body string `json:"body"`
}

//Messages is an array of message addresses
type Messages []*Message

//ToJSON marshals a message struct into JSON
func (m *Message) ToJSON() ([]byte, error) {
	buffer, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

//Validate validates a message
func (m *Message) Validate() error {
	if len(m.Body) < 1 {
		return errors.New("Body must contain content")
	}

	return nil
}
