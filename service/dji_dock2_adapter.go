// 定义插件接口与插件元信息

package service

// 新增OSD数据结构
type DroneOSD struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Altitude    float64 `json:"altitude"`
	Battery     int     `json:"battery"`
	Speed       float64 `json:"speed"`
	FlightState int     `json:"flight_state"` // 飞行状态
}

// 定义适配器接口
type DJIDock2DroneAdapter interface {
	/**  上云api  **/
	/**  无人机业务  **/
	// 设备绑定机场
	AirportBindStatus() (bool, error) // 是否绑定机场
	// 设备绑定机场组织
	AirportOrganizationBind() (bool, error) // 是否绑定机场组织
	// 一键起飞
	TakeOff() (string, error)
	// 一键起飞结果事件通知
	TakeOffToPointProgress() (string, error)
	// 一键降落
	Land() (string, error)
	// 一键降落结果事件通知
	LandToPointProgress() (string, error)
	// DRC链路状态通知
	DrcStatusNotify() (string, error)
	// DRC飞行控制无效原因通知
	JoystickInvalidNotify() (string, error)
	// DRC杆量控制
	StickControl(stickX, stickY int) (string, error)
	// DRC飞行器急停
	DroneEmergencyStop() (string, error)
	// 飞行控制权争夺
	FlightAuthorityGrab() (string, error)
	// 飞行控制权释放
	FlightAuthorityRelease() (string, error)
	// 切换相机模式
	CameraModeSwitch(mode string) (string, error)
	// 云台挂载开始拍照
	CameraPhotoTake() (string, error)
	// 云台挂载停止拍照
	CameraPhotoStop() (string, error)
	// 云台挂载开始录像
	CameraRecordingStart() (string, error)
	// 云台挂载停止录像
	CameraRecordingStop() (string, error)
	// 云台挂载框选变焦
	CameraFrameZoom(frameX, frameY int) (string, error)
	// 云台挂载变焦
	CameraFocalLengthSet(focalLength int) (string, error)
	// 画面拖动控制（云台与机身是否一起转动）
	CameraScreenDrag(dragX, dragY int, isFollow bool) (string, error)
	// 云台照片存储设置
	CameraPhotoStorageSet(storageType string) (string, error)
	// 云台录像存储设置
	CameraRecordingStorageSet(storageType string) (string, error)
	// 获取无人机名称
	GetDroneName() string
	// 获取无人机类型
	GetDroneType() string
	// 获取无人机状态
	GetDroneStatus() (string, error)
	// 获取无人机位置信息
	GetDroneLocation() (string, error)
	// 获取无人机电池状态
	GetDroneBatteryLevel() (int, error)
	// 获取无人机摄像头状态
	GetDroneCameraStatus() (string, error)
	// 获取无人机飞行时间
	GetDroneFlightTime() (int, error)
	// 获取无人机最大飞行高度
	GetDroneMaxAltitude() (int, error)
	// 设备属性
	GetDroneMaxSpeed() (int, error)
	// 获取无人机最大载重
	GetDronePayloadCapacity() (int, error)
	// 获取无人机制造商
	GetDroneManufacturer() (string, error)
	// 获取无人机型号
	GetDroneModel() (string, error)
	// 获取无人机固件版本
	GetDroneFirmwareVersion() (string, error)
	// 获取无人机序列号
	GetDroneLastMaintenanceDate() (string, error)
	/**  mqtt相关  **/
	//是否支持MQTT
	SupportsMqtt() bool
	/**  图传  **/
	// 获取图传流地址
	GetLiveStreamURL() (string, error)
}
