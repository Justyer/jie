# 自定义协议框架

寒蝉鸣泣之时：解

![jie](cover.jpg)

## 使用

```go
package main

import (
	"github.com/Justyer/jie"
	ptl "github.com/Justyer/jie/protocol"
)

func main() {
	j := jie.New()

	// 设置协议
	j.SetProtocol(func() jie.Protocol {
		return ptl.NewDataPack_2_2_4()
	})

	// 设定路由的具体实现
	j.SetRouter(func(c *jie.Context, p jie.Protocol) {
		// dp := p.(*ptl.DataPack_2_2_4)

		// 这里写路由业务
	})

	// 开启内网监听
	j.ListenAndInnerServe("9595")
}
```