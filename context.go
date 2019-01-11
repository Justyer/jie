package jie

import (
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/proto"
)

type Context struct {
	Link *Link
	DP   IProtocol
	RT   IRouter
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

// 给本连接发送数据
func (self *Context) Send(d []byte) (int, error) {
	return self.Link.Conn.Write(d)
}

// 重定向到某一条路由
func (self *Context) Redirect(rs ...interface{}) {
	self.RT.Do(self, rs...)
}

// 广播数据
func (self *Context) Broadcast(d []byte, cs []net.Conn) string {
	var errs_str []string
	for i, _ := range cs {
		if _, err := cs[i].Write(d); err != nil {
			err_str := fmt.Sprintf("%s:%s", cs[i].RemoteAddr(), err.Error())
			errs_str = append(errs_str, err_str)
		}
	}
	return strings.Join(errs_str, ";")
}

func (self *Context) Get(k string) interface{} {
	return self.Link.Cache[k]
}

func (self *Context) Put(k string, v interface{}) {
	self.Link.Cache[k] = v
}
