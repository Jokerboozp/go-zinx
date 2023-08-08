package main

import (
	"fmt"
	"net"
	"time"
	"zinx/znet"
)

/*
*
模拟客户端
*/
func main() {
	fmt.Println("client0 start ...")

	time.Sleep(time.Second * 1)
	//直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("dial err:", err)
		return
	}
	//链接调用Write写数据
	for {
		//发送封包的message消息
		dp := znet.NewDataPack()
		msg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx V0.9 client0 Test Message")))
		if err != nil {
			fmt.Println("pack err:", err)
			return
		}
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("write err:", err)
			return
		}
		//服务器回复的消息
		//先读取流中的head部分，得到ID和dataLen
		headData := make([]byte, dp.GetHeadLen())
		_, err = conn.Read(headData)
		if err != nil {
			fmt.Println("read head err:", err)
			return
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msgHead.GetMsgLen() > 0 {
			//msg是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			//根据dataLen从io中读取字节流
			_, err := conn.Read(msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			fmt.Println("==> Recv Server Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		//cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
