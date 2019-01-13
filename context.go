package jie

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/proto"
)

// Context ： 上下文
type Context struct {
	Link *Link
	DP   IProtocol
	RT   IRouter
}

// NewContext ： 实例化
func NewContext() *Context {
	return &Context{}
}

// BindProtoReq ： 绑定proto请求
func (c *Context) BindProtoReq(p proto.Message) error {
	return proto.Unmarshal(c.DP.Data(), p)
}

// PackProtoResp : 封装proto响应
func (c *Context) PackProtoResp(p proto.Message) ([]byte, error) {
	return proto.Marshal(p)
}

// Send : 给本连接发送数据
func (c *Context) Send(d []byte) (int, error) {
	return c.Link.Conn.Write(d)
}

// Redirect : 重定向到某一条路由
func (c *Context) Redirect(rs ...interface{}) {
	c.RT.Do(c, rs...)
}

// Broadcast : 广播数据
func (c *Context) Broadcast(d []byte, cs []net.Conn) error {
	var errsStr []string
	for i := 0; i < len(cs); i++ {
		if _, err := cs[i].Write(d); err != nil {
			errStr := fmt.Sprintf("%s:%s", cs[i].RemoteAddr(), err.Error())
			errsStr = append(errsStr, errStr)
		}
	}
	return errors.New(strings.Join(errsStr, ";"))
}

// Get : 获取上下文缓存
func (c *Context) Get(k string) interface{} {
	return c.Link.Cache[k]
}

// Put : 设置上下文缓存
func (c *Context) Put(k string, v interface{}) {
	c.Link.Cache[k] = v
}
