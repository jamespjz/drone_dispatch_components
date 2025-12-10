package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// Config 结构体保持不变
type Config struct {
	Drone          *Drone      `mapstructure:"Drone"`
	Mqtt           *MqttConfig `mapstructure:"Mqtt"`
	Database       *Database   `mapstructure:"Database"`
	RtmpURL        string      `mapstructure:"RtmpURL"`
	AmapKey        string      `mapstructure:"AmapKey"`
	TokenExpiresIn int         `mapstructure:"TokenExpiresIn"`
	FH2            *FH2        `mapstructure:"FH"`
}

type Drone struct {
	Dji *UAV `mapstructure:"Dji"`
	MMC *UAV `mapstructure:"MMC"`
	XAG *UAV `mapstructure:"XAG"`
}

type UAV struct {
	AppId        string `mapstructure:"appId"`
	AppKey       string `mapstructure:"appKey"`
	AppLicense   string `mapstructure:"appLicense"`
	Url          string `mapstructure:"url"`
	GatewaySn    string `mapstructure:"GatewaySn"`
	DjiWebsocket string `mapstructure:"DjiWebsocket"`
	DockSn       string `mapstructure:"DockSn"`
	ClientId     string `mapstructure:"ClientId"`
}

type MqttConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Database struct {
	Def   *Def `mapstructure:"default"`
	CK    *Def `mapstructure:"clickhouse"`
	Redis *Def `mapstructure:"redis"`
}

type Def struct {
	Link string `mapstructure:"link"`
}

type FH2 struct {
	Host       string `mapstructure:"host"`
	Q          string `mapstructure:"q"`
	XUserToken string `mapstructure:"xUserToken"`
}

// 全局变量和同步控制
var (
	cfg           Config
	configOnce    sync.Once
	configErr     error
	isInitialized bool
	autoInitOnce  sync.Once
	autoInitErr   error

	// 全局设置变量
	DjiSettings            map[string]string
	MMCSettings            map[string]string
	XAGSettings            map[string]string
	MqttSettings           map[string]string
	DatabaseSettings       map[string]string
	CKSettings             map[string]string
	RedisSettings          map[string]string
	RtmpURLSettings        string
	AmapKeySettings        string
	TokenExpiresInSettings int
	FH2Settings            map[string]string
)

// 默认配置路径
const (
	DefaultConfigPath = "./config.yaml"
)

// InitConfig 初始化配置，支持自定义配置文件路径
func InitConfig(configPath string) error {
	configOnce.Do(func() {
		if configPath == "" {
			configPath = DefaultConfigPath
		}
		configErr = loadConfig(configPath)
		if configErr == nil {
			isInitialized = true
		}
	})
	return configErr
}

// InitDefaultConfig 使用默认路径初始化配置
func InitDefaultConfig() error {
	return InitConfig(DefaultConfigPath)
}

// IsConfigInitialized 检查配置是否已成功初始化
func IsConfigInitialized() bool {
	return isInitialized
}

// GetConfig 安全获取配置实例
func GetConfig() (*Config, error) {
	if !isInitialized {
		return nil, fmt.Errorf("配置未初始化，请先调用 InitConfig 或 InitDefaultConfig")
	}
	if configErr != nil {
		return nil, fmt.Errorf("配置初始化失败: %w", configErr)
	}
	return &cfg, nil
}

// GetSetting 提供统一的方法获取特定配置项
func GetSetting(section, key string) (string, error) {
	if !isInitialized {
		return "", fmt.Errorf("配置未初始化")
	}

	// 这里可以根据需要实现更精细的配置项获取逻辑
	switch section {
	case "Drone":
		if value, exists := DjiSettings[key]; exists {
			return value, nil
		}
		// 可以扩展其他section
	}

	return "", fmt.Errorf("配置项不存在: %s.%s", section, key)
}

