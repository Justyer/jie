package jie

import (
	"fmt"
	"io"
	"net"

	"github.com/Justyer/JekoServer/plugin/log"
	"github.com/Justyer/lingo/ip"
)

type RouterFunc func(*Context)

type Engine struct {
	// 所有连接管理
	conns []net.Conn

	Router IRouter

	Protocol IProtocol
}

func New() *Engine {
	e := &Engine{}
	return e
}

// ************
//  初始化设置
// ************

// 设置自定义协议
func (self *Engine) SetProtocol(p IProtocol) {
	self.Protocol = p
}

// 设置自定义路由
func (self *Engine) SetRouter(r IRouter) {
	self.Router = r
}

// ************
//    监听
// ************

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

// ************
// 每条连接的处理
// ************

func (self *Engine) dealData(conn net.Conn) {
	defer conn.Close()

	lnk := NewLink()
	lnk.Conn = conn

	for {
		if err := lnk.Read(512); err != nil {
			if err == io.EOF {
				log.Err("[link-close-iport]: (%s)", lnk.Conn.RemoteAddr().String())
				return
			}
			log.Err("[Error reading]: %s, on (%s)", err.Error(), lnk.Conn.RemoteAddr().String())
			return
		}

		p := self.Protocol.New()
		l, err := p.Get(lnk.BufPool)
		if err != nil {
			log.Err("[buf not enough]: %v %s", lnk.BufPool, err.Error())
			continue
		}
		c := NewContext()
		c.Link = lnk
		c.DP = p
		c.RT = self.Router
		self.Router.Deal(c)

		lnk.BufPop(l)
	}
}
