package channels

import (
	"chat/server/models/messages"
	"chat/server/models/users"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

//SQLChannelStore is a store of channels
type SQLChannelStore struct {
	DB *sql.DB
}

//GetAllChannels gets all channels
func (store *SQLChannelStore) GetAllChannels() ([]*Channel, error) {
	q := "select channelid from channels"
	rows, err := store.DB.Query(q)
	if err != nil {
		return nil, err
	}
	var channelSlice []*Channel
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
func (store *SQLChannelStore) Init() error {
	q := "insert ignore into channels (channelname, channeldescription, isprivate, createdat) values (general, general, false, ?)"
	res, err1 := store.DB.Exec(q, time.Now())
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
func (store *SQLChannelStore) Edit(id int64, edit *ChannelEdit) (*Channel, error) {
	if len(edit.Name) != 0 {
		return nil, errors.New("enter a name over 0 char long")
	}
	queryF := "update channels set name = ? , description = ? , editedAt = ? where channelid = ?"

	if _, err := store.DB.Exec(queryF, edit.Name, edit.Description, time.Now(), id); err != nil {
		return nil, err
	}
	return store.GetByID(id)
}

//GetByID returns a channel struct that can be encoded to JSON
func (store *SQLChannelStore) GetByID(id int64) (*Channel, error) {
	q := "select channelid, channelname, channeldescription, isprivate, createdat, creatorid, editedat, from channels where channelid=?"
	row := store.DB.QueryRow(q, id)

	c := &Channel{}
	err := row.Scan(c.ID, c.Name, c.Description, c.Private, c.CreatedAt, c.Creator.ID, c.EditedAt)

	if err != nil {
		return nil, err
	}

	//get list of members :(

	memberQuery := "select channelid, userid from userchannel where channelid=?"
	rows, err := store.DB.Query(memberQuery, id)

	if err != nil {
		return nil, err
	}

	userstore := &users.MySQLStore{}
	userstore.DB = store.DB

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
func (store *SQLChannelStore) Insert(c *Channel) (*Channel, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	insq := "insert ignore into channels (channelname, channeldescription, isprivate, createdat, creatorid) values (?,?,?,?,?)"

	res, err1 := store.DB.Exec(insq, c.Name, c.Description, c.Private, time.Now(), c.Creator.ID)
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
func (store *SQLChannelStore) AddMember(userid int64, channelid int64) error {
	insq2 := "insert into userchannel (userid, channelid) values (?,?)"
	res2, err1 := store.DB.Exec(insq2, userid, channelid)
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
func (store *SQLChannelStore) DeleteMember(userid int64, channelid int64) error {
	//delete memberlist
	ex3 := "delete from userchannel where userid = ? and channelid = ?"

	_, err := store.DB.Exec(ex3, userid, channelid)

	if err != nil {
		return err
	}

	return nil
}

//Delete deletes a channel from the store. As well as all messages associated with channel
func (store *SQLChannelStore) Delete(id int64) error {
	ex := "delete from channels where channelid = ?"

	_, err := store.DB.Exec(ex, id)

	if err != nil {
		return err
	}

	//delete messages
	ex2 := "delete from messages where channelid = ?"

	_, err = store.DB.Exec(ex2, id)

	if err != nil {
		return err
	}

	//delete memberlist
	ex3 := "delete from userchannel where channelid = ?"

	_, err = store.DB.Exec(ex3, id)

	if err != nil {
		return err
	}

	return nil
}

//GetRecipients gets the ids of every user that a message is being sent to
func (store *SQLChannelStore) GetRecipients(message *messages.Message) ([]int64, error) {
	channel, err := store.GetByID(message.ChannelID)
	if err != nil {
		return nil, err
	}
	return channel.GetMemberIDs()
}
