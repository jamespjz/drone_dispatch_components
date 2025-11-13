package main

import (
	"fmt"
	"log"

	"gitee.com/jamespi/drone_dispatch/config"
	"gitee.com/jamespi/drone_dispatch/plugin"
	_ "gitee.com/jamespi/drone_dispatch/plugin/plugins" // 自动注册插件
	"gitee.com/jamespi/drone_dispatch/service"
)

func main() {
	// 初始化配置 - 这是必须的第一步
	if err := config.InitDefaultConfig(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	// 司空2调用
	// 启用指定插件
	plugin.LoadEnableList([]string{"fh2", "dji_dock2"})
	// 获取并使用
	if fh2, ok := plugin.Get[service.FH2DroneAdapter](plugin.FH2Plugin); ok {
		// 调用适配器方法
		// 获取组织下的项目列表
		projectList, err := fh2.GetprojectList()
		if err != nil {
			fmt.Println("获取项目列表失败:", err)
			return
		}
		fmt.Println("项目列表:", projectList)
		// 获取项目下的设备列表
		//deviceList, err := fh2.GetDeviceList("c33595a4-3996-481d-9d81-459d435ade84")
		//if err != nil {
		//	fmt.Println("获取设备列表失败:", err)
		//	return
		//}
		//fmt.Println("设备列表:", deviceList)
		//// 获取设备HMS信息
		//stsTokenList, err := fh2.GetDeviceHms("c33595a4-3996-481d-9d81-459d435ade84", "7CTXN4A00B096H和,1581F6QAD241600B4FP8")
		//if err != nil {
		//	fmt.Println("获取设备HMS列表失败:", err)
		//	return
		//}
		//fmt.Println("HMS列表:", stsTokenList)
		//// 获取设备物模型信息
		//stsTokenList, err := fh2.GetStsToken("c33595a4-3996-481d-9d81-459d435ade84", "7CTXN4A00B096H")
		//if err != nil {
		//	fmt.Println("获取设备物模型列表失败:", err)
		//	return
		//}
		//fmt.Println("物模型列表:", stsTokenList)
		// 获取飞行任务信息
		//taskTrackList, err := fh2.GetFlightTaskInfo("c33595a4-3996-481d-9d81-459d435ade84", "d542b848-6f0a-46ac-b072-8f68af7531c31")
		//if err != nil {
		//	fmt.Println("获取设备任务轨迹信息失败:", err)
		//	return
		//}
		//fmt.Println("任务轨迹信息:", taskTrackList)
		// 获取飞行任务信息
		//flightTaskInfo, err := fh2.GetFlightTaskInfo("c33595a4-3996-481d-9d81-459d435ade84", "4143bf37-63a1-42bd-85e9-a7ce48003947")
		//if err != nil {
		//	fmt.Println("获取飞行任务信息失败:", err)
		//	return
		//}
		//fmt.Println("获取飞行任务信息:", flightTaskInfo)
		// 创建飞行任务
		//	payLoad := `{
		//  "name": "测试任务",
		//  "wayline_uuid": "6d88fbe5-a399-485a-86ba-7bbdbb99edec",
		//  "sn": "7CTXN4A00B096H",
		//  "rth_altitude": 80,
		//  "rth_mode": "optimal",
		//  "wayline_precision_type": "gps",
		//  "out_of_control_action_in_flight": "return_home",
		//  "resumable_status": "auto",
		//  "task_type": "immediate",
		//  "time_zone": "Asia/Shanghai",
		//  "repeat_type": "nonrepeating",
		//  "min_battery_capacity": 60
		//}`
		//	taskStatus, err := fh2.CreateFlightTask("c33595a4-3996-481d-9d81-459d435ade84", payLoad)
		//	if err != nil {
		//		fmt.Println("获取设备物模型列表失败:", err)
		//		return
		//	}
		//	fmt.Println("飞行任务状态:", taskStatus)
		// 获取飞行任务信息
		// 模型重建
		//payLoad := `{
		//  "name": "测试任务",
		//  "wayline_uuid": "6d88fbe5-a399-485a-86ba-7bbdbb99edec",
		//  "sn": "7CTXN4A00B096H",
		//  "rth_altitude": 80,
		//  "rth_mode": "optimal",
		//  "wayline_precision_type": "gps",
		//  "out_of_control_action_in_flight": "return_home",
		//  "resumable_status": "auto",
		//  "task_type": "immediate",
		//  "time_zone": "Asia/Shanghai",
		//  "repeat_type": "nonrepeating",
		//  "min_battery_capacity": 60
		//}`
		//taskStatus, err := fh2.CreateModel("c33595a4-3996-481d-9d81-459d435ade84", payLoad)
		//if err != nil {
		//	fmt.Println("获取设备物模型列表失败:", err)
		//	return
		//}
		//fmt.Println("设备物模型:", taskStatus)
		// 获取项目的存储上传凭证
		stsTokenInfo, err := fh2.GetProjectStsToken("c33595a4-3996-481d-9d81-459d435ade84")
		if err != nil {
			fmt.Println("获取存储上传凭证失败:", err)
			return
		}
		fmt.Println("获取存储上传凭证信息:", stsTokenInfo)

	} else {
		fmt.Println("插件未启用或不存在")
	}
	// 机场2调用
	// 获取并使用
	//if dock2, ok := plugin.Get[service.DJIDock2DroneAdapter](plugin.DJIDock2Plugin); ok {
	//	// 调用适配器方法
	//	// 一键起飞
	//	projectList, err := dock2.TakeOff()
	//	if err != nil {
	//		fmt.Println("获取设备物模型列表失败:", err)
	//		return
	//	}
	//	fmt.Println("起飞状态:", projectList)
	//} else {
	//	fmt.Println("插件未启用或不存在")
	//}
}
