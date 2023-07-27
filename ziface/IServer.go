package ziface

// IServer 定义一个服务器接口
type IServer interface {
	// Start 启动一个服务器
	Start()

	// Stop 停止一个服务器
	Stop()

	// Serve 运行服务器
	Serve()
}
