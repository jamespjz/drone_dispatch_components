package plugins

import (
	"encoding/json"
	"fmt"
	"gitee.com/jamespi/lecheng-drone/config"
	"gitee.com/jamespi/lecheng-drone/plugin"
	"gitee.com/jamespi/lecheng-drone/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
	"time"
)

// DJI dock 2 适配器
type DJIDock2Adapter struct {
	Client        mqtt.Client
	AccessToken   string           // 访问令牌
	GateWaySn     string           // 网关序列号
	Mutex         sync.RWMutex     // 互斥锁，确保线程安全
	DroneSn       string           // 无人机序列号
	DockSn        string           // 机场序列号
	latestOSD     service.DroneOSD // 存储最新OSD数据
	osdSubscribed bool             // 是否已订阅OSD主题
}

// 初始化链接和订阅
type InitParams struct {
	AccessToken string // 访问令牌
	GateWaySn   string // 网关序列号
	DockSn      string // 无人机序列号
	ClientId    string // 客户端ID
	MqttHost    string // MQTT主机地址
	UserName    string // mqtt用户名
	Password    string // mqtt密码
}

/**  鉴权  **/
func InitializationDJIDock2Adapter(params InitParams) (djiDock2Adapter *DJIDock2Adapter, err error) {

	opts := mqtt.NewClientOptions(). // 设置MQTT客户端选项
						AddBroker(params.MqttHost).   // 添加MQTT代理地址
						SetProtocolVersion(4).        // 设置MQTT协议版本
						SetClientID(params.ClientId). // 设置客户端ID
						SetUsername(params.UserName). // 设置认证用户名
						SetPassword(params.Password). // 设置认证密码
						SetAutoReconnect(true).       // 启用断线自动重连
						SetCleanSession(true)         // 设置清除会话模式

	client := mqtt.NewClient(opts) // 创建新的MQTT客户端

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error() // 如果连接失败，返回错误
	}

	djiDock2Adapter = &DJIDock2Adapter{
		Client:      client,
		AccessToken: params.AccessToken,
		GateWaySn:   params.GateWaySn,
		DockSn:      params.DockSn,
	}

	// 订阅数传OSD
	if err := djiDock2Adapter.subscribeOSD(); err != nil {
		return nil, fmt.Errorf("订阅数传OSD失败: %w", err)
	}

	return djiDock2Adapter, nil
}

// 订阅数传OSD
func (d *DJIDock2Adapter) subscribeOSD() error {
	topic := fmt.Sprintf("thing/product/%s/services/osd", d.DockSn) // 构建OSD主题
	token := d.Client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var osd service.DroneOSD
		if err := json.Unmarshal(msg.Payload(), &osd); err != nil {
			d.Mutex.Lock()    // 锁定互斥锁
			d.latestOSD = osd // 更新最新OSD数据
			d.Mutex.Unlock()  // 解锁互斥锁
		}
	})
	if token.Wait() && token.Error() != nil {
		return token.Error() // 如果订阅失败，返回错误
	}
	d.osdSubscribed = true // 标记已订阅OSD主题
	return nil
}

/**  无人机业务  **/
// 设备绑定机场
func (d *DJIDock2Adapter) AirportBindStatus() (bool, error) {
	//TODO implement me
	panic("implement me")
}

// 设备绑定机场组织
func (d *DJIDock2Adapter) AirportOrganizationBind() (bool, error) {
	//TODO implement me
	panic("implement me")
}

// 一键起飞
func (d *DJIDock2Adapter) TakeOff() (string, error) {
	request := map[string]interface{}{
		"method":      "takeoff_to_point",
		"params":      map[string]string{"sn": d.GateWaySn},
		"timestamp":   time.Now().UnixMilli(),
		"clientToken": d.AccessToken,
	}
	// 发送一键起飞命令到MQTT主题
	topic := fmt.Sprintf("thing/product/%s/services/takeoff", d.GateWaySn)
	payload, _ := json.Marshal(request)
	token := d.Client.Publish(topic, 1, false, payload)
	token.Wait()                          // 等待消息发送完成
	return string(payload), token.Error() // 返回发送结果或错误信息
}

// 一键起飞结果事件通知
func (d *DJIDock2Adapter) TakeOffToPointProgress() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 一键降落
func (d *DJIDock2Adapter) Land() (string, error) {
	request := map[string]interface{}{
		"method":      "land_to_point",
		"params":      map[string]string{"sn": d.GateWaySn},
		"timestamp":   time.Now().UnixMilli(),
		"clientToken": d.AccessToken,
	}
	// 发送一键降落命令到MQTT主题
	topic := fmt.Sprintf("thing/product/%s/services/land", d.DockSn)
	payload, _ := json.Marshal(request)
	token := d.Client.Publish(topic, 1, false, payload)
	token.Wait() // 等待消息发送完成

	return string(payload), token.Error() // 返回发送结果或错误信息
}

