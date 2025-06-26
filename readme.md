## 插件说明
- [drone_dispatch](https://gitee.com/jamespi/drone_dispatch) - 这是一个无人机控制的核心库，提供了无人机的基本操作接口。
- 插件调用：go get gitee.com/jamespi/drone_dispatch@latest

## 插件调用示例

```azure
package main

import (
	"fmt"
	"gitee.com/jamespi/drone_dispatch/plugin"
	_ "gitee.com/jamespi/drone_dispatch/plugin/plugins" // 自动注册插件
	"gitee.com/jamespi/drone_dispatch/service"
)

func main() {
	// 启用指定插件
	plugin.LoadEnableList([]string{"fh2"})
	// 获取并使用
if fh2, ok := plugin.Get[service.FH2DroneAdapter](plugin.FH2Plugin); ok {
		// 调用适配器方法
		projectList, err := fh2.GetprojectList()
if err != nil {
			fmt.Println("获取项目列表失败:", err)
    return
        }
		fmt.Println("项目列表:", projectList)

		deviceList, err := fh2.GetDeviceList("example-project-uuid")
if err != nil {
			fmt.Println("获取设备列表失败:", err)
    return
        }
		fmt.Println("设备列表:", deviceList)
	} else {
		fmt.Println("插件未启用或不存在")
	}

}


```


## 依赖插件
- go get gopkg.in/yaml.v3 （废弃）
- go get github.com/spf13/viper



