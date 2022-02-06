package db

import (
	"chat/server/models/channels"
	"chat/server/models/messages"
	"chat/server/models/users"
	"errors"
	"fmt"
	"time"
)

type ChannelStore struct {
	conn *Connection
}

func (c *Connection) InitChannelStore() *ChannelStore {
	return &ChannelStore{conn: c}
}

//GetAllChannels gets all channels
func (store *ChannelStore) GetAllChannels() ([]*channels.Channel, error) {
	db := store.conn.db
	q := "select channelid from channels"
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	var channelSlice []*channels.Channel
	var currentChannelID int64
	for rows.Next() {

		if err89 := rows.Scan(currentChannelID); err89 != nil {
			return nil, err89
		}

		currentChannel, err45 := store.GetByID(currentChannelID)

		if err45 != nil {
			return nil, err45
		}

		channelSlice = append(channelSlice, currentChannel)
	}

	return channelSlice, nil
}

//Init adds the general channel
func (store *ChannelStore) Init() error {
	db := store.conn.db
	q := "insert ignore into channels (channelname, channeldescription, isprivate, createdat) values (general, general, false, ?)"
	res, err1 := db.Exec(q, time.Now())
	if err1 != nil {
		fmt.Printf("error inserting new row: %v\n", err1)
		return err1
	}

	id, err2 := res.LastInsertId()
	if err2 != nil {
		fmt.Printf("error getting new ID: %v\n", id)
		return err2
	}

	return nil
}

//Edit edits a channel.
func (store *ChannelStore) Edit(id int64, edit *channels.ChannelEdit) (*channels.Channel, error) {
	db := store.conn.db
	if len(edit.Name) != 0 {
		return nil, errors.New("enter a name over 0 char long")
	}
	queryF := "update channels set name = ? , description = ? , editedAt = ? where channelid = ?"

	if _, err := db.Exec(queryF, edit.Name, edit.Description, time.Now(), id); err != nil {
		return nil, err
	}
	return store.GetByID(id)
}

//GetByID returns a channel struct that can be encoded to JSON
func (store *ChannelStore) GetByID(id int64) (*channels.Channel, error) {
	db := store.conn.db
	q := "select channelid, channelname, channeldescription, isprivate, createdat, creatorid, editedat, from channels where channelid=?"
	row := db.QueryRow(q, id)

	c := &channels.Channel{}
	err := row.Scan(c.ID, c.Name, c.Description, c.Private, c.CreatedAt, c.Creator.ID, c.EditedAt)

	if err != nil {
		return nil, err
	}

	//get list of members :(

	memberQuery := "select channelid, userid from userchannel where channelid=?"
	rows, err := db.Query(memberQuery, id)

	if err != nil {
		return nil, err
	}

	userstore := store.conn.InitUserStore()

	var userSlice []*users.User
	throwaway := -1
	var currentUserID int64
	for rows.Next() {

		if err = rows.Scan(throwaway, currentUserID); err != nil {
			return nil, err
		}

		currentUser, err45 := userstore.GetByID(currentUserID)

		if err45 != nil {
			return nil, err45
		}

		userSlice = append(userSlice, currentUser)
	}
	c.Members = userSlice
	return c, nil
}

//Insert inserts a channel into the store.
func (store *ChannelStore) Insert(c *channels.Channel) (*channels.Channel, error) {
	db := store.conn.db

	if err := c.Validate(); err != nil {
		return nil, err
	}

	insq := "insert ignore into channels (channelname, channeldescription, isprivate, createdat, creatorid) values (?,?,?,?,?)"

	res, err1 := db.Exec(insq, c.Name, c.Description, c.Private, time.Now(), c.Creator.ID)
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

	c.ID = id

	store.AddMember(c.Creator.ID, c.ID)

	return c, nil
}

//AddMember adds a new member to user-channel
func (store *ChannelStore) AddMember(userid int64, channelid int64) error {
	db := store.conn.db
	insq2 := "insert into userchannel (userid, channelid) values (?,?)"
	res2, err1 := db.Exec(insq2, userid, channelid)
	if err1 != nil {
		return err1
	}

	_, err3 := res2.LastInsertId()
	if err3 != nil {
		return err3
	}

	return nil
}

//DeleteMember deletes a channel from the store. As well as all messages associated with channel
func (store *ChannelStore) DeleteMember(userid int64, channelid int64) error {
	db := store.conn.db

	//delete memberlist
	ex3 := "delete from userchannel where userid = ? and channelid = ?"

	_, err := db.Exec(ex3, userid, channelid)

	if err != nil {
		return err
	}

	return nil
}

//Delete deletes a channel from the store. As well as all messages associated with channel
func (store *ChannelStore) Delete(id int64) error {
	db := store.conn.db

	ex := "delete from channels where channelid = ?"

	_, err := db.Exec(ex, id)

	if err != nil {
		return err
	}

	//delete messages
	ex2 := "delete from messages where channelid = ?"

	_, err = db.Exec(ex2, id)

	if err != nil {
		return err
	}

	//delete memberlist
	ex3 := "delete from userchannel where channelid = ?"

	_, err = db.Exec(ex3, id)

	if err != nil {
		return err
	}

	return nil
}

//GetRecipients gets the ids of every user that a message is being sent to
func (store *ChannelStore) GetRecipients(message *messages.Message) ([]int64, error) {
	channel, err := store.GetByID(message.ChannelID)
	if err != nil {
		return nil, err
	}
	return channel.GetMemberIDs()
}
