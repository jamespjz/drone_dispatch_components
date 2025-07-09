package config

import (
	"fmt"
	"github.com/spf13/viper"
)

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
	GatewaySn    string `mapstructure:"GatewaySn"`    //网关序列号
	DjiWebsocket string `mapstructure:"DjiWebsocket"` //websocket地址
	DockSn       string `mapstructure:"DockSn"`       //机场序列号
	ClientId     string `mapstructure:"ClientId"`     //客户端ID
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

var cfg Config
var DjiSettings map[string]string
var MMCSettings map[string]string
var XAGSettings map[string]string
var MqttSettings map[string]string
var DatabaseSettings map[string]string
var CKSettings map[string]string
var RedisSettings map[string]string
var RtmpURLSettings string
var AmapKeySettings string
var TokenExpiresInSettings int
var FH2Settings map[string]string

func init() {
	viper.SetConfigName("config") // 文件名（不含扩展名）
	viper.SetConfigType("yaml")   // 文件类型
	viper.AddConfigPath(".")      // 搜索路径

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("配置文件读取失败: %w", err))
	}

	err = viper.Unmarshal(&cfg) // 自动绑定到结构体
	if err != nil {
		panic(fmt.Errorf("配置解析失败: %w", err))
	}

	fmt.Printf("应用端口: %s\n", cfg.Drone.Dji.AppId)

	if cfg.Drone != nil && cfg.Drone.Dji != nil {
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

	if cfg.Drone != nil && cfg.Drone.MMC != nil {
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

	if cfg.Drone != nil && cfg.Drone.XAG != nil {
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

	if cfg.Database != nil && cfg.Database.Def != nil {
		DatabaseSettings = map[string]string{
			"link": cfg.Database.Def.Link,
		}
	} else {
		DatabaseSettings = make(map[string]string)
	}

	if cfg.Database != nil && cfg.Database.CK != nil {
		CKSettings = map[string]string{
			"link": cfg.Database.CK.Link,
		}
	} else {
		CKSettings = make(map[string]string)
	}

	if cfg.Database != nil && cfg.Database.Redis != nil {
		RedisSettings = map[string]string{
			"link": cfg.Database.Redis.Link,
		}
	} else {
		RedisSettings = make(map[string]string)
	}

	if cfg.RtmpURL != "" {
		RtmpURLSettings = cfg.RtmpURL
	}

	if cfg.AmapKey != "" {
		AmapKeySettings = cfg.AmapKey
	}

	if cfg.TokenExpiresIn > 0 {
		TokenExpiresInSettings = cfg.TokenExpiresIn
	}
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
