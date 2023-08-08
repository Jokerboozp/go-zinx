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

// DoConnectionBegin 创建链接之后执行的钩子函数
func DoConnectionBegin(connection ziface.IConnection) {
	fmt.Println("do connection begin")
	err := connection.SendMsg(202, []byte("DoConnectionBegin"))
	if err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置一些属性
	connection.SetProperty("name", "jokerboozp")
	connection.SetProperty("host", "jokerboozp.top")
}

// DoConnectionLost 链接断开前需要执行的函数
func DoConnectionLost(connection ziface.IConnection) {
	fmt.Println("do connection lost call")
	fmt.Println("connID = ", connection.GetConnID(), " is lost")

	//获取链接属性
	if name, err := connection.GetProperty("name"); err == nil {
		fmt.Println("name = ", name)
	}

	if host, err := connection.GetProperty("host"); err == nil {
		fmt.Println("host = ", host)
	}
}

func main() {
	//创建一个Server具柄，基于zinx的api
	s := znet.NewServer("[zinx v1.0]")

	//注册链接的hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//启动Server
	s.Serve()
}
