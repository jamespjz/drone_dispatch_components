## 插件说明
- [lecheng-drone](https://gitee.com/jamespi/lecheng-drone) - 这是一个无人机控制的核心库，提供了无人机的基本操作接口。
- 插件调用：go get gitee.com/jamespi/drone_dispatch@latest

## 插件调用示例

```azure
package main

import (
	"fmt"
	"my-drone-project/plugin"
	_ "gitee.com/jamespi/lecheng-drone/plugin/plugins" // 自动注册插件
)

func main() {
	// 启用指定插件
	plugin.LoadEnableList([]string{"dji_dock2"})

	// 获取并使用
	drone := plugin.Get("dji_dock2")
	if drone != nil {
		drone.TakeOff()
	} else {
		fmt.Println("插件未启用或不存在")
	}
}

```


## 依赖插件
- go get gopkg.in/yaml.v3 （废弃）
- go get github.com/spf13/viper



