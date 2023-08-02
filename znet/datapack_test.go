package znet

import (
	"fmt"
	"net"
	"testing"
)

/**
这段代码实现了一个简单的TCP服务器和客户端,演示了自定义的拆包器来处理粘包问题。

主要的流程是:

1. 服务器端监听端口,启动一个goroutine循环Accept新连接。

2. 对每个新连接,启动一个goroutine做业务处理。

3. 业务处理goroutine中,使用自定义的拆包器(DataPack结构)来拆分数据。

4. 客户端模拟粘包情况,连续不断发送了两个数据包。

5. 服务器端成功地把粘包的数据包拆分开来,打印出了两个消息。

关键点:

1. 为什么需要两个goroutine?

   - 主goroutine监听端口,接收新连接。

   - 每个连接一个goroutine处理业务,避免 blocking,实现并发。

2. 拆包器的作用?

   - TCP是流式协议,没有界限的字节流,需要自定义协议来划分数据包。

   - 拆包器实现了数据包的头部格式,提供了Unpack、Pack方法。

3. 拆包的流程?

   - 读取头部,解析出数据长度。

   - 再根据数据长度读取 Package 数据。

4. 粘包是如何产生和处理的?

   - 发送端连续不断发送小数据包,接收端一次性读取,形成大数据块,就是粘包。

   - 使用拆包器沿包头制定的格式去拆分。

所以使用两个goroutine,既能并发处理连接,又能优雅地解决TCP流的粘包问题。
*/

func TestDataPack(t *testing.T) {

	//创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("listen error") // listen error
		return
	}

	//从客户端读取数据，拆包处理
	go func() {
		for {
			accept, err := listener.Accept()
			if err != nil {
				fmt.Println("accept error") // accept error
				return                      // accept error
			}

			//创建一个go承载，负责从客户端处理业务
			go func(conn net.Conn) {
				//处理客户端的请求
				//----拆包的过程----
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//1.第一次从conn读，把包的head读出来，得到datalen
					headData := make([]byte, dp.GetHeadLen())
					_, err := conn.Read(headData)
					if err != nil {
						fmt.Println("read head error") // read head error
						break                          // read head error
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err) // server unpack err: EOF
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//2.第二次从conn读，把dataLen的内容读出来，得到data
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := conn.Read(msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err) // server unpack data err: EOF
							return
						}
						//封装到msg中
						fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}
				}
			}(accept)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial error") // client dial error
		return
	}

	//创建一个封包对象dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err) // client pack msg1 error: <nil>
		return
	}
	//封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error:", err) // client pack msg2 error: <nil>
		return
	}
	//将sendData1，和 sendData2 拼接一起，组成粘包
	sendData1 = append(sendData1, sendData2...)
	//一次性发送给服务端
	conn.Write(sendData1)

	//客户端阻塞
	select {}

	//==> Recv Msg: ID= 1 , len= 4 , data= zinx
	//==> Recv Msg: ID= 2 , len= 7 , data= nihao!!
}
