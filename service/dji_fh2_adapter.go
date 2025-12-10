// 定义插件接口与插件元信息

package service

import (
	"context"
	"io"
)

// 定义适配器接口
type FH2DroneAdapter interface {
	/**  司空2  **/
	/**  无人机业务  **/
	// GetprojectList 获取组织下的项目列表
	GetprojectList(ctx context.Context) (string, error)
	// GetDeviceList 获取项目下的设备列表
	GetDeviceList(ctx context.Context) (string, error)
	// GetProjectStsToken 获取项目的存储上传凭证
	GetProjectStsToken(ctx context.Context) (string, error)
	// GetStsToken 获取设备物模型信息
	GetStsToken(ctx context.Context, deviceSn string) (string, error)
	// GetDeviceHms 获取设备HMS信息
	GetDeviceHms(ctx context.Context, deviceSnList string) (string, error)
	// UpdateDeviceCommand 实时控制指令下发
	UpdateDeviceCommand(ctx context.Context, deviceSn string, payLoad io.Reader) (string, error)
	// UpdateDeviceChangeCamera 机场相机切换
	UpdateDeviceChangeCamera(ctx context.Context, payLoad io.Reader) (string, error)
	// UpdateDeviceChangeLens 设备飞行器镜头切换
	UpdateDeviceChangeLens(ctx context.Context, payLoad io.Reader) (string, error)
	// GetDeviceControl 获取设备控制权
	GetDeviceControl(ctx context.Context, payLoad io.Reader) (string, error)
	// DeleteDeviceControl 设备控制权释放
	DeleteDeviceControl(ctx context.Context, payLoad io.Reader) (string, error)
	// UpdateDeviceStreamQuality 图传清晰度设置
	UpdateDeviceStreamQuality(ctx context.Context, payLoad io.Reader) (string, error)
	// CreateDeviceRTK 自定义网络RTK标定
	CreateDeviceRTK(ctx context.Context, deviceSn string, payLoad io.Reader) (string, error)
	// LiveStreamStart 开启直播
	LiveStreamStart(ctx context.Context, payLoad io.Reader) (string, error)
	// CreateFlightTask 创建飞行任务
	CreateFlightTask(ctx context.Context, payLoad io.Reader) (string, error)
	// UpdateFlightTaskStatus 更新飞行状态
	UpdateFlightTaskStatus(ctx context.Context, taskUUID string, payLoad io.Reader) (string, error)
	// GetFlightTaskInfo 获取飞行任务信息
	GetFlightTaskInfo(ctx context.Context, taskUUID string) (string, error)
	// GetFlightTask 获取飞行任务列表
	GetFlightTask(ctx context.Context, sn string, name string, beginAt int, endAt int, taskType string, status string) (string, error)
	// GetFlightTaskMedia 获取飞行任务产生的媒体资源
	GetFlightTaskMedia(ctx context.Context, taskUUID string) (string, error)
	// GetFlightTaskTrack 获取飞行任务轨迹信息
	GetFlightTaskTrack(ctx context.Context, taskUUID string) (string, error)
	// SetFinishUpload 航线上传完成通知
	SetFinishUpload(ctx context.Context, objectKeyPrefix string, fileName string) (string, error)
	// GetWayLine 获取项目下航线列表
	GetWayLine(ctx context.Context) (string, error)
	// GetWayLineInfo 获取项目下的航线详情
	GetWayLineInfo(ctx context.Context, wayLineUuid string) (string, error)
	// CreateModel 模型重建
	CreateModel(ctx context.Context, payLoad io.Reader) (string, error)
	// GetModelInfo 获取模型详情
	GetModelInfo(ctx context.Context, modelId int64) (string, error)
	// GetModelList 获取项目下的模型列表
	GetModelList(ctx context.Context) (string, error)
}
