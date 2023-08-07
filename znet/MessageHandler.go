package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

/**
消息处理模块的实现
*/

type MessageHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandle 初始化/创建MsgHandle方法
func NewMsgHandle() *MessageHandle {
	return &MessageHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (m *MessageHandle) DoMsgHandler(request ziface.IRequest) {
	//从Request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not found.")
	}
	//根据MsgID调度对应Router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MessageHandle) AddRouter(msgID uint32, router ziface.IRouter) {

	//判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		//id已经注册
		panic("repeat api,msg id = " + strconv.Itoa(int(msgID)))
	}
	// 添加msg与API的绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add msg id = ", msgID, " success!")
}
