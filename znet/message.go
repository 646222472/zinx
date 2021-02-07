package znet

// Message 消息的封装
type Message struct {
	ID      uint32 // 消息的 ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}

// NewMessage 创建一个 Message 消息包
func NewMessage(msgID uint32, data []byte) *Message {
	return &Message{
		ID:      msgID,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// GetMsgID 获取消息的 ID
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

// GetDataLen 获取消息的长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// SetMsgID 设置消息的 ID
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

// SetDataLen 设置消息的长度
func (m *Message) SetDataLen(datalen uint32) {
	m.DataLen = datalen
}

// SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
