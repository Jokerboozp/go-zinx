package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

/**
消息处理模块的实现
*/

type MessageHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作worker池的worker数量
	WorkerPoolSize uint32
}

// NewMsgHandle 初始化/创建MsgHandle方法
func NewMsgHandle() *MessageHandle {
	return &MessageHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
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

// StartWorkerPool 启动一个Worker工作池(开启工作池的动作只能发生一次，因为一个zinx框架只能有一个worker工作池)
func (m *MessageHandle) StartWorkerPool() {
	//根据workerPoolSize分别开启worker，每个worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		//一个worker被启动
		//当前的worker对应的channel消息队列，开辟空间。第0个worker就用第0个channel
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前的worker，阻塞等待消息从channel传递过来
		go m.startOneWorker(i, m.TaskQueue[i])
	}
}

// startOneWorker 启动一个Worker工作流程
func (m *MessageHandle) startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker id = ", workerID, " is started")
	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的request,执行当前request绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker进行处理
func (m *MessageHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//将消息平均分配给不通过的worker
	//根据客户端建立的connID来进行分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("add connID = ", request.GetConnection().GetConnID(), " request msgID = ", request.GetMsgID(), " to workerID = ", workerID)
	//将消息发送给对应的worker的taskQueue即可
	m.TaskQueue[workerID] <- request
}
