### v0.1
[![pCvRwTO.png](https://s1.ax1x.com/2023/07/27/pCvRwTO.png)](https://imgse.com/i/pCvRwTO)
#### cnt, err := conn.Read(buf)
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

#### select {}
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

#### func NewServer(name string) ziface.IServer

```text
Server结构体实现了IServer接口所要求的所有方法,所以它满足IServer接口的要求。

在NewServer函数中,返回一个Server指针,该指针满足IServer接口,所以可以赋值给IServer类型。

这是Go语言接口的一个常见用法,通过接口定义一组方法,然后自定义结构体实现这些方法,这样该结构体就实现了该接口,可以赋值给接口变量。

所以NewServer返回一个Server指针,可以赋值给IServer类型,这是因为Server实现了IServer接口。
```

### v0.2
[![pCvh5tO.png](https://s1.ax1x.com/2023/07/27/pCvh5tO.png)](https://imgse.com/i/pCvh5tO)

#### Connection.go中StartReader方法

```text
为什么这段代码里面第一个if用continue，第二个if用break

原因是因为Read方法在读取客户端数据时,可能会出现不同类型的错误。

第一个if处理的是读取过程中的可恢复错误,比如临时的网络问题导致读不到数据。这种情况下我们打印错误,然后通过continue跳过本次循环,等待下一次循环再尝试读取。

第二个if处理的是不可恢复的错误,比如客户端断开了连接等致命错误。这种情况下我们直接break跳出循环,结束这个goroutine。

总结一下:

- continue用于跳过当前循环,进入下一次循环。适用于可恢复的错误。

- break用于直接跳出for循环,结束goroutine。适用于不可恢复的致命错误。

使用continue和break的不同取决于错误的严重性和可恢复性。continue让循环继续,break终止循环。根据实际情况选择合适的控制结构。

在网络编程中,一些比较轻微的临时错误例如:

- 由于网络闪断,短时间读取数据失败
- 对端暂停发送数据,本次读取返回0字节
- 数据包暂时阻塞,读超时但并非严重问题
- 资源短暂不可用,例如内存不足导致读失败

这些错误都是暂时性的,中断当前的读取操作,但并不会真正断开连接或造成致命问题。

对于这类轻微错误,我们可以打印出来,然后continue跳过当前循环,等待一段时间后再次尝试,有可能就能恢复正常读取了。

而严重错误通常是:

- 对端关闭了连接
- 本端关闭或超时了连接 
- TCP连接断开了
- 出现了严重内部错误,例如解析错误导致后续处理无法继续

这类错误通常意味着连接已经完全断开了,或者程序进入了一个无法恢复的错误状态。对于这类错误,continue是无法修复的,必须立即break跳出循环,防止程序进入异常状态。

总结一下,轻微错误可以retry,严重错误则需要立即处理。判断一个错误是否严重需要根据具体场景来决定。
```

#### 为什么要在NewConnection方法里设置HandleFunc回调函数

```text
设置回调函数的目的主要有以下几点:

NewConnection函数是创建连接的函数,不应该处理具体的数据。但是后续需要一个函数来处理数据,所以在创建时就设置回调函数,方便后续使用。
Go语言倡导面向接口编程,这里的回调函数实际上是实现一个接口ziface.Handler的Component方法。通过回调注入不同的实现,使连接的行为可配置和扩展。
回调函数能够解耦连接创建和业务数据处理逻辑。连接本身不需要知道具体的处理逻辑,由调用者注入不同的回调实现即可。
设置回调函数实现了依赖注入模式,Connection依赖一个回调函数来处理数据,通过接口注入具体实现,降低了耦合。
最后一点,由于连接是长连接,需要持续不断地处理数据,所以需要在创建时就设置好回调函数,后续直接使用,避免每次读到数据时都要重新指定回调函数。
总结一下,设置回调函数是一种编程模式,可以实现解耦,依赖注入,提高代码扩展性和可维护性。不将具体实现写在连接内部,而是通过回调注入,可以提高连接对象的通用性和复用性。
```

#### 创建时设置的回调函数在后续都可以直接使用

```text
在NewConnection函数中设置的回调函数handleAPI,在后续的连接生命周期里都可以直接使用,不会被删除或改变。

具体来说:

在NewConnection中,会将传入的handleAPI回调函数保存在Connection对象的handleAPI字段中。
在后续的StartReader goroutine 中,会通过connection.handleAPI(conn, buf, cnt)来调用这个回调函数。
这个回调函数在连接的整个生命周期内都可以被重复调用,直到连接关闭。
调用者传入什么回调函数,在后续这个连接上就是什么数据处理逻辑。
回调函数本身不会改变或者被删除,一直持有连接对象的引用,除非连接关闭。
这样可以确保整个连接生命周期内,都使用初始化时注入的那个回调函数来处理数据,避免反复绑定的开销。
综上所述,在连接初始化时注入的回调函数会贯穿整个连接生命周期,可以被持续重复使用,而不会被删除或改变,除非连接关闭。这充分利用了Go语言函数都是第一类对象的特点。
```

#### NewConnection中的handleAPI和StartReader中的c.HandleAPI是一个回调函数

```text
在StartReader方法中,调用c.handleAPI时就是使用的在NewConnection中传入并保存的那个回调函数。

整个流程是:

在NewConnection中,调用者传入回调函数,保存到Connection对象的handleAPI字段。
在后续的StartReader中,从Connection对象中取出这个handleAPI字段,调用它。
这个handleAPI字段是一开始保存的那个回调函数,不会变化。
所以StartReader中的c.handleAPI,调用的就是一开始注入的回调函数。
通过这种代码设计,连接对象本身不需要硬编码任何数据处理逻辑,全部通过回调函数实现。

连接对象负责数据读取,回调函数负责业务处理。这样可以很好地实现关注点分离和解耦。

可以总结为:

NewConnection中注入回调函数 → 连接对象保存回调函数 → StartReader中调用回调函数
```

#### 回调函数handleAPI并不是直接调用,而是需要传入参数并得到返回值。

```text
回调函数的类型是:
type HandleFunc func(connection Connection, data []byte, cnt int) error

它需要传入connection, data, cnt三个参数。

所以在调用回调函数时,不能简单地写成:

c.handleAPI()

而需要传递PROPER参数:

err = c.handleAPI(c.Conn, buf, cnt)

只有传入匹配的参数,回调函数才能得到需要的信息,进行处理并返回结果。

这与其他语言如JS中的回调函数直接调用有区别。Go需要参数匹配调用。

所以即使有了回调函数,也必须传参正确调用,不能简单地直接调用。
```

### v0.3
![b799e8d1874559a7b0bb8f0e195cc6fd.png](https://i.mji.rip/2023/07/31/b799e8d1874559a7b0bb8f0e195cc6fd.png)
