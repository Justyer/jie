package jie

// IRouter ： 路由接口，实现此接口就可以自定义路由
type IRouter interface {
	// 将数据分发到对应的路由上
	Forward(*Context)
	// 执行路由
	Do(*Context, ...interface{})
	// 保存路由与控制器之间的对应
	GET(RouterFunc, ...interface{}) IRouter
}
