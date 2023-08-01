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
