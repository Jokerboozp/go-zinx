- cnt, err := conn.Read(buf)
```text

conn.Read(buf)这个语句的作用是从conn网络连接中读取数据并放入buf缓冲区。它返回两个值:

- cnt: 读取到的字节数
- err: 读取是否发生错误

举个例子:

```go
buf := make([]byte, 1024)

cnt, err := conn.Read(buf)

if err != nil {
  // 处理错误
} else {
  // 在buf[:cnt]中获取读取的数据 
}

具体分析:

- 先创建一个1024字节的缓冲区buf
- 调用conn.Read(),尝试从连接读取数据到buf中
- 它会返回实际读取了多少字节cnt,存放在buf[:cnt]
- 还会返回是否出现错误err

所以通过检查cnt,就知道读取了多少数据。如果err不是nil,则表示出错了。

这种读取方式可以处理TCP数据Confirm的读取,也就是讲缓冲区传递给conn.Read,让它自行读取适当大小的数据,非常方便。

conn.Read()会被阻塞,直到有数据可读或者出现错误。所以它是异步读取网络数据的常用方法。
```

- select {}
```text
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
```

- func NewServer(name string) ziface.IServer

```text
Server结构体实现了IServer接口所要求的所有方法,所以它满足IServer接口的要求。

在NewServer函数中,返回一个Server指针,该指针满足IServer接口,所以可以赋值给IServer类型。

这是Go语言接口的一个常见用法,通过接口定义一组方法,然后自定义结构体实现这些方法,这样该结构体就实现了该接口,可以赋值给接口变量。

所以NewServer返回一个Server指针,可以赋值给IServer类型,这是因为Server实现了IServer接口。
```