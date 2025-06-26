package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Drone          *Drone      `yarml:"drone"`
	Mqtt           *MqttConfig `yarml:"Mqtt"`
	Database       *Database   `yarml:"Database"`
	RtmpURL        string      `yarml:"RtmpURL"`
	AmapKey        string      `yarml:"AmapKey"`
	TokenExpiresIn int         `yarml:"TokenExpiresIn"`
	FH2            *FH2        `yarml:"FH"`
}

type Drone struct {
	Dji *UAV `yarml:"Dji"`
	MMC *UAV `yarml:"MMC"`
	XAG *UAV `yarml:"XAG"`
}

type UAV struct {
	appId        string `yarml:"appId"`
	appKey       string `yarml:"appKey"`
	appLicense   string `yarml:"appLicense"`
	url          string `yarml:"url"`
	GatewaySn    string `yarml:"GatewaySn"`    //网关序列号
	DjiWebsocket string `yarml:"DjiWebsocket"` //websocket地址
	DockSn       string `yarml:"DockSn"`       //机场序列号
	ClientId     string `yarml:"ClientId"`     //客户端ID
}

type MqttConfig struct {
	host     string `yarml:"host"`
	port     string `yarml:"port"`
	username string `yarml:"username"`
	password string `yarml:"password"`
}

type Database struct {
	Def   *Def `yarml:"default"`
	CK    *Def `yarml:"clickhouse"`
	Redis *Def `yarml:"redis"`
}

type Def struct {
	Link string `yarml:"link"`
}

type FH2 struct {
	host       string `yarml:"host"`
	port       string `yarml:"port"`
	q          string `yarml:"q"`
	XUserToken string `yarml:"xUserToken"`
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

	fmt.Printf("应用端口: %d\n", cfg.Drone.Dji.appId)

	DjiSettings = map[string]string{
		"appId":        cfg.Drone.Dji.appId,
		"appKey":       cfg.Drone.Dji.appKey,
		"appLicense":   cfg.Drone.Dji.appLicense,
		"url":          cfg.Drone.Dji.url,
		"GatewaySn":    cfg.Drone.Dji.GatewaySn,
		"DjiWebsocket": cfg.Drone.Dji.DjiWebsocket,
		"DockSn":       cfg.Drone.Dji.DockSn,
		"ClientId":     cfg.Drone.Dji.ClientId,
	}

	MMCSettings = map[string]string{
		"appId":        cfg.Drone.MMC.appId,
		"appKey":       cfg.Drone.MMC.appKey,
		"appLicense":   cfg.Drone.MMC.appLicense,
		"url":          cfg.Drone.MMC.url,
		"GatewaySn":    cfg.Drone.MMC.GatewaySn,
		"DjiWebsocket": cfg.Drone.MMC.DjiWebsocket,
		"DockSn":       cfg.Drone.MMC.DockSn,
		"ClientId":     cfg.Drone.MMC.ClientId,
	}

	XAGSettings = map[string]string{
		"appId":        cfg.Drone.XAG.appId,
		"appKey":       cfg.Drone.XAG.appKey,
		"appLicense":   cfg.Drone.XAG.appLicense,
		"url":          cfg.Drone.XAG.url,
		"GatewaySn":    cfg.Drone.XAG.GatewaySn,
		"DjiWebsocket": cfg.Drone.XAG.DjiWebsocket,
		"DockSn":       cfg.Drone.XAG.DockSn,
		"ClientId":     cfg.Drone.XAG.ClientId,
	}

	MqttSettings = map[string]string{
		"host":     cfg.Mqtt.host,
		"port":     cfg.Mqtt.port,
		"username": cfg.Mqtt.username,
		"password": cfg.Mqtt.password,
	}

	DatabaseSettings = map[string]string{
		"link": cfg.Database.Def.Link,
	}

	CKSettings = map[string]string{
		"link": cfg.Database.CK.Link,
	}

	RedisSettings = map[string]string{
		"link": cfg.Database.Redis.Link,
	}

	RtmpURLSettings = cfg.RtmpURL

	AmapKeySettings = cfg.AmapKey

	TokenExpiresInSettings = cfg.TokenExpiresIn

	FH2Settings = map[string]string{
		"host":     cfg.FH2.host,
		"port":     cfg.FH2.port,
		"q":        cfg.FH2.q,
		"username": cfg.FH2.XUserToken,
	}
}
