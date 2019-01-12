package main

import (
	"github.com/Justyer/jie"
	protocol "github.com/Justyer/jie/protocol/ptl_2_2_4"
)

func main() {
	j := jie.New()

	// 设置协议
	j.SetProtocol(protocol.NewProtocol())

	// 设置路由
	j.SetRouter(protocol.NewRouter())

	// 添加路由
	j.Router.GET(func(c *jie.Context) {
		c.Send([]byte("JekoWorld~"))
	}, uint16(1), uint16(2))

	// 开启内网监听
	j.ListenAndInnerServe("9595")
}
