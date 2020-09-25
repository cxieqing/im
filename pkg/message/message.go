package message

type MessageType uint8

const (
	UserMessage MessageType = 1

	GroupMessage MessageType = 2
)

type ReadType uint8

const (
	HasRead ReadType = 1

	UnRead ReadType = 0
)

type ContenType uint8

const (
	ImageContent ContenType = iota
	TextContent
	VdioContent
)

type Message struct {
	Id         int
	Type       MessageType
	Content    []byte
	From       int
	To         int
	IsRead     ReadType
	CreateTime int64
}

func (m Message) save() bool {
	return true
}

func GetUserUnReadMsgById(uid int) {

}
