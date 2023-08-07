package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/**
基于Zinx框架来开发的服务器端应用程序
*/

// PingRouter ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

// PreHandle Test PreHandle
//func (p *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("call router PreHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
//	if err != nil {
//		fmt.Println("call back before ping err,", err)
//	}
//}

// Handle Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call router Handle")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func (p *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("call hello router Handle")
	// 先读取客户端的数据，再回写hello...hello...hello
	fmt.Println("recv from client: msgID=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println(err)
	}
}

//// PostHandle Test PostHandle
//func (p *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("call router PostHandle")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
//	if err != nil {
//		fmt.Println("call back after ping err,", err)
//	}
//}

func main() {
	//创建一个Server具柄，基于zinx的api
	s := znet.NewServer("[zinx v0.8]")

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//启动Server
	s.Serve()
}
