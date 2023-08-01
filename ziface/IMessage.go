package ziface

/*
*
将请求的消息封装到一个Message中，定义抽象的接口
*/

type IMessage interface {
	//GetMsgId 获取消息ID
	GetMsgId() uint32
	//GetData 获取消息的内容
	GetData() []byte
	//GetMsgLen 获取消息的长度
	GetMsgLen() uint32
	//SetMsgId 设置消息的ID
	SetMsgId(id uint32)
	//SetData 设置消息的内容
	SetData(data []byte)
	//SetMsgLen 设置消息的长度
	SetMsgLen(len uint32)
}
