package jie

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Justyer/JekoServer/plugin/log"
	"github.com/Justyer/lingo/ip"
)

var (
	// MaxBufferSize : 从TCP缓冲区一次读取数据的最大容量
	MaxBufferSize = 512

	// HeartBeatInterval : 设置心跳间隔
	HeartBeatInterval = 6 * time.Second
)

// RouterFunc : 自定义路由
type RouterFunc func(*Context)

// Engine : 主引擎
type Engine struct {
	// 所有连接管理
	conns []net.Conn

	Router IRouter

	Protocol IProtocol
}

// New : 实例化引擎
func New() *Engine {
	e := &Engine{}
	return e
}

// ************
//  初始化设置
// ************

// SetProtocol : 设置自定义协议
func (e *Engine) SetProtocol(p IProtocol) {
	e.Protocol = p
}

// SetRouter : 设置自定义路由
func (e *Engine) SetRouter(r IRouter) {
	e.Router = r
}

// ************
//    监听
// ************

// ListenAndLocalServe : 本地监听
func (e *Engine) ListenAndLocalServe(port string) {
	e.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port))
}

// ListenAndInnerServe : 内网监听
func (e *Engine) ListenAndInnerServe(port string) {
	ip := ip.MustInnerIP()
	e.ListenAndServe(fmt.Sprintf("%s:%s", ip, port))
}

// ListenAndServe : 自定义地址监听
func (e *Engine) ListenAndServe(addr string) {
	listener, err := net.Listen("tcp", addr)
	defer listener.Close()
	if err != nil {
		return
	}
	log.Info("[listener]: (%s)", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Err("[error accept]: %s", err.Error())
			continue
		}

		log.Info("[link start]: (%s)", conn.RemoteAddr().String())
		e.conns = append(e.conns, conn)
		go e.dealData(conn)
	}
}

// ************
// 每条连接的处理
// ************

func (e *Engine) dealData(conn net.Conn) {
	defer conn.Close()

	lnk := NewLink()
	lnk.Conn = conn

	for {
		thisB, err := lnk.Read(MaxBufferSize)
		if err != nil {
			if err == io.EOF {
				log.Err("[link close]: (%s)", lnk.Conn.RemoteAddr().String())
				return
			}
			log.Err("[error read]: %s, on (%s)", err.Error(), lnk.Conn.RemoteAddr().String())
			return
		}

		go e.heartbeat(conn, thisB, HeartBeatInterval)

		// 一次读取的缓冲区数据可能包含多个数据包，所以要循环处理
		for {
			p := e.Protocol.New()
			l, err := p.Get(lnk.BufPool)
			if err != nil {
				log.Err("[buf not enough]: %v %s", lnk.BufPool, err.Error())
				break
			}
			c := NewContext()
			c.Link = lnk
			c.DP = p
			c.RT = e.Router
			go e.Router.Forward(c)

			lnk.BufPop(l)
		}
	}
}

// heartbeat : 心跳
func (e *Engine) heartbeat(conn net.Conn, bs []byte, timeout time.Duration) {
	msg := make(chan byte)
	select {
	case <-msg:
		conn.SetDeadline(time.Now().Add(timeout))
		break
	case <-time.After(timeout):
		log.Err("[link cut]: heart attack")
		conn.Close()
	}
	msg <- bs[0]
	close(msg)
}
