// 定义插件接口与插件元信息

package service

import "io"

// 定义适配器接口
type FH2DroneAdapter interface {
	/**  司空2  **/
	/**  无人机业务  **/
	// 添加设置token的方法
	SetUserToken(token string)
	// 获取组织下的项目列表
	GetprojectList() (string, error)
	// 获取项目下的设备列表
	GetDeviceList(projectUuid string) (string, error)
	// 获取项目的存储上传凭证
	GetProjectStsToken(projectUuid string) (string, error)
	// 获取设备物模型信息
	GetStsToken(projectUuid string, deviceSn string) (string, error)
	// 获取设备HMS信息
	GetDeviceHms(projectUuid string, deviceSnList string) (string, error)
	// 实时控制指令下发
	UpdateDeviceCommand(projectUuid string, deviceSn string, payLoad io.Reader) (string, error)
	// 机场相机切换
	UpdateDeviceChangeCamera(projectUuid string, payLoad io.Reader) (string, error)
	// 设备飞行器镜头切换
	UpdateDeviceChangeLens(projectUuid string, payLoad io.Reader) (string, error)
	// 获取设备控制权
	GetDeviceControl(projectUuid string, payLoad io.Reader) (string, error)
	// 设备控制权释放
	DeleteDeviceControl(projectUuid string, payLoad io.Reader) (string, error)
	// 图传清晰度设置
	UpdateDeviceStreamQuality(projectUuid string, payLoad io.Reader) (string, error)
	// 自定义网络RTK标定
	CreateDeviceRTK(projectUuid string, deviceSn string, payLoad string) (string, error)
	// 开启直播
	LiveStreamStart(projectUuid string, payLoad string) (string, error)
	// 创建飞行任务
	CreateFlightTask(projectUuid string, payLoad string) (string, error)
	// 更新飞行状态
	UpdateFlightTaskStatus(projectUuid string, task_uuid string, payLoad io.Reader) (string, error)
	// 获取飞行任务信息
	GetFlightTaskInfo(projectUuid string, task_uuid string) (string, error)
	// 获取飞行任务列表
	GetFlightTask(projectUuid string, sn string, name string, begin_at int, end_at int, task_type string, status string) (string, error)
	// 获取飞行任务产生的媒体资源
	GetFlightTaskMedia(projectUuid string, task_uuid string) (string, error)
	// 获取飞行任务轨迹信息
	GetFlightTaskTrack(projectUuid string, task_uuid string) (string, error)
	// 航线上传完成通知
	SetFinishUpload(projectUuid string, objectKeyPrefix string, fileName string) (string, error)
	// 获取项目下航线列表
	GetWayLine(projectUuid string) (string, error)
	// 获取项目下的航线详情
	GetWayLineInfo(projectUuid string, wayLineUuid string) (string, error)
	// 模型重建
	CreateModel(projectUuid string, payLoad string) (string, error)
	// 获取模型详情
	GetModelInfo(projectUuid string, modelId int64) (string, error)
	// 获取项目下的模型列表
	GetModelList(projectUuid string) (string, error)
}
