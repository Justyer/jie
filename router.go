package jie

type IRouter interface {
	// 将数据分发到对应的路由上
	Deal(*Context)
	//执行路由对应的控制器
	GET(RouterFunc, ...interface{}) IRouter
}
