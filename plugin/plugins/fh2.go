package plugins

import (
	"fmt"
	"gitee.com/jamespi/lecheng-drone/config"
	"gitee.com/jamespi/lecheng-drone/plugin"
	"gitee.com/jamespi/lecheng-drone/service"
	"github.com/google/uuid"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// 司空2适配器
// 司空2openapi接口文档地址：https://apifox.com/apidoc/shared/6b4ca90b-233f-48ac-818c-d694acb0663a/api-221842037
type FH2Adapter struct {
	Client *http.Client
}

// NewFH2Adapter 创建一个新的FH2适配器
func NewFH2Adapter() *FH2Adapter {
	return &FH2Adapter{
		Client: &http.Client{},
	}
}

// 获取组织下的项目列表
func (F *FH2Adapter) GetprojectList() (string, error) {
	url := config.FH2Settings["host"] + "/openapi/v0.1/project?page=1&page_size=10&q=" + config.FH2Settings["q"] + "&prj_authorized_status=project-status-authorized&usage=simple&sort_column=created_at&sort_type=ASC"
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", config.FH2Settings["user_token"])

	res, err := F.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}

// 获取项目下的设备列表
func (F *FH2Adapter) GetDeviceList(projectUuid string) (string, error) {
	url := config.FH2Settings["host"] + "/openapi/v0.1/project/device"
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", config.FH2Settings["user_token"])
	req.Header.Add("X-Project-Uuid", projectUuid)

	res, err := F.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil
}

// 获取设备物模型信息
func (F *FH2Adapter) GetStsToken(projectUuid string) (string, error) {
	url := config.FH2Settings["host"] + "/openapi/v0.1/sts/token"
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", config.FH2Settings["user_token"])
	req.Header.Add("X-Project-Uuid", projectUuid)

	res, err := F.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil

}

// 获取设备HMS信息

// 获取设备控制权

// 设备控制权释放

// 设备飞行器镜头切换

// 机场相机切换

// 图传清晰度设置

// 实时控制指令下发

// 获取项目的存储上传凭证

// 航线上传
func (F *FH2Adapter) SetFinishUpload(projectUuid string, objectKeyPrefix string, fileName string) (string, error) {
	url := config.FH2Settings["host"] + "/openapi/v0.1/wayline/finish-upload"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{
    "name": "%s",
    "object_key": "%s"
}`, fileName, objectKeyPrefix))
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", config.FH2Settings["user_token"])
	req.Header.Add("X-Project-Uuid", projectUuid)

	res, err := F.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil

}

// 创建飞行任务
func (F *FH2Adapter) CreateFlightTask(projectUuid string, payLoad string) (string, error) {
	url := config.FH2Settings["host"] + "/openapi/v0.1/flight-task"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`"%s"`, payLoad))
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", config.FH2Settings["user_token"])
	req.Header.Add("X-Project-Uuid", projectUuid)

	res, err := F.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil

}

// 更新飞行状态

// 获取飞行任务轨迹信息

// 获取飞行任务产生的媒体资源

// 获取飞行任务信息

// 获取飞行任务列表

// 获取项目下航线列表
func (F *FH2Adapter) GetWayLine(projectUuid string) (string, error) {
	url := config.FH2Settings["host"] + "/openapi/v0.1/wayline"
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", config.FH2Settings["user_token"])
	req.Header.Add("X-Project-Uuid", projectUuid)

	res, err := F.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil

}

// 获取项目下的航线详情
func (F *FH2Adapter) GetWayLineInfo(projectUuid string, wayLineUuid string) (string, error) {
	url := config.FH2Settings["host"] + "/openapi/v0.1/wayline/" + wayLineUuid
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}
	// 设置请求头
	req.Header.Add("X-Request-Id", uuid.New().String())
	req.Header.Add("X-Language", "zh")
	req.Header.Add("X-User-Token", config.FH2Settings["user_token"])
	req.Header.Add("X-Project-Uuid", projectUuid)

	res, err := F.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	return string(body), nil

}

// 实例化 FH2Adapter 并注册到插件系统（自动注册）
func init() {
	FH2Adapter := NewFH2Adapter()
	// 注册 FH2 适配器插件
	plugin.RegisterPlugin(plugin.FH2Plugin, reflect.TypeOf((*service.FH2DroneAdapter)(nil)).Elem(), FH2Adapter)
}
