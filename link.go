package jie

import (
	"net"

	"github.com/Justyer/lingo/bytes"
)

// Link : TCP连接信息
type Link struct {
	// TCP连接
	Conn net.Conn
	// 字节缓冲区
	BufPool []byte
	// 一次连接过程的缓存
	Cache map[string]interface{}
}

// NewLink ： 实例化连接
func NewLink() *Link {
	lnk := &Link{}
	lnk.Cache = make(map[string]interface{})
	return lnk
}

// Read : 从TCP缓冲区中读取数据字节
func (lnk *Link) Read(bs int) error {
	buf := make([]byte, bs)
	l, err := lnk.Conn.Read(buf)
	if err != nil {
		return err
	}

	// 去掉多余的0字节
	buf = buf[:l]

	lnk.BufPool = bytes.Extend(lnk.BufPool, buf)

	return nil
}

// BufPop : 将匹配好的数据包字节从连接缓冲区打出
func (lnk *Link) BufPop(l int) {
	lnk.BufPool = lnk.BufPool[l:]
}
