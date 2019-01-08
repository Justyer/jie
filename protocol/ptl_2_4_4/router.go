package protocol

import (
	"github.com/Justyer/JekoServer/plugin/log"
	"github.com/Justyer/jie"
)

type Node struct {
	Value uint16
	Func  jie.RouterFunc
	Node  map[uint16]*Node
}

type Router struct {
	Protocol jie.IProtocol
	Tree     map[uint16]*Node
}

func NewRouter() *Router {
	r := &Router{}
	r.Tree = make(map[uint16]*Node)
	return r
}

func (self *Router) Deal(c *jie.Context) {
	dp := c.DP.(*Protocol)
	msg_type := dp.MsgType
	msg_cmd := dp.MsgCmd

	var rs = []uint16{
		msg_type, msg_cmd,
	}

	var i int = 0
	var l int = len(rs) - 1
	tree := self.Tree
	var f jie.RouterFunc
	for {
		r := rs[i]
		node, ok := tree[r]
		if i == l {
			if ok {
				f = node.Func
			} else {
				f = NoRouter
			}
			break
		}
		i++
	}
	f(c)
}

func (self *Router) GET(f jie.RouterFunc, rs ...interface{}) jie.IRouter {
	var i int = 0
	var l int = len(rs) - 1
	tree := self.Tree
	for {
		r := rs[i].(uint16)
		node, ok := tree[r]
		if i == l {
			if ok {
				node.Value = r
				node.Func = f
				node.Node = make(map[uint16]*Node)
				break
			} else {
				var n Node
				n.Value = r
				n.Func = f
				n.Node = make(map[uint16]*Node)
				tree[r] = &n
				break
			}
		}
		i++
	}
	return self
}

func (self *Router) SetProtocol(p jie.IProtocol) {
	self.Protocol = p
}

func NoRouter(c *jie.Context) {
	p := c.DP.(*Protocol)
	log.Tx("[no_router]: %d %d", p.MsgType, p.MsgCmd)
}
