package protocol

import (
	"github.com/Justyer/JekoServer/plugin/log"
	"github.com/Justyer/jie"
)

// Node : 路由树节点
type Node struct {
	Value uint16
	Func  jie.RouterFunc
	Node  map[uint16]*Node
}

// Router : 自定义路由
// 由协议和路由树组成
type Router struct {
	// 路由对应的协议
	Protocol jie.IProtocol
	// 路由树
	Tree map[uint16]*Node
}

// NewRouter : 实例化
func NewRouter() *Router {
	r := &Router{}
	r.Tree = make(map[uint16]*Node)
	return r
}

// Forward : 当请求来临时，转发给对应的路由
func (r *Router) Forward(c *jie.Context) {
	dp := c.DP.(*Protocol)
	msgType := dp.MsgType
	msgCmd := dp.MsgCmd

	r.Do(c, msgType, msgCmd)
}

// Do : 执行对应路由
func (r *Router) Do(c *jie.Context, rs ...interface{}) {
	i, l := 0, len(rs)-1
	tree := r.Tree
	var f jie.RouterFunc
	for {
		msg := rs[i].(uint16)
		node, ok := tree[msg]
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

// GET : 保存路由与控制器之间的关系
func (r *Router) GET(f jie.RouterFunc, rs ...interface{}) jie.IRouter {
	i, l := 0, len(rs)-1
	tree := r.Tree
	for {
		msg := rs[i].(uint16)
		node, ok := tree[msg]
		if i == l {
			if ok {
				node.Value = msg
				node.Func = f
				node.Node = make(map[uint16]*Node)
				break
			} else {
				var n Node
				n.Value = msg
				n.Func = f
				n.Node = make(map[uint16]*Node)
				tree[msg] = &n
				break
			}
		}
		i++
	}
	return r
}

// SetProtocol : 设置协议
func (r *Router) SetProtocol(p jie.IProtocol) {
	r.Protocol = p
}

// NoRouter : 当路由不存在时的默认返回方法
func NoRouter(c *jie.Context) {
	p := c.DP.(*Protocol)
	log.Tx("[no_router]: %d %d", p.MsgType, p.MsgCmd)
}
