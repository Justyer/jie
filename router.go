package jie

type IRouter interface {
	// 将数据分发到对应的路由上
	Deal(*Context)
	// 执行路由
	Do(*Context, ...interface{})
	// 保存路由与控制器之间的对应
	GET(RouterFunc, ...interface{}) IRouter
}
