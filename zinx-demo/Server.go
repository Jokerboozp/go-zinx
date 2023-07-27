package main

import "zinx/znet"

/**
基于Zinx框架来开发的服务器端应用程序
*/

func main() {
	//创建一个Server具柄，基于zinx的api
	s := znet.NewServer("[zinx b0.2]")
	//启动Server
	s.Serve()
}
