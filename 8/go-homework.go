package main

import (
	"encoding/binary"
	"errors"
)

/// 1. 总结几种 socket 粘包的解包方式: fix length/delimiter based/length field based frame decoder。尝试举例其应用
//    1. fix length
//       1. 客户端和服务端约定每次信息交互的长度，客户端每次发送固定长度的数据，服务端每次接受固定长度的内容
//       2. 这种策略，会造成资源的浪费，以及消息过长（消息长度大于固定长度）时使用固定长度会导致消息被分割
//    2. delimiter based
//       1. 使用特殊字符当作消息结束信号，如`\n`
//       2. 这种策略要求消息正文不允许出现该特殊字符，在某些场景下不太适合，如要保留消息的原始（原始消息可能是个文本消息，本身含有换行符）格式的时候
//    3. length field based frame decoder
//       1. 约定消息体的前多少byte为消息长度字段。
//       2. 解决了使用特殊换行符的场景限制，比`fix length`方案使用更少的资源，也问题解决的更好

/// 2. 实现一个从 socket connection 中解码出 goim 协议的解码器。
// PacketLen	HeaderLen	Version	Operation	Sequence	Body
// 4bytes	    2bytes	    2bytes	4bytes	    4bytes	    PacketLen - HeaderLen

type Goim struct {
	PacketLen uint32
	HeaderLen uint16
	Version   uint16
	Operation uint32
	Sequence  uint32
	Body      []byte
}

var (
	ErrNotGoimData = errors.New("goim: input data wrong")
)

func NewGoim(data []byte) (*Goim, error) {
	if len(data) < 16 {
		return nil, ErrNotGoimData
	}

	im := &Goim{
		PacketLen: binary.BigEndian.Uint32(data[:4]),
		HeaderLen: binary.BigEndian.Uint16(data[4:6]),
		Version:   binary.BigEndian.Uint16(data[6:8]),
		Operation: binary.BigEndian.Uint32(data[8:12]),
		Sequence:  binary.BigEndian.Uint32(data[12:16]),
	}
	im.Body = data[im.HeaderLen:im.PacketLen]
	return im, nil
}

func main() {
	testData := func() []byte {
		data := "Hello, World."
		dataSize := len(data)
		allSize := 16+dataSize
		d := make([]byte, allSize+4)
		binary.BigEndian.PutUint32(d[:4], uint32(allSize))
		binary.BigEndian.PutUint16(d[4:6], uint16(16))
		binary.BigEndian.PutUint16(d[6:8], uint16(1))
		binary.BigEndian.PutUint32(d[8:12], uint32(0))
		binary.BigEndian.PutUint32(d[12:16], uint32(1))
		data += "12"
		copy(d[16:], []byte(data))
		return d
	}()

	goim, err := NewGoim(testData)
	if err != nil {
		panic("something wrong")
	}
	println(goim.PacketLen, string(goim.Body))
}
