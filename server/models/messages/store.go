package messages

//Store retrieves information about messages.
type Store interface {
	//GetByID gets messages by ID
	GetByID(id int64) (*Message, error)
	Insert(message *Message) (*Message, error)
	Delete(id int64) error
	Edit(id int64, edit *MessageEdit) (*Message, error)
	GetRecentX(x int64, channelid int64) (Messages, error)
	GetRecentXBeforeY(x int64, y int64, channelid int64) (Messages, error)
}
