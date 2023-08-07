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

### v0.4
![b9e1a9bf53667eaa63fc948ffb092de5.png](https://i.mji.rip/2023/08/01/b9e1a9bf53667eaa63fc948ffb092de5.png)

### v0.5

#### TCP粘包问题

```text
TCP 粘包是在 TCP 传输过程中出现的一个现象。因为 TCP 是面向流的协议，数据像流水一样经过 TCP 连接，所以在应用程序中，无法明确区分每次 write 写入的数据在 TCP 到达对端后是否会被当成独立的消息进行处理。
上述这种情况会导致以下两种现象：
多个发送操作的数据会拼接在一起： 如果你连续发送了两个数据包，可能在对端接收时，被当成一个数据包接收。
一个发送操作的数据被拆分： 如果你发送的数据包过大，超过了 TCP 的滑动窗口大小或者 MTU（最大传输单元），那么这个包会被拆分成多个小包在网络中传输，而在接收端可能会被分多次接收。
这就是所谓的 "TCP 粘包" 问题。为了解决这个问题，通常我们需要在应用层增加帧边界信息，最常见的就是增加一个包头，包含了包的长度或者结束标记，这样在接收端就可以正确地拼装或者区分各个数据包了。
```

#### TLV序列化

![8bdc26e68f7d8510c2855b8919df1d20.png](https://i.mji.rip/2023/08/01/8bdc26e68f7d8510c2855b8919df1d20.png)

```text
TLV（Type-Length-Value）是一种数据序列化的格式，广泛用于各种网络协议和文件格式。TLV 不仅可以编码一对一的键值对，还可表示复杂的、层次化的数据结构。
如其名称所示，一个 TLV 数据单位包含三部分：
Type (类型): 表示数据的类型，例如整数、浮点数、字符串等。它通常用于解析程序了解如何解释接下来的值。
Length (长度): 表示后面的值占用了多少字节。
Value (值): 实际的数据

例如，我们有一个字符串 "Hello, world!"，其 TLV 表现形式可能为：
T : 0x01 (假设我们规定0x01标识字符串)
L : 0x0D (字符串长度为13字节)
V : "Hello, world!"

TLV 格式具有很好的灵活性与扩展性，因为类型字段可以容易地增加新的类型，而长度字段可以使得值的大小变为动态，因此 TLV 在很多协议中都能看到它的身影
```

#### dataBuf := bytes.NewReader(binaryData)

![0e159ad6d3347a7cb1760542c5b500e9.png](https://i.mji.rip/2023/08/01/0e159ad6d3347a7cb1760542c5b500e9.png)

```text
bytes.NewReader(binaryData) 是创建一个新的 Reader，并让这个 Reader 的缓冲区内容为 binaryData 这个切片。
bytes.Reader 是一个可以从中读取数据的对象，实现了 io.Reader, io.ReaderAt, io.Seeker, io.WriterTo, io.ByteScanner,和 io.RuneScanner 接口。创建这个 reader 之后，你就可以用各种方法从中读取数据，比如 Read(), ReadAt(), ReadByte() 等等。
```

#### binary.LittleEndian和binary.BigEndian

![0e159ad6d3347a7cb1760542c5b500e9.png](https://i.mji.rip/2023/08/01/0e159ad6d3347a7cb1760542c5b500e9.png)

```text
大端字节序（Big-Endian）和小端字节序（Little-Endian）是计算机科学中关于存储或传输多字节数据的两种方法。这两种方式的主要区别在于以哪种顺序存储或传输最低位和最高位。
大端字节序（Big-Endian）： 最低有效字节（Least Significant Byte，简称 LSB）存储在内存的最高地址上，最高有效字节（Most Significant Byte，简称 MSB）存储在内存的最低地址上。也就是说，高位字节在前，低位字节在后。这种存储方式类似于人类阅读数字和文字的习惯，第一个字节是最高位字节，因此称为 "大端"。
小端字节序（Little-Endian）： 最低有效字节存储在内存的最低地址上，最高有效字节存储在内存的最高地址上。也就是说，低位字节在前，高位字节在后。因此称为 "小端"。
这两种字节序常在网络编程和硬件编程中遇到，由于不同的机器和网络协议可能采用不同的字节序，因此在处理多字节数值时，需要注意转换问题。
例如，对于16位十六进制数0x1234：
在大端字节序中存储为：12 34
在小端字节序中存储为：34 12
所以，在处理多字节数据时，需要明确数据的字节序并进行正确的转换，以保证数据的正确性。
```

#### append(sendData1, sendData2...)

```text
在 Go 语言中，append 函数用于向一个 slice 添加元素。其第一个参数是你想要添加元素的 slice，之后的参数是你想要添加的元素。
sendData1 = append(sendData1, sendData2...) 这行代码的含义是，将 sendData2 中的所有元素添加到 sendData1 的末尾，并将结果赋值给 sendData1。
... 是 Go 语言的一种语法糖，被称作"可变参数"或"展开操作符"。当它用在 slice 前面时，它会将 slice 展开为一个元素列表。因此，append(sendData1, sendData2...) 就等同于将 sendData2 中的所有元素一个一个添加到 sendData1 中
```

#### 为什么要在空结构体上写方法

```text
在Go语言中，我们通常通过方法来组织和结构化代码，即使这些方法可能并不需要访问结构体的字段。把这些函数放在DataPack上，会使得你的代码更容易测试、复用和理解。
DataPack结构体在这里的作用相当于一个命名空间，它使得所有跟数据包打包和解包相关的操作被组织到一起。当你看到dp.UnPack()或dp.Pack()时，很明显这些方法是在做数据包的打包和解包操作。
此外，这种模式可以让你在未来更容易地扩展代码。例如，如果你在未来想要增加一些字段到DataPack中来影响打包和解包操作，那么你只需要在已有的方法中添加这些字段，而不需要更改函数签名或者在全局范围内添加状态。
另外要注意的是，DataPack实现了ziface.IDataPack接口，这个接口约定了DataPack需要实现哪些方法。有了这个接口，我们就可以编写其他实现了这个接口的结构体，为数据打包和解包提供不同的实现。这是面向接口的编程思想，它可以让我们的代码更具灵活性和可维护性。
```

### v0.6

![6ec73f2eec6bbe31da69aa506577c08a.png](https://i.miji.bid/2023/08/07/6ec73f2eec6bbe31da69aa506577c08a.png)

#### PingRouter和HelloZinxRouter为什么都可以使用Handle方法

````text
ingRouter和HelloZinxRouter都含有一个匿名的znet.BaseRouter字段。这是Go语言的结构体嵌入和方法提升的机制。
在Go中，当访问一个嵌入字段的方法时，如果该结构体没有定义此方法，会自动提升嵌入字段的同名方法。这就是你可以对PingRouter和HelloZinxRouter直接调用Handle方法的原因，即便你并没有在这两个结构体中显式的定义这个方法。
然而在你的代码示例中，你覆盖了（或者说重新实现了）PingRouter和HelloZinxRouter的Handle方法。这使得当你在这两个类型上调用Handle方法时，实际上调用的是你自定义的方法，而不是他们嵌入的znet.BaseRouter字段的方法。
这种设计也是一种常见的面向对象的设计模式：模板方法模式。基类定义了一套操作的框架，具体的步骤则由子类来实现。
````

#### 为什么修改了HelloZinxRouter的Handle方法名称就接收不到返回消息

```text
这是因为在你的框架设计中，Golang的接口方法被用作了消息的处理方法。
在你的代码中，“PingRouter” 和 “HelloZinxRouter” 都是实现了 "znet.BaseRouter" 中定义的接口方法 "Handle"。根据你的架构设计，Handle 方法被当作了处理客户端发送消息的主要逻辑。
当“HelloZinxRouter”的“Handle”方法被重命名，这个结构不再完全实现 "znet.BaseRouter"，使得不会调用到你的特定逻辑，而可能只会运行基础的或者是默认的处理逻辑。
这就是为什么你发现把 “HelloZinxRouter”的“Handle”方法改名后就收不到返回消息了，因为你的特定逻辑不再被调用。
函数名非常重要，如果你实现了一个接口，函数的名称、接收器、参数列表和返回参数都必须和接口定义的完全一致，才算真正实现了该接口，才能被框架正确的识别和调用。
```

### v0.7

![8988ff7040852d9b14af0c5180ba135c.png](https://i.miji.bid/2023/08/07/8988ff7040852d9b14af0c5180ba135c.png)

### v0.8

![a7f97e2f28641800f6b0ac92cf06c044.png](https://i.miji.bid/2023/08/07/a7f97e2f28641800f6b0ac92cf06c044.png)
