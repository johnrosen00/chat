package messages

import (
	"chat/server/gateway/models/users"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

//SQLMessageStore is a store of messages
type SQLMessageStore struct {
	DB *sql.DB
}

//GetRecentX gets X most recent messages in channel
func (store *SQLMessageStore) GetRecentX(x int64, channelid int64) (Messages, error) {
	q := "select messageid, channelid, createdat from messages where channelid=? order by createdat limit ?"
	rows, err := store.DB.Query(q, channelid, x)

	if err != nil {
		return nil, err
	}

	var messageSlice Messages
	var currentMessageID int64
	for rows.Next() {

		if err89 := rows.Scan(currentMessageID); err89 != nil {
			return nil, err89
		}

		currentMessage, err45 := store.GetByID(currentMessageID)

		if err45 != nil {
			return nil, err45
		}

		messageSlice = append(messageSlice, currentMessage)
	}
	return messageSlice, nil
}

//GetRecentXBeforeY gets X most recent messages in channel BeforeY
func (store *SQLMessageStore) GetRecentXBeforeY(x int64, y int64, channelid int64) (Messages, error) {
	q := "select messageid, channelid, createdat from messages where channelid=? and messageid < ? order by createdat limit ?"
	rows, err := store.DB.Query(q, channelid, y, x)

	if err != nil {
		return nil, err
	}

	var messageSlice Messages
	var currentMessageID int64
	for rows.Next() {

		if err89 := rows.Scan(currentMessageID); err89 != nil {
			return nil, err89
		}

		currentMessage, err45 := store.GetByID(currentMessageID)

		if err45 != nil {
			return nil, err45
		}

		messageSlice = append(messageSlice, currentMessage)
	}
	return messageSlice, nil
}

//GetByID gets Message by ID
func (store *SQLMessageStore) GetByID(id int64) (*Message, error) {
	q := "select messageid, channelid, messagebody, createdat, creatorid, editedat from messages where messageid=?"
	row := store.DB.QueryRow(q, id)

	m := &Message{}
	var userID int64
	err := row.Scan(m.ID, m.ChannelID, m.Body, m.CreatedAt, userID, m.EditedAt)

	if err != nil {
		return nil, err
	}

	userStore := &users.MySQLStore{}
	userStore.DB = store.DB
	user, err1 := userStore.GetByID(userID)
	if err1 != nil {
		return nil, err1
	}

	m.Creator = user

	return m, nil
}

//Insert inserts a new message
func (store *SQLMessageStore) Insert(m *Message) (*Message, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	insq := "insert into messages(channelid, body, createdat, creator) values(?,?,?,?)"

	res, err1 := store.DB.Exec(insq, m.ChannelID, m.Body, time.Now(), m.Creator.ID)
	if err1 != nil {
		fmt.Printf("error inserting new row: %v\n", err1)
		return nil, err1
	}

	//get the auto-assigned ID for the new row
	id, err2 := res.LastInsertId()
	if err2 != nil {
		fmt.Printf("error getting new ID: %v\n", id)
		return nil, err2
	}

	m.ID = id
	return m, nil
}

//Delete deletes a channel from the store.
func (store *SQLMessageStore) Delete(id int64) error {
	ex := "delete from messages where messageid = ?"

	_, err := store.DB.Exec(ex, id)

	if err != nil {
		return err
	}

	return nil
}

//Edit modifies an existing message
func (store *SQLMessageStore) Edit(id int64, edit *MessageEdit) (*Message, error) {
	queryF := "update messages set body = ? , editedAt = ? where messageid = ?"

	if len(edit.Body) < 1 {
		return nil, errors.New("submit valid edit")
	}

	if _, err := store.DB.Exec(queryF, edit.Body, time.Now(), id); err != nil {
		return nil, err
	}

	return store.GetByID(id)
}
