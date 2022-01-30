package dbwrapper

type ConnType int

const (
	UserStore ConnType = iota
	MessageStore
	ChannelStore
)
