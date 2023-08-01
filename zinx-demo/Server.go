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

// PreHandle Test PreHandle
func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("call router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping err,", err)
	}
}

// Handle Test Handle
func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("call router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping ping...\n"))
	if err != nil {
		fmt.Println("call back ping err,", err)
	}
}

// PostHandle Test PostHandle
func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("call router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping err,", err)
	}
}

func main() {
	//创建一个Server具柄，基于zinx的api
	s := znet.NewServer("[zinx v0.4]")

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})

	//启动Server
	s.Serve()
}
