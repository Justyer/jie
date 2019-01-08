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

func (self *Context) PackProtoResp() []byte {

}
