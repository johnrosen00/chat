package db

import (
	"chat/server/models/messages"

	"errors"
	"fmt"
	"time"
)

//MessageStore is a store of messages
type MessageStore struct {
	conn *Connection
}

func (c *Connection) InitMessageStore() *MessageStore {
	return &MessageStore{conn: c}
}

//GetRecentX gets X most recent messages in channel
func (store *MessageStore) GetRecentX(x int64, channelid int64) (messages.Messages, error) {
	db := store.conn.db
	q := "select messageid, channelid, createdat from messages where channelid=? order by createdat limit ?"
	rows, err := db.Query(q, channelid, x)

	if err != nil {
		return nil, err
	}

	var messageSlice messages.Messages
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
func (store *MessageStore) GetRecentXBeforeY(x int64, y int64, channelid int64) (messages.Messages, error) {
	db := store.conn.db
	q := "select messageid, channelid, createdat from messages where channelid=? and messageid < ? order by createdat limit ?"
	rows, err := db.Query(q, channelid, y, x)

	if err != nil {
		return nil, err
	}

	var messageSlice messages.Messages
	var currentMessageID int64
	for rows.Next() {

		if err := rows.Scan(currentMessageID); err != nil {
			return nil, err
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
func (store *MessageStore) GetByID(id int64) (*messages.Message, error) {
	db := store.conn.db
	q := "select messageid, channelid, messagebody, createdat, creatorid, editedat from messages where messageid=?"
	row := db.QueryRow(q, id)

	m := &messages.Message{}
	var userID int64
	err := row.Scan(m.ID, m.ChannelID, m.Body, m.CreatedAt, userID, m.EditedAt)

	if err != nil {
		return nil, err
	}

	userStore := store.conn.InitUserStore()
	user, err1 := userStore.GetByID(userID)
	if err1 != nil {
		return nil, err1
	}

	m.Creator = user

	return m, nil
}

//Insert inserts a new message
func (store *MessageStore) Insert(m *messages.Message) (*messages.Message, error) {
	db := store.conn.db
	if err := m.Validate(); err != nil {
		return nil, err
	}

	insq := "insert into messages(channelid, body, createdat, creator) values(?,?,?,?)"

	res, err1 := db.Exec(insq, m.ChannelID, m.Body, time.Now(), m.Creator.ID)
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
func (store *MessageStore) Delete(id int64) error {
	db := store.conn.db
	ex := "delete from messages where messageid = ?"

	_, err := db.Exec(ex, id)

	if err != nil {
		return err
	}

	return nil
}

//Edit modifies an existing message
func (store *MessageStore) Edit(id int64, edit *messages.MessageEdit) (*messages.Message, error) {
	db := store.conn.db
	queryF := "update messages set body = ? , editedAt = ? where messageid = ?"

	if len(edit.Body) < 1 {
		return nil, errors.New("submit valid edit")
	}

	if _, err := db.Exec(queryF, edit.Body, time.Now(), id); err != nil {
		return nil, err
	}

	return store.GetByID(id)
}
