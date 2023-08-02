package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

/**
封包、拆包的具体模块
在Go语言中，我们通常通过方法来组织和结构化代码，即使这些方法可能并不需要访问结构体的字段。把这些函数放在DataPack上，会使得你的代码更容易测试、复用和理解。
DataPack结构体在这里的作用相当于一个命名空间，它使得所有跟数据包打包和解包相关的操作被组织到一起。当你看到dp.UnPack()或dp.Pack()时，很明显这些方法是在做数据包的打包和解包操作。
此外，这种模式可以让你在未来更容易地扩展代码。例如，如果你在未来想要增加一些字段到DataPack中来影响打包和解包操作，那么你只需要在已有的方法中添加这些字段，而不需要更改函数签名或者在全局范围内添加状态。
另外要注意的是，DataPack实现了ziface.IDataPack接口，这个接口约定了DataPack需要实现哪些方法。有了这个接口，我们就可以编写其他实现了这个接口的结构体，为数据打包和解包提供不同的实现。这是面向接口的编程思想，它可以让我们的代码更具灵活性和可维护性。
*/

type DataPack struct{}

// NewDataPack 拆包、封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// datalen是uint32类型（4字节），dataID是uint32类型（4字节），所以直接返回8
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})
	//将datalen写进databuf中
	err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}
	//将dataId写进databuf中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}
	//将data数据写进databuf中
	err = binary.Write(dataBuf, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

// UnPack 拆包方法(只需要将包的head信息读出来就可以了，之后再根据head信息里的data的长度进行一次读取)
func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个新的 Reader，并让这个 Reader 的缓冲区内容为 binaryData 这个切片
	dataBuf := bytes.NewReader(binaryData)

	//只解压head信息，得到datalen和msgID
	msg := &Message{}

	//binary.Read() 是Go语言标准库 encoding/binary 中的一个函数，它用于从一个输入源（在这里就是 dataBuf）中读取二进制数据并存储到给定的数据结构（在这里就是 &msg.DataLen）。
	err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}
	//读msgId
	err = binary.Read(dataBuf, binary.LittleEndian, &msg.Id)
	if err != nil {
		return nil, err
	}

	//判断datalen是否已经超出允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("包长度过长")
	}
	return msg, nil
}
