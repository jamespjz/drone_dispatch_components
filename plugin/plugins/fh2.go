package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/jamespi/drone_dispatch/config"
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

// 司空2适配器
// 司空2openapi接口文档地址：https://apifox.com/apidoc/shared/6b4ca90b-233f-48ac-818c-d694acb0663a/api-221842037
type FH2Adapter struct {
	Client     *http.Client
	xUserToken string
	mu         sync.RWMutex // 并发线程安全
}

// NewFH2Adapter 创建一个新的FH2适配器
func NewFH2Adapter() *FH2Adapter {
	return &FH2Adapter{
		xUserToken: "",
		Client:     &http.Client{},
	}
}

// 添加设置token的方法
func (f *FH2Adapter) SetUserToken(token string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.xUserToken = token
}

// 定义与API响应对应的结构体
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (F *FH2Adapter) doRequest(ctx context.Context, method, url string, body io.Reader, projectUUID string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", F.xUserToken)
	if projectUUID != "" {
		req.Header.Add("X-Project-Uuid", projectUUID)
	}

	res, err := F.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	//return io.ReadAll(res.Body)
	// 读取响应体
	bodyBytes, err := io.ReadAll(res.Body)
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

// 获取组织下的项目列表
func (F *FH2Adapter) GetprojectList() (string, error) {
	encodedQ := url.QueryEscape(config.FH2Settings["q"])
	url := fmt.Sprintf("%s/openapi/v0.1/project?page=1&page_size=10&q=%s&prj_authorized_status=project-status-authorized&usage=simple&sort_column=created_at&sort_type=ASC", config.FH2Settings["host"], encodedQ)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, "")
	return string(resp), err
}

// 获取项目的存储上传凭证
func (F *FH2Adapter) GetProjectStsToken(projectUuid string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/project/sts-token", config.FH2Settings["host"])
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取项目下的设备列表
func (F *FH2Adapter) GetDeviceList(projectUuid string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/project/device", config.FH2Settings["host"])
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取设备物模型信息
func (F *FH2Adapter) GetStsToken(projectUuid string, deviceSn string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/device/%s/state", config.FH2Settings["host"], deviceSn)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取设备HMS信息
func (F *FH2Adapter) GetDeviceHms(projectUuid string, deviceSnList string) (string, error) {
	encodedQ := url.QueryEscape(deviceSnList)
	url := fmt.Sprintf("%s/openapi/v0.1/device/hms?device_sn_list=%s", config.FH2Settings["host"], encodedQ)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取设备控制权

// 设备控制权释放

// 设备飞行器镜头切换

// 机场相机切换

// 图传清晰度设置

// 实时控制指令下发

// 获取项目的存储上传凭证

// 航线上传
func (F *FH2Adapter) SetFinishUpload(projectUuid string, objectKeyPrefix string, fileName string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/wayline/finish-upload", config.FH2Settings["host"])
	payload := strings.NewReader(fmt.Sprintf(`{"name":"%s","object_key":"%s"}`, fileName, objectKeyPrefix))
	resp, err := F.doRequest(context.Background(), http.MethodPost, url, payload, projectUuid)
	return string(resp), err
}

// 创建飞行任务
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
func (F *FH2Adapter) CreateFlightTask(projectUuid string, payLoad string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task", config.FH2Settings["host"])
	resp, err := F.doRequest(context.Background(), http.MethodPost, url, strings.NewReader(payLoad), projectUuid)
	return string(resp), err
}

// 更新飞行状态

// 获取飞行任务轨迹信息
func (F *FH2Adapter) GetFlightTaskTrack(projectUuid string, task_uuid string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/%s/track", config.FH2Settings["host"], task_uuid)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取飞行任务信息
func (F *FH2Adapter) GetFlightTaskInfo(projectUuid string, task_uuid string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/%s", config.FH2Settings["host"], task_uuid)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取飞行任务列表
func (F *FH2Adapter) GetFlightTask(projectUuid string, sn string, name string, begin_at int, end_at int, task_type string, status string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/list?sn=%s&name=%s&begin_at=%d&end_at=%d&task_type=%s&status=%s", config.FH2Settings["host"], sn, name, begin_at, end_at, task_type, status)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取飞行任务产生的媒体资源
func (F *FH2Adapter) GetFlightTaskMedia(projectUuid string, task_uuid string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/flight-task/%s/media", config.FH2Settings["host"], task_uuid)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 开启直播
func (F *FH2Adapter) LiveStreamStart(projectUuid string, payLoad string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/live-stream/start", config.FH2Settings["host"])
	resp, err := F.doRequest(context.Background(), http.MethodPost, url, strings.NewReader(payLoad), projectUuid)
	return string(resp), err
}

// 获取项目下航线列表
func (F *FH2Adapter) GetWayLine(projectUuid string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/wayline", config.FH2Settings["host"])
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 获取项目下的航线详情
func (F *FH2Adapter) GetWayLineInfo(projectUuid string, wayLineUuid string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/wayline/%s", config.FH2Settings["host"], wayLineUuid)
	resp, err := F.doRequest(context.Background(), http.MethodGet, url, nil, projectUuid)
	return string(resp), err
}

// 模型重建
func (F *FH2Adapter) CreateModel(projectUuid string, payLoad string) (string, error) {
	url := fmt.Sprintf("%s/openapi/v0.1/model/create", config.FH2Settings["host"])
	resp, err := F.doRequest(context.Background(), http.MethodPost, url, strings.NewReader(payLoad), projectUuid)
	return string(resp), err
}

// 实例化 FH2Adapter 并注册到插件系统（自动注册）
func init() {
	FH2Adapter := NewFH2Adapter()
	// 注册 FH2 适配器插件
	plugin.RegisterPlugin(plugin.FH2Plugin, reflect.TypeOf((*service.FH2DroneAdapter)(nil)).Elem(), FH2Adapter)
}
