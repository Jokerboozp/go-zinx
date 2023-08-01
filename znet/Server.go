package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
	//当前的Server添加一个Router，server注册的链接对应的处理业务
	Router ziface.IRouter
}

// AddRouter 路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("add router success")
}

// NewServer 初始化Server模块
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}

// CallBackToClient 定义当前客户端链接所绑定的handle api(目前这个handle是写死的，以后优化应该由用户自定义handle方法)
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	fmt.Println("connection handle callback")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write err,", err)
//		return err
//	}
//	return nil
//}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Println("[Start] Server Listener at IP", s.IP, "Port ", s.Port, "is starting")

	go func() {
		//获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tco addr error:", err)
			return
		}
		//监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err ", err)
			return
		}
		fmt.Println("start zinx server success, ", s.Name, " success,Listening")

		var connID uint32
		connID = 0

		//阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			//如果客户端链接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			//将处理新链接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, connID, s.Router)
			connID++

			//启动当前的业务链接处理
			go dealConn.Start()

			//已经与客户端建立链接，做一个最基本的最大512字节长度的回显业务
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err := conn.Read(buf)
			//		if err != nil {
			//			fmt.Println("receive buf err:", err)
			//			continue
			//		}
			//		fmt.Printf("receive client buf %s,cnt %d\n", buf, cnt)
			//
			//		//回显功能
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back buf error:", err)
			//			continue
			//		}
			//	}
			//}()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {

}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	/**
	阻塞状态
	select{}会阻塞当前的Goroutine,但不会阻止程序继续运行。原因是:

	main函数在启动Server时,会用go关键字启动一个新的Goroutine去执行Serve方法
	这个Goroutine被select{}阻塞了,但是main函数本身不受影响
	main函数可以继续往下执行,做其他初始化工作
	当初始化工作完成,main函数就可以退出了
	但这个时候Serve方法还在另一个Goroutine中运行并阻塞着
	所以程序整体并不会退出,会持续运行
	而Serve中的其他代码在select{}之前,包括启动监听端口,创建处理连接的Goroutine等都可以继续工作
	所以总结一下:

	select{}阻塞了Serve所在的Goroutine,但不影响其他Goroutine
	主程序main函数还可以继续执行其他逻辑
	服务器已启动的部分仍可以处理连接请求
	因此服务器可以持续运行,程序不会退出
	*/
	select {}
}
