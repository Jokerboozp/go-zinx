package znet

type Message struct {
	Id      uint32 //消息ID
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

// GetMsgId 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// GetData 获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

// GetMsgLen 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// SetMsgId 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// SetData 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// SetMsgLen 设置消息的长度
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}
