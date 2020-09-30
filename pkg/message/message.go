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

type SendType uint8

const (
	HasSend SendType = 1

	UnSend SendType = 0
)

type ContenType uint8

const (
	ImageContent ContenType = iota
	TextContent
	VdioContent
)

type Message struct {
	Id          int
	ContentType ContenType
	Content     string
	From        int
	To          int
	IsRead      ReadType
	IsSend      SendType
	Len         float32
	CreateTime  int64
}

func (m Message) save(mtype MessageType) bool {
	return true
}

func GetUserUnsendMsgByUserId(uid int) []Message {
	return []Message{}
}

func GetUserUnsendMsgByGroupId(gid int) []Message {
	return []Message{}
}

func SaveUnsendGroupMsg(m Message) bool {
	return true
}

func SaveUnsendUserMsg(m Message) bool {
	return true
}

func SaveSendUserMsg(m Message) bool {
	return true
}
