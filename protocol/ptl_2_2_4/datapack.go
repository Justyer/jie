package protocol

import (
	"errors"

	"github.com/Justyer/jie"

	"github.com/Justyer/lingo/bytes"
)

// Protocol : 自定义协议
// 协议构成 : 2-2-4-x
// 2byte : 大类型
// 2byte : 小类型
// 4byte : 数据字段长度
// xbyte : 数据字段
type Protocol struct {
	OgnByte  []byte
	MsgType  uint16
	MsgCmd   uint16
	DataLen  uint32
	DataByte []byte
}

// NewProtocol : 实例化
func NewProtocol() *Protocol {
	return &Protocol{}
}

// AssertProtocol : 断言接口中的协议是否是该协议
func AssertProtocol(p jie.IProtocol) *Protocol {
	return p.(*Protocol)
}

// New : 实例化
func (p *Protocol) New() jie.IProtocol {
	return &Protocol{}
}

// IsWhole : 判断字节数组是否完整
// 这个方法只能判断这个字节数组是否可以解析，不能保证一定包含这个协议
func (p *Protocol) IsWhole(buf []byte) bool {
	bufLen := len(buf)
	// 如果字节长度大于8，才可能是完整的数据包（不是一定）
	if bufLen < 8 {
		return false
	}
	var lenI uint32
	bytes.ByteToForLE(buf[4:8], &lenI)
	if bufLen < int(lenI)+8 {
		return false
	}
	return true
}

// Parse : 解析数据包
func (p *Protocol) Parse(b []byte) {
	var typeB, cmdB, lenB []byte

	for i := 0; i < 8; i++ {
		switch i {
		case 0, 1:
			typeB = append(typeB, b[i])
		case 2, 3:
			cmdB = append(cmdB, b[i])
		case 4, 5, 6, 7:
			lenB = append(lenB, b[i])
		}
	}

	// 解析大分类和小分类
	bytes.ByteToForLE(typeB, &p.MsgType)
	bytes.ByteToForLE(cmdB, &p.MsgCmd)
	bytes.ByteToForLE(lenB, &p.DataLen)

	p.DataByte = b[8 : p.DataLen+8]
}

// Get : 从字节流中获取数据包
func (p *Protocol) Get(b []byte) (int, error) {
	if p.IsWhole(b) {
		p.Parse(b)
		p.OgnByte = b[:p.DataLen+8]
		return int(p.DataLen) + 8, nil
	}
	return 0, errors.New("buf is not a whole datapack")
}

// Data : 获取协议中数据字段
func (p *Protocol) Data() []byte {
	return p.DataByte
}
