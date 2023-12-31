package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

/**
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数是可以通过zinx.json由用户进行配置
*/

type GlobalObj struct {
	//当前zinx全局的server对象
	TcpServer ziface.IServer
	//当前服务器主机监听的IP
	Host string
	//当前服务器主机监听的端口号
	TcpPort int
	//当前服务器名称
	Name string
	//当前zinx的版本号
	Version string
	//当前服务器允许的最大连接数
	MaxConn int
	//当前zinx框架数据包的最大值
	MaxPackageSize uint32
	//worker工作池的队列的个数
	WorkerPoolSize uint32
	//每个worker对应的消息队列的任务数量的最大值
	MaxWorkerTaskLen uint32
}

/*
GlobalObject 定义一个全局的对外GlobalObj对象
*/
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		return
	}
	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 提供一个init方法，初始化当前的GlobalObject
func init() {
	//如果配置文件没有家在，默认的值
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "ZinxServerApp",
		Version:          "V0.6",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	//应该尝试从conf/zinx.json中去加载一些用户自定义的参数
	GlobalObject.Reload()
}
