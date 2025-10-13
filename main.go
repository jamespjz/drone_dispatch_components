package main

import (
	"fmt"
	"gitee.com/jamespi/drone_dispatch/plugin"
	_ "gitee.com/jamespi/drone_dispatch/plugin/plugins" // 自动注册插件
	"gitee.com/jamespi/drone_dispatch/service"
)

func main() {
	// 司空2调用
	// 启用指定插件
	plugin.LoadEnableList([]string{"fh2"})
	//	// 获取并使用
	//	if fh2, ok := plugin.Get[service.FH2DroneAdapter](plugin.FH2Plugin); ok {
	//		// 调用适配器方法
	//		// 获取组织下的项目列表
	//		projectList, err := fh2.GetprojectList()
	//		if err != nil {
	//			fmt.Println("获取项目列表失败:", err)
	//			return
	//		}
	//		fmt.Println("项目列表:", projectList)
	//		// 获取项目下的设备列表
	//		deviceList, err := fh2.GetDeviceList("c33595a4-3996-481d-9d81-459d435ade84")
	//		if err != nil {
	//			fmt.Println("获取设备列表失败:", err)
	//			return
	//		}
	//		fmt.Println("设备列表:", deviceList)
	//		// 获取设备物模型信息
	//		stsTokenList, err := fh2.GetStsToken("c33595a4-3996-481d-9d81-459d435ade84", "7CTXN4A00B096H")
	//		if err != nil {
	//			fmt.Println("获取设备物模型列表失败:", err)
	//			return
	//		}
	//		fmt.Println("物模型列表:", stsTokenList)
	//		// 创建飞行任务
	//		payLoad := `{
	//    "name": "测试任务",
	//    "wayline_uuid": "6d88fbe5-a399-485a-86ba-7bbdbb99edec",
	//    "sn": "7CTXN4A00B096H",
	//    "rth_altitude": 80,
	//    "rth_mode": "optimal",
	//    "wayline_precision_type": "gps",
	//    "out_of_control_action_in_flight": "return_home",
	//    "resumable_status": "auto",
	//    "task_type": "immediate",
	//    "time_zone": "Asia/Shanghai",
	//    "repeat_type": "nonrepeating",
	//    "min_battery_capacity": 60
	//}`
	//		taskStatus, err := fh2.CreateFlightTask("c33595a4-3996-481d-9d81-459d435ade84", payLoad)
	//		if err != nil {
	//			fmt.Println("获取设备物模型列表失败:", err)
	//			return
	//		}
	//		fmt.Println("飞行任务状态:", taskStatus)
	//
	//	} else {
	//		fmt.Println("插件未启用或不存在")
	//	}
	// 机场2调用
	// 获取并使用
	if dock2, ok := plugin.Get[service.DJIDock2DroneAdapter](plugin.DJIDock2Plugin); ok {
		// 调用适配器方法
		// 一键起飞
		projectList, err := dock2.TakeOff()
		if err != nil {
			fmt.Println("获取设备物模型列表失败:", err)
			return
		}
		fmt.Println("起飞状态:", projectList)
	} else {
		fmt.Println("插件未启用或不存在")
	}
}
