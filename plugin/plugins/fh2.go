package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/jamespi/drone_dispatch/config"
	"gitee.com/jamespi/drone_dispatch/pkg/httpclient"
	"gitee.com/jamespi/drone_dispatch/pkg/tenant"
	"gitee.com/jamespi/drone_dispatch/pkg/validator"
	"gitee.com/jamespi/drone_dispatch/plugin"
	"gitee.com/jamespi/drone_dispatch/service"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
)

// TenantInfo 租户上下文基础信息
type TenantInfo struct {
	XUserToken  string `json:"xUserToken"`
	ProjectUUID string `json:"project_uuid"`
	TenantId    int64  `json:"tenant_id"`
}

// FH2Adapter 司空2适配器
// 司空2openapi接口文档地址：https://apifox.com/apidoc/shared/6b4ca90b-233f-48ac-818c-d694acb0663a/api-221842037
type FH2Adapter struct {
	tenantInfo   *TenantInfo
	validator    *validator.InputValidator
	secureClient *httpclient.SecureHTTPClient
	mu           sync.RWMutex // 并发线程安全
}

// NewFH2Adapter 创建一个新的FH2适配器
func NewFH2Adapter() *FH2Adapter {
	return &FH2Adapter{
		tenantInfo: &TenantInfo{},
		validator:  validator.GetValidator(),
	}
}

// APIResponse 定义与API响应对应的结构体
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// doRequestWithTenant 基于租户执行http请求
func (F *FH2Adapter) doRequestWithTenant(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	// 从上下文获取租户信息
	tenantInfo, err := tenant.GetTenantFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取租户信息失败: %w", err)
	}

	F.mu.RLock()
	F.tenantInfo = &TenantInfo{
		XUserToken:  tenantInfo.UserToken,
		TenantId:    tenantInfo.TenantId,
		ProjectUUID: tenantInfo.ProjectUUID,
	}
	F.mu.RUnlock()

	// 构建header请求头
	headers := map[string]string{
		"X-User-Token": F.tenantInfo.XUserToken,
		"X-Request-Id": tenant.GetRequestIDFromContext(ctx),
		"X-Language":   "zh",
	}
	if F.tenantInfo.ProjectUUID != "" {
		headers["X-Project-Uuid"] = F.tenantInfo.ProjectUUID
	}
	// 如果请求ID为空，生成一个
	if headers["X-Request-Id"] == "" {
		headers["X-Request-Id"] = uuid.New().String()
	}

	// 执行请求
	result, err := F.doRequestWithHeaders(ctx, method, url, body, headers)

	return result, err

}

// doRequestWithHeaders 执行指定请求头的HTTP请求
func (F *FH2Adapter) doRequestWithHeaders(ctx context.Context, method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	// 执行安全的HTTP请求
	resp, err := F.secureClient.DoRequest(ctx, method, url, body, headers)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP请求返回错误状态码: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 解析JSON并检查业务状态码
	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %w\n原始响应: %s", err, string(bodyBytes))
	}

	// 检查业务代码，如果大于0表示错误
	if apiResp.Code > 0 {
		return nil, fmt.Errorf("API错误[%d]: %s", apiResp.Code, apiResp.Message)
	}

	// 业务代码为0或负数表示成功，返回原始响应体
	return bodyBytes, nil
}

