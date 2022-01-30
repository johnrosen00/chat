package channels

import "chat/server/gateway/models/messages"

//Store accesses Channels
type Store interface {
	//GetByID channel
	GetByID(id int64) (*Channel, error)
	Insert(*Channel) (*Channel, error)
	Delete(id int64) error
	Edit(id int64, edit *ChannelEdit) (*Channel, error)
	Init() error
	GetAllChannels() ([]*Channel, error)
	AddMember(userid int64, channelid int64) error
	DeleteMember(userid int64, channelid int64) error
	GetRecipients(message *messages.Message) ([]int64, error)
}
