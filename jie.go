package jie

import (
	"fmt"
	"io"
	"net"

	"github.com/Justyer/JekoServer/plugin/log"
	"github.com/Justyer/lingo/ip"
)

type ProtocolFunc func() Protocol

type RouterFunc func(*Context, Protocol)

type Engine struct {
	// 路由管理
	RouterGroup

	// 所有连接管理
	conns []net.Conn

	// 自定义协议创建
	protocolFunc ProtocolFunc
	routerFunc   RouterFunc
}

func New() *Engine {
	e := &Engine{}
	return e
}

// 设置自定义协议
func (self *Engine) SetProtocol(f ProtocolFunc) {
	self.protocolFunc = f
}

// 设置路由规则
func (self *Engine) SetRouter(r RouterFunc) {
	self.routerFunc = r
}

// 实例化协议
func (self *Engine) NewProtocol() Protocol {
	return self.protocolFunc()
}

// 本地监听
func (self *Engine) ListenAndLocalServe(port string) {
	self.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port))
}

// 内网监听
func (self *Engine) ListenAndInnerServe(port string) {
	ip := ip.MustInnerIP()
	self.ListenAndServe(fmt.Sprintf("%s:%s", ip, port))
}

// 自定义地址监听
func (self *Engine) ListenAndServe(addr string) {
	listener, err := net.Listen("tcp", addr)
	defer listener.Close()
	if err != nil {
		return
	}
	log.Info("[listener]: (%s)", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Err("[Error accepting]: %s", err.Error())
			continue
		}

		log.Info("[link-start-iport]: (%s)", conn.RemoteAddr().String())
		self.conns = append(self.conns, conn)
		go self.dealData(conn)
	}
}

func (self *Engine) dealData(conn net.Conn) {
	defer conn.Close()

	lnk := NewLink()
	lnk.Conn = conn

	for {
		if err := lnk.Read(50); err != nil {
			if err == io.EOF {
				log.Err("[link-close-iport]: (%s)", lnk.Conn.RemoteAddr().String())
				return
			}
			log.Err("[Error reading]: %s, on (%s)", err.Error(), lnk.Conn.RemoteAddr().String())
			return
		}

		p := self.NewProtocol()
		l, err := p.Get(lnk.BufPool)
		if err != nil {
			log.Err("[buf not enough]: %s", err.Error())
		}
		c := NewContext()
		c.Link = lnk
		self.routerFunc(c, p)

		lnk.BufPop(l)
	}
}
