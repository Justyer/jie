package jie

import (
	"github.com/golang/protobuf/proto"
)

type Context struct {
	Link *Link
	DP   IProtocol
}

func NewContext() *Context {
	return &Context{}
}

func (self *Context) BindProtoReq(p proto.Message) error {
	return proto.Unmarshal(self.DP.Data(), p)
}

func (self *Context) PackProtoResp(p proto.Message) ([]byte, error) {
	return proto.Marshal(p)
}

func (self *Context) Send(d []byte) (int, error) {
	return self.Link.Conn.Write(d)
}
