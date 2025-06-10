package plugins

import (
	"gitee.com/jamespi/lecheng-drone/plugin"
	"gitee.com/jamespi/lecheng-drone/service"
)

type DJIPilotAdapter struct {
}

/**  无人机业务  **/
// 设备绑定机场
func (d *DJIPilotAdapter) AirportBindStatus() (bool, error) {
	//TODO implement me
	panic("implement me")
}

// 设备绑定机场组织
func (d *DJIPilotAdapter) AirportOrganizationBind() (bool, error) {
	//TODO implement me
	panic("implement me")
}

// 一键起飞
func (d *DJIPilotAdapter) TakeOff() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 一键起飞结果事件通知
func (d *DJIPilotAdapter) TakeOffToPointProgress() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 一键降落
func (d *DJIPilotAdapter) Land() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 一键降落结果事件通知
func (d *DJIPilotAdapter) LandToPointProgress() (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC链路状态通知
func (d *DJIPilotAdapter) DrcStatusNotify() (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC飞行控制无效原因通知
func (d *DJIPilotAdapter) JoystickInvalidNotify() (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC杆量控制
func (d *DJIPilotAdapter) StickControl(stickX, stickY int) (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC飞行器急停
func (d *DJIPilotAdapter) DroneEmergencyStop() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 飞行控制权争夺
func (d *DJIPilotAdapter) FlightAuthorityGrab() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 飞行控制权释放
func (d *DJIPilotAdapter) FlightAuthorityRelease() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 切换相机模式
func (d *DJIPilotAdapter) CameraModeSwitch(mode string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载开始拍照
func (d *DJIPilotAdapter) CameraPhotoTake() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载停止拍照
func (d *DJIPilotAdapter) CameraPhotoStop() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载开始录像
func (d *DJIPilotAdapter) CameraRecordingStart() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载停止录像
func (d *DJIPilotAdapter) CameraRecordingStop() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载框选变焦
func (d *DJIPilotAdapter) CameraFrameZoom(frameX, frameY int) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载变焦
func (d *DJIPilotAdapter) CameraFocalLengthSet(focalLength int) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 画面拖动控制（云台与机身是否一起转动）
func (d *DJIPilotAdapter) CameraScreenDrag(dragX, dragY int, isFollow bool) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台照片存储设置
func (d *DJIPilotAdapter) CameraPhotoStorageSet(storageType string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台录像存储设置
func (d *DJIPilotAdapter) CameraRecordingStorageSet(storageType string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机名称
func (d *DJIPilotAdapter) GetDroneName() string {
	//TODO implement me
	panic("implement me")
}

// 获取无人机类型
func (d *DJIPilotAdapter) GetDroneType() string {
	//TODO implement me
	panic("implement me")
}

// 获取无人机状态
func (d *DJIPilotAdapter) GetDroneStatus() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机位置信息
func (d *DJIPilotAdapter) GetDroneLocation() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机电池状态
func (d *DJIPilotAdapter) GetDroneBatteryLevel() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机摄像头状态
func (d *DJIPilotAdapter) GetDroneCameraStatus() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机飞行时间
func (d *DJIPilotAdapter) GetDroneFlightTime() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机最大飞行高度
func (d *DJIPilotAdapter) GetDroneMaxAltitude() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 设备属性
func (d *DJIPilotAdapter) GetDroneMaxSpeed() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机最大载重
func (d *DJIPilotAdapter) GetDronePayloadCapacity() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机制造商
func (d *DJIPilotAdapter) GetDroneManufacturer() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机型号
func (d *DJIPilotAdapter) GetDroneModel() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机固件版本
func (d *DJIPilotAdapter) GetDroneFirmwareVersion() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机序列号
func (d *DJIPilotAdapter) GetDroneLastMaintenanceDate() (string, error) {
	//TODO implement me
	panic("implement me")
}

/**  mqtt相关  **/
// 是否支持MQTT
func (d *DJIPilotAdapter) SupportsMqtt() bool {
	//TODO implement me
	panic("implement me")
}

// 发布消息到MQTT主题
func (d *DJIPilotAdapter) Subscribe(topic string, callback func(message string)) error {
	//TODO implement me
	panic("implement me")
}

/**  图传相关业务  **/
func (d *DJIPilotAdapter) GetLiveStreamURL() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 实例化 DJIPilotAdapter 并注册到插件系统（自动注册）
func init() {
	plugin.RegisterPlugin("dji_pilot", func() service.DroneAdapter {
		return &DJIPilotAdapter{}
	})
}
