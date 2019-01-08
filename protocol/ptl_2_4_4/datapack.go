package protocol

import (
	"errors"

	"github.com/Justyer/jie"

	"github.com/Justyer/lingo/bytes"
)

type Protocol struct {
	OgnByte  []byte
	MsgType  uint16
	MsgCmd   uint16
	DataLen  uint32
	DataByte []byte
}

func NewProtocol() *Protocol {
	return &Protocol{}
}

func AssertProtocol(p jie.IProtocol) *Protocol {
	return p.(*Protocol)
}

func (self *Protocol) New() jie.IProtocol {
	return &Protocol{}
}

// 判断字节数组是否完整
// 这个方法只能判断这个字节数组是否可以解析，不能保证一定包含这个协议
func (self *Protocol) IsWhole(buf []byte) bool {
	buf_len := len(buf)
	// 如果字节长度大于8，才可能是完整的数据包（不是一定）
	if buf_len < 8 {
		return false
	}
	var len_i uint32
	bytes.ByteToForLE(buf[4:8], &len_i)
	if buf_len < int(len_i)+8 {
		return false
	}
	return true
}

// 解析数据包
func (self *Protocol) Parse(b []byte) {
	var type_b, cmd_b, len_b []byte

	for i := 0; i < 8; i++ {
		switch i {
		case 0, 1:
			type_b = append(type_b, b[i])
		case 2, 3:
			cmd_b = append(cmd_b, b[i])
		case 4, 5, 6, 7:
			len_b = append(len_b, b[i])
		}
	}

	// 解析大分类和小分类
	bytes.ByteToForLE(type_b, &self.MsgType)
	bytes.ByteToForLE(cmd_b, &self.MsgCmd)
	bytes.ByteToForLE(len_b, &self.DataLen)

	self.DataByte = b[8 : self.DataLen+8]
}

// 从字节流中获取数据包
func (self *Protocol) Get(b []byte) (int, error) {
	if self.IsWhole(b) {
		self.Parse(b)
		self.OgnByte = b[:self.DataLen+8]
		return int(self.DataLen) + 8, nil
	}
	return 0, errors.New("buf is not a whole datapack")
}

func (self *Protocol) Data() []byte {
	return self.DataByte
}
