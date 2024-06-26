package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// 将数据包编码（即加上包头再转为二进制）
func Encode(mes []byte) ([]byte, error) {
	//获取发送数据的长度，并转换为四个字节的长度，即int32
	len := uint16(len(mes))
	//创建数据包
	dataPackage := new(bytes.Buffer) //使用字节缓冲区，一步步写入性能更高

	//先向缓冲区写入包头
	//大小端口诀：大端：尾端在高位，小端：尾端在低位
	//编码用小端写入，解码也要从小端读取，要保持一致
	err := binary.Write(dataPackage, binary.LittleEndian, len) //往存储空间小端写入数据
	if err != nil {
		return nil, err
	}
	//写入消息
	err = binary.Write(dataPackage, binary.LittleEndian, mes)
	if err != nil {
		return nil, err
	}

	return dataPackage.Bytes(), nil
}

// 解码数据包
func Decode(reader *bufio.Reader) ([]byte, error) {
	//读取数据包的长度（从包头获取）
	lenByte, err := reader.Peek(2) //读取前四个字节的数据
	if err != nil {
		return []byte{}, err
	}
	//转成Buffer对象,设置为从小端读取
	buff := bytes.NewBuffer(lenByte)

	var len uint16                                     //读取的数据大小，初始化为0
	err = binary.Read(buff, binary.LittleEndian, &len) //从小端读取
	if err != nil {
		return []byte{}, err
	}

	//读取消息
	pkg := make([]byte, int(len)+2)
	//Buffered返回缓冲区中现有的可读取的字节数
	if reader.Buffered() < int(len)+2 { //如果读取的包头的数据大小和读取到的不符合
		hr := 0
		for hr < int(len)+2 {
			l, err := reader.Read(pkg[hr:])
			if err != nil {
				return []byte{}, err
			}
			hr += l
		}
	} else {
		_, err := reader.Read(pkg)
		if err != nil {
			return []byte{}, err
		}
	}

	return pkg[2:], nil

}