// loadConfig 加载和解析配置文件
func loadConfig(configPath string) error {
	viper.SetConfigFile(configPath)

	// 支持多种配置格式
	if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
		viper.SetConfigType("yaml")
	} else if strings.HasSuffix(configPath, ".json") {
		viper.SetConfigType("json")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("配置文件读取失败: %w", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("配置解析失败: %w", err)
	}

	// 验证必要配置
	if err := validateConfig(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// 初始化全局设置
	initGlobalSettings()

	return nil
}

// validateConfig 验证配置的必要项
func validateConfig() error {
	var missingFields []string

	if cfg.Drone == nil {
		missingFields = append(missingFields, "Drone配置")
	} else {
		if cfg.Drone.Dji == nil {
			missingFields = append(missingFields, "Dji无人机配置")
		}
		// 可以添加其他必要字段的验证
	}

	if cfg.Mqtt == nil {
		missingFields = append(missingFields, "MQTT配置")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("缺少必要配置项: %s", strings.Join(missingFields, ", "))
	}

	return nil
}

// initGlobalSettings 初始化全局设置变量
func initGlobalSettings() {
	// 初始化无人机配置
	initDroneSettings()

	// 初始化MQTT配置
	if cfg.Mqtt != nil {
		MqttSettings = map[string]string{
			"host":     cfg.Mqtt.Host,
			"port":     cfg.Mqtt.Port,
			"username": cfg.Mqtt.Username,
			"password": cfg.Mqtt.Password,
		}
	} else {
		MqttSettings = make(map[string]string)
	}

	// 初始化数据库配置
	initDatabaseSettings()

	// 初始化其他配置
	if cfg.RtmpURL != "" {
		RtmpURLSettings = cfg.RtmpURL
	}

	if cfg.AmapKey != "" {
		AmapKeySettings = cfg.AmapKey
	}

	if cfg.TokenExpiresIn > 0 {
		TokenExpiresInSettings = cfg.TokenExpiresIn
	}

	// 初始化FH2配置
	if cfg.FH2 != nil {
		FH2Settings = map[string]string{
			"host":       cfg.FH2.Host,
			"q":          cfg.FH2.Q,
			"xUserToken": cfg.FH2.XUserToken,
		}
	} else {
		FH2Settings = make(map[string]string)
	}
}

// initDroneSettings 初始化无人机相关设置
func initDroneSettings() {
	if cfg.Drone != nil {
		if cfg.Drone.Dji != nil {
			DjiSettings = map[string]string{
				"appId":        cfg.Drone.Dji.AppId,
				"appKey":       cfg.Drone.Dji.AppKey,
				"appLicense":   cfg.Drone.Dji.AppLicense,
				"url":          cfg.Drone.Dji.Url,
				"GatewaySn":    cfg.Drone.Dji.GatewaySn,
				"DjiWebsocket": cfg.Drone.Dji.DjiWebsocket,
				"DockSn":       cfg.Drone.Dji.DockSn,
				"ClientId":     cfg.Drone.Dji.ClientId,
			}
		} else {
			DjiSettings = make(map[string]string)
		}

		if cfg.Drone.MMC != nil {
			MMCSettings = map[string]string{
				"appId":        cfg.Drone.MMC.AppId,
				"appKey":       cfg.Drone.MMC.AppKey,
				"appLicense":   cfg.Drone.MMC.AppLicense,
				"url":          cfg.Drone.MMC.Url,
				"GatewaySn":    cfg.Drone.MMC.GatewaySn,
				"DjiWebsocket": cfg.Drone.MMC.DjiWebsocket,
				"DockSn":       cfg.Drone.MMC.DockSn,
				"ClientId":     cfg.Drone.MMC.ClientId,
			}
		} else {
			MMCSettings = make(map[string]string)
		}

		if cfg.Drone.XAG != nil {
			XAGSettings = map[string]string{
				"appId":        cfg.Drone.XAG.AppId,
				"appKey":       cfg.Drone.XAG.AppKey,
				"appLicense":   cfg.Drone.XAG.AppLicense,
				"url":          cfg.Drone.XAG.Url,
				"GatewaySn":    cfg.Drone.XAG.GatewaySn,
				"DjiWebsocket": cfg.Drone.XAG.DjiWebsocket,
				"DockSn":       cfg.Drone.XAG.DockSn,
				"ClientId":     cfg.Drone.XAG.ClientId,
			}
		} else {
			XAGSettings = make(map[string]string)
		}
	} else {
		DjiSettings = make(map[string]string)
		MMCSettings = make(map[string]string)
		XAGSettings = make(map[string]string)
	}
}

// initDatabaseSettings 初始化数据库相关设置
func initDatabaseSettings() {
	if cfg.Database != nil {
		if cfg.Database.Def != nil {
			DatabaseSettings = map[string]string{
				"link": cfg.Database.Def.Link,
			}
		} else {
			DatabaseSettings = make(map[string]string)
		}

		if cfg.Database.CK != nil {
			CKSettings = map[string]string{
				"link": cfg.Database.CK.Link,
			}
		} else {
			CKSettings = make(map[string]string)
		}

		if cfg.Database.Redis != nil {
			RedisSettings = map[string]string{
				"link": cfg.Database.Redis.Link,
			}
		} else {
			RedisSettings = make(map[string]string)
		}
	} else {
		DatabaseSettings = make(map[string]string)
		CKSettings = make(map[string]string)
		RedisSettings = make(map[string]string)
	}
}
