package jie

// IProtocol ： 协议接口，实现此接口就可以自定义协议
type IProtocol interface {
	// 判断该字节数据是不是一个完整的协议
	IsWhole([]byte) bool
	// 解析协议
	Parse([]byte)
	// 将数据流转化为数据包，依赖IsWhole和Parse两个方法
	Get([]byte) (int, error)
	// 实例化
	New() IProtocol
	// 获取数据包中的数据字段
	Data() []byte
}
