package channels

import (
	"chat/server/models/users"
	"encoding/json"
	"errors"
	"time"
)

//Channel contains channel information ready to export to JSON format.
type Channel struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Private     bool          `json:"private"`
	Members     []*users.User `json:"members"`
	CreatedAt   time.Time     `json:"createdat"`
	Creator     *users.User   `json:"creator"`
	EditedAt    time.Time     `json:"editedat"`
}

//ChannelEdit Structs are used to edit channel rows.
type ChannelEdit struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

//ToJSON returns channel to JSON
func (c *Channel) ToJSON() ([]byte, error) {
	buffer, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

//Validate verifies if channel is valid
func (c *Channel) Validate() error {
	if len(c.Name) < 2 {
		return errors.New("Invalid name length")
	}

	return nil
}

//ContainsUser checks to see if channel contains user
func (c *Channel) ContainsUser(id int64) bool {
	for i := 0; i < len(c.Members); i++ {
		if c.Members[i].ID == id {
			return true
		}
	}

	return false
}

//GetMemberIDs gets all the members from a channel
func (c *Channel) GetMemberIDs() ([]int64, error) {
	if c.Private {
		members := c.Members
		var idslice []int64
		for i := 0; i < len(members); i++ {
			currentID := members[i].ID
			idslice = append(idslice, currentID)
		}
		return idslice, nil
	}

	return nil, nil
}