// 一键降落结果事件通知
func (d *DJIDock2Adapter) LandToPointProgress() (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC链路状态通知
func (d *DJIDock2Adapter) DrcStatusNotify() (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC飞行控制无效原因通知
func (d *DJIDock2Adapter) JoystickInvalidNotify() (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC杆量控制
func (d *DJIDock2Adapter) StickControl(stickX, stickY int) (string, error) {
	//TODO implement me
	panic("implement me")
}

// DRC飞行器急停
func (d *DJIDock2Adapter) DroneEmergencyStop() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 飞行控制权争夺
func (d *DJIDock2Adapter) FlightAuthorityGrab() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 飞行控制权释放
func (d *DJIDock2Adapter) FlightAuthorityRelease() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 切换相机模式
func (d *DJIDock2Adapter) CameraModeSwitch(mode string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载开始拍照
func (d *DJIDock2Adapter) CameraPhotoTake() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载停止拍照
func (d *DJIDock2Adapter) CameraPhotoStop() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载开始录像
func (d *DJIDock2Adapter) CameraRecordingStart() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载停止录像
func (d *DJIDock2Adapter) CameraRecordingStop() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载框选变焦
func (d *DJIDock2Adapter) CameraFrameZoom(frameX, frameY int) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台挂载变焦
func (d *DJIDock2Adapter) CameraFocalLengthSet(focalLength int) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 画面拖动控制（云台与机身是否一起转动）
func (d *DJIDock2Adapter) CameraScreenDrag(dragX, dragY int, isFollow bool) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台照片存储设置
func (d *DJIDock2Adapter) CameraPhotoStorageSet(storageType string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 云台录像存储设置
func (d *DJIDock2Adapter) CameraRecordingStorageSet(storageType string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机名称
func (d *DJIDock2Adapter) GetDroneName() string {
	//TODO implement me
	panic("implement me")
}

// 获取无人机类型
func (d *DJIDock2Adapter) GetDroneType() string {
	//TODO implement me
	panic("implement me")
}

// 获取无人机状态
func (d *DJIDock2Adapter) GetDroneStatus() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机位置信息
func (d *DJIDock2Adapter) GetDroneLocation() (string, error) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock() // 确保在访问最新OSD数据时锁定互斥锁
	location := fmt.Sprintf("Latitude: %f, Longitude: %f, Altitude: %f",
		d.latestOSD.Latitude,
		d.latestOSD.Longitude,
		d.latestOSD.Altitude)

	return location, nil
}

// 获取无人机电池状态
func (d *DJIDock2Adapter) GetDroneBatteryLevel() (int, error) {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.latestOSD.Battery, nil
}

// 获取无人机摄像头状态
func (d *DJIDock2Adapter) GetDroneCameraStatus() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机飞行时间
func (d *DJIDock2Adapter) GetDroneFlightTime() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机最大飞行高度
func (d *DJIDock2Adapter) GetDroneMaxAltitude() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 设备属性
func (d *DJIDock2Adapter) GetDroneMaxSpeed() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机最大载重
func (d *DJIDock2Adapter) GetDronePayloadCapacity() (int, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机制造商
func (d *DJIDock2Adapter) GetDroneManufacturer() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机型号
func (d *DJIDock2Adapter) GetDroneModel() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机固件版本
func (d *DJIDock2Adapter) GetDroneFirmwareVersion() (string, error) {
	//TODO implement me
	panic("implement me")
}

// 获取无人机序列号
func (d *DJIDock2Adapter) GetDroneLastMaintenanceDate() (string, error) {
	//TODO implement me
	panic("implement me")
}

/**  mqtt相关  **/
// 是否支持MQTT
func (d *DJIDock2Adapter) SupportsMqtt() bool {
	return true
}

// 发布消息到MQTT主题
func (d *DJIDock2Adapter) Subscribe(topic string, callback func(message string)) error {
	return d.Client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		callback(string(msg.Payload()))
	}).Error()
}

/**  图传相关业务  **/
func (d *DJIDock2Adapter) GetLiveStreamURL() (string, error) {

	return config.RtmpURLSettings, nil

}

// 实例化 DJIDock2Adapter 并注册到插件系统（自动注册）
func init() {
	// 注册 DJI Dock2 适配器
	plugin.RegisterPlugin("dji_dock2", func() service.DroneAdapter {
		//从配置获取参数
		tm := plugin.NewTokenManager("dji_dock2")
		params := InitParams{
			AccessToken: tm.GetAccessToken(),                                                        // 获取访问令牌
			GateWaySn:   config.DjiSettings["GatewaySn"],                                            // 获取网关序列号
			DockSn:      config.DjiSettings["DockSn"],                                               // 获取机场序列号
			ClientId:    config.DjiSettings["ClientId"],                                             // 获取客户端ID
			MqttHost:    "tcp://" + config.MqttSettings["host"] + ":" + config.MqttSettings["port"], // 获取MQTT主机地址
			UserName:    config.MqttSettings["username"],                                            // 获取MQTT用户名
			Password:    config.MqttSettings["password"],                                            // 获取MQTT密码
		}
		// 设置令牌
		tm.SetAccessToken(params.AccessToken, "", config.TokenExpiresInSettings) // 设置访问令牌，刷新令牌和过期时间

		adapter, err := InitializationDJIDock2Adapter(params)
		if err != nil {
			log.Fatalf("DJI Dock2初始化失败: %v", err)
		}
		// 启动令牌刷新协程
		go adapter.startTokenRefreshScheduler(tm)
		return adapter
	})
}

// 令牌刷新调度器
func (d *DJIDock2Adapter) startTokenRefreshScheduler(tm *plugin.TokenManager) {
	ticker := time.NewTicker(5 * time.Minute) // 每5分钟刷新一次令牌
	defer ticker.Stop()
	for range ticker.C {
		_, err := tm.RefreshAccessToken()
		if err != nil {
			log.Printf("令牌刷新失败: %v", err)
		}
	}
}
