// 定义插件接口与插件元信息

package service

// 定义适配器接口
type FH2DroneAdapter interface {
	/**  司空2  **/
	/**  无人机业务  **/
	// 获取组织下的项目列表
	GetprojectList() (string, error)
	// 获取项目下的设备列表
	GetDeviceList(projectUuid string) (string, error)
	// 获取设备HMS信息
	GetDeviceHms(projectUuid string, deviceSnList string) (string, error)
	// 获取设备物模型信息
	GetStsToken(projectUuid string, deviceSn string) (string, error)
	// 航线上传
	SetFinishUpload(projectUuid string, objectKeyPrefix string, fileName string) (string, error)
	// 创建飞行任务
	CreateFlightTask(projectUuid string, payLoad string) (string, error)
	// 获取项目下航线列表
	GetWayLine(projectUuid string) (string, error)
	// 获取项目下的航线详情
	GetWayLineInfo(projectUuid string, wayLineUuid string) (string, error)
}