// GetprojectList 获取组织下的项目列表
func (F *FH2Adapter) GetprojectList(ctx context.Context) (string, error) {
	encodedQ := url.QueryEscape(config.FH2Settings["q"])
	if err := F.validator.ValidateQueryParam(encodedQ); err != nil {
		return "", fmt.Errorf("查询参数验证失败: %w", err)
	}

	url := fmt.Sprintf("%s/openapi/v0.1/project?page=1&page_size=10&q=%s&prj_authorized_status=project-status-authorized&usage=simple&sort_column=created_at&sort_type=ASC", config.FH2Settings["host"], encodedQ)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetDeviceList 获取项目下的设备列表
func (F *FH2Adapter) GetDeviceList(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/project/device", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetProjectStsToken 获取项目的存储上传凭证
func (F *FH2Adapter) GetProjectStsToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/project/sts-token", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetStsToken 获取设备物模型信息
func (F *FH2Adapter) GetStsToken(ctx context.Context, deviceSn string) (string, error) {
	if err := F.validator.ValidateDeviceSN(deviceSn); err != nil {
		return "", fmt.Errorf("设备序列号验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/device/%s/state", config.FH2Settings["host"], deviceSn)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetDeviceHms 获取设备HMS信息
func (F *FH2Adapter) GetDeviceHms(ctx context.Context, deviceSnList string) (string, error) {
	encodedQ := url.QueryEscape(deviceSnList)
	if err := F.validator.ValidateQueryParam(encodedQ); err != nil {
		return "", fmt.Errorf("查询参数验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/device/hms?device_sn_list=%s", config.FH2Settings["host"], encodedQ)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// UpdateDeviceCommand 实时控制指令下发
// 在航线飞行或者手动飞行过程中，可以调用此命令对飞行器进行实时控制，比如返航，取消返航，飞行任务暂停和恢复等操作。
//
//	payLoad := `{
//					  "device_command": "return_home", // 控制指令，支持的指令有：return_home（返航）、return_specific_home（蛙跳任务指定目标机场返航）、return_home_cancel（取消返航）、flighttask_pause（飞行任务暂停）、flighttask_recovery（飞行任务恢复）
//					}`
func (F *FH2Adapter) UpdateDeviceCommand(ctx context.Context, deviceSn string, payLoad io.Reader) (string, error) {
	if err := F.validator.ValidateDeviceSN(deviceSn); err != nil {
		return "", fmt.Errorf("设备序列号验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/device/%s/command", config.FH2Settings["host"], deviceSn)
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// UpdateDeviceChangeCamera 机场相机切换
// 支持机场相机的切换，可切换舱内和舱外相机。
// camera_index 和 camera_position 信息可通过获取组织下的设备列表接口下的 data.list.gateway:.camera_list 中获取。
//
//	payLoad := `{
//					  "sn": "1581F6Q8D242100CPWEK", // 设备编码
//					  "camera_index": "165-0-7",  // 相机索引，0表示舱外相机，1表示舱内相机
//					  "camera_position": "indoor",  // 相机位置，indoor表示舱内摄像头，outdoor表示舱外摄像头
//				   }`
func (F *FH2Adapter) UpdateDeviceChangeCamera(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/device/change-camera", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// UpdateDeviceChangeLens 设备飞行器镜头切换
// 支持飞行器相机镜头的切换，可切换广角，变焦和红外。
// 需要先调用控制权获取接口获取飞行器控制权后才能进行飞行器镜头切换。
// camera_index 和 lens_type 信息可通过获取组织下的设备列表接口下的 data.list.drone:.camera_list 中获取。
//
//	payLoad := `{
//					   "sn": "1581F6Q8D242100CPWEK", // 设备编码
//					   "camera_index": "81-0-0",  // 相机索引，0表示主相机，1表示辅助相机
//					   "lens_type": "zoom"   // 镜头类型，支持 "wide"（广角）、"zoom"（变焦）、"ir"（红外）
//					}`
func (F *FH2Adapter) UpdateDeviceChangeLens(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/device/change-lens", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// GetDeviceControl 获取设备控制权
// 如果你需要针对特定飞行器设备的负载进行操作，需要先获取该负载设备的控制权，然后对这个负载进行操作。比如对飞行器的某个负载相机的镜头进行切换。
// payload_index 参数可通过获取设备列表中data.list.drone.camera_list.camera_index传入。
//
//	payLoad := `{
//				   "sn": "1581F6Q8D242100CPWEK", // 设备编码
//				   "payload_index": ["81-0-0"]  // 负载索引，此参数可通过获取设备列表中data.list.drone.camera_list.camera_index传入
//				}`
func (F *FH2Adapter) GetDeviceControl(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/device/control", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// DeleteDeviceControl 设备控制权释放
// 释放指定设备的指定负载的控制权。
// payload_index 参数可通过获取设备列表中data.list.drone.camera_list.camera_index传入。
//
//	payLoad := `{
//				   "sn": "1581F6Q8D242100CPWEK", // 设备编码
//				   "payload_index": ["81-0-0"]  // 负载索引，此参数可通过获取设备列表中data.list.drone.camera_list.camera_index传入
//				}`
func (F *FH2Adapter) DeleteDeviceControl(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/device/control", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodDelete, url, payLoad)
	return string(resp), err
}

// UpdateDeviceStreamQuality 图传清晰度设置(支持自动，流畅，超清设置)
// 需要先调用控制权获取接口获取飞行器控制权后才能进行图传清晰度设置
// camera_index 信息可通过获取组织下的设备列表接口下的 data.list.drone:.camera_list 中获取
//
//	payLoad := `{
//				   "sn": "1581F6Q8D242100CPWEK", // 设备编码
//				   "camera_index": 81-0-0,  // 相机索引，0表示主相机，1表示辅助相机
//				   "quality": "ultra_high_definition"   // 清晰度选项，支持 "adaptive"（自动）、"smooth"（流畅）、"ultra_high_definition"（超清）
//				}`
func (F *FH2Adapter) UpdateDeviceStreamQuality(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/device/stream/quality", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodPut, url, payLoad)
	return string(resp), err
}

// CreateDeviceRTK 自定义网络RTK标定
//
//	payLoad := `{
//			   "host": "10.53.226.97",  // 服务器地址
//			   "port": 8002,   // 端口
//			   "account": "vntcqkj1CfCO", // 账号
//			   "password": "Sbii1qoJBows", // 密码
//			   "mount_point": "RTCM33_GRCEJ" // 挂载点
//			}`
func (F *FH2Adapter) CreateDeviceRTK(ctx context.Context, deviceSn string, payLoad io.Reader) (string, error) {
	if err := F.validator.ValidateDeviceSN(deviceSn); err != nil {
		return "", fmt.Errorf("设备序列号验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/device/%s/rtk", config.FH2Settings["host"], deviceSn)
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// LiveStreamStart 开启直播
// 启用设备直播功能后，设备将自动向媒体服务器推流，并返回对应直播供应商的鉴权信息。需集成供应商的SDK（如火山引擎、声网、SRS等），通过鉴权信息调用其接口拉取直播流。
// camera_index 信息可通过获取组织下的设备列表接口下的 data.list.drone.camera_list 中获取。
// 没有拉流观众5分钟后，将停止直播推流。
//
//	payLoad := `{
//				  "sn": "1581F6Q8D242100CPWEK", // 设备编码
//				  "camera_index": "81-0-0",  // 相机索引，0表示主相机，1表示辅助相机
//				  "video_expire": 7200   // 直播推流Token有效期，超过这个有限期直播将中止。
//				  "quality_type": "adaptive"   // 直播清晰度，adaptive（自动）、smooth（流畅）、ultra_high_definition（超清）
//				}`
func (F *FH2Adapter) LiveStreamStart(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/live-stream/start", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// CreateFlightTask 创建飞行任务
//
//	payLoad := `{
//			   "name": "测试任务",
//			   "wayline_uuid": "6d88fbe5-a399-485a-86ba-7bbdbb99edec",
//			   "sn": "7CTXN4A00B096H",
//			   "rth_altitude": 80,
//			   "rth_mode": "optimal",
//			   "wayline_precision_type": "gps",
//			   "out_of_control_action_in_flight": "return_home",
//			   "resumable_status": "auto",
//			   "task_type": "immediate",
//			   "time_zone": "Asia/Shanghai",
//			   "repeat_type": "nonrepeating",
//			   "min_battery_capacity": 60
//			}`
func (F *FH2Adapter) CreateFlightTask(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// UpdateFlightTaskStatus 更新飞行状态
//
//	payLoad := `{
//				   "status": restored/restored  //任务挂起&任务恢复
//				}`
func (F *FH2Adapter) UpdateFlightTaskStatus(ctx context.Context, taskUUID string, payLoad io.Reader) (string, error) {
	if err := F.validator.ValidateUUID(taskUUID); err != nil {
		return "", fmt.Errorf("UUID验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/%s/status", config.FH2Settings["host"], taskUUID)
	resp, err := F.doRequestWithTenant(ctx, http.MethodPut, url, payLoad)
	return string(resp), err
}

// GetFlightTaskInfo 获取飞行任务信息
func (F *FH2Adapter) GetFlightTaskInfo(ctx context.Context, taskUUID string) (string, error) {
	if err := F.validator.ValidateUUID(taskUUID); err != nil {
		return "", fmt.Errorf("UUID验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/%s", config.FH2Settings["host"], taskUUID)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetFlightTask 获取飞行任务列表
func (F *FH2Adapter) GetFlightTask(ctx context.Context, sn string, name string, beginAt int, endAt int, taskType string, status string) (string, error) {
	if err := F.validator.ValidateDeviceSN(sn); err != nil {
		return "", fmt.Errorf("设备序列号验证失败: %w", err)
	}
	if err := F.validator.ValidateProjectName(name); err != nil {
		return "", fmt.Errorf("任务名称验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/list?sn=%s&name=%s&begin_at=%d&end_at=%d&task_type=%s&status=%s", config.FH2Settings["host"], sn, name, beginAt, endAt, taskType, status)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetFlightTaskMedia 获取飞行任务产生的媒体资源
func (F *FH2Adapter) GetFlightTaskMedia(ctx context.Context, taskUUID string) (string, error) {
	if err := F.validator.ValidateUUID(taskUUID); err != nil {
		return "", fmt.Errorf("UUID验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/%s/media", config.FH2Settings["host"], taskUUID)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetFlightTaskTrack 获取飞行任务轨迹信息
func (F *FH2Adapter) GetFlightTaskTrack(ctx context.Context, taskUUID string) (string, error) {
	if err := F.validator.ValidateUUID(taskUUID); err != nil {
		return "", fmt.Errorf("UUID验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/%s/track", config.FH2Settings["host"], taskUUID)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// SetFinishUpload 航线上传完成通知
func (F *FH2Adapter) SetFinishUpload(ctx context.Context, objectKeyPrefix string, fileName string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/wayline/finish-upload", config.FH2Settings["host"])
	payload := strings.NewReader(fmt.Sprintf(`{"name":"%s","object_key":"%s"}`, fileName, objectKeyPrefix))
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payload)
	return string(resp), err
}

// GetWayLine 获取项目下航线列表
func (F *FH2Adapter) GetWayLine(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/wayline", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetWayLineInfo 获取项目下的航线详情
func (F *FH2Adapter) GetWayLineInfo(ctx context.Context, wayLineUuid string) (string, error) {
	if err := F.validator.ValidateUUID(wayLineUuid); err != nil {
		return "", fmt.Errorf("UUID验证失败: %w", err)
	}
	url := fmt.Sprintf("%s/openapi/v0.1/wayline/%s", config.FH2Settings["host"], wayLineUuid)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// CreateModel 模型重建
func (F *FH2Adapter) CreateModel(ctx context.Context, payLoad io.Reader) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/model/create", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodPost, url, payLoad)
	return string(resp), err
}

// GetModelInfo 获取模型详情
func (F *FH2Adapter) GetModelInfo(ctx context.Context, modelId int64) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/model/%s", config.FH2Settings["host"], modelId)
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// GetModelList 获取项目下模型列表
func (F *FH2Adapter) GetModelList(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s//openapi/v0.1/model", config.FH2Settings["host"])
	resp, err := F.doRequestWithTenant(ctx, http.MethodGet, url, nil)
	return string(resp), err
}

// 实例化 FH2Adapter 并注册到插件系统（自动注册）
func init() {
	// 注册 FH2 适配器插件
	plugin.RegisterPlugin(plugin.FH2Plugin, reflect.TypeOf((*service.FH2DroneAdapter)(nil)).Elem(), func() interface{} {
		return NewFH2Adapter()
	})
}
