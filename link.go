package jie

import (
	"net"

	"github.com/Justyer/lingo/bytes"
)

type Link struct {
	Conn    net.Conn
	BufPool []byte
	Cache   map[string]interface{}
}

func NewLink() *Link {
	lnk := &Link{}
	lnk.Cache = make(map[string]interface{})
	return lnk
}

func (self *Link) Read(bs int) error {
	buf := make([]byte, bs)
	l, err := self.Conn.Read(buf)
	if err != nil {
		return err
	}

	// 去掉多余的0字节
	buf = buf[:l]

	self.BufPool = bytes.Extend(self.BufPool, buf)

	return nil
}

func (self *Link) BufPop(l int) {
	self.BufPool = self.BufPool[l:]
}
