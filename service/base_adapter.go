package service

type PluginStatus string

// 定义适配器插件生命周期状态
const (
	PluginRegistered PluginStatus = "registered" //插件已注册
	PluginUnloaded   PluginStatus = "unloaded"   //插件已卸载
	PluginDisabled   PluginStatus = "disabled"   //插件已禁用
	PluginEnabled    PluginStatus = "enabled"    //插件已启用
)

type BaseAdapter interface {
}
