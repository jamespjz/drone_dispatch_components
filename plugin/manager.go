// 插件生命周期管理器
package plugin

import (
	"gitee.com/jamespi/lecheng-drone/service"
	"reflect"
)

type PluginInfo struct {
	PluginType PluginType
	Interfaces []reflect.Type
	Status     service.PluginStatus
}

// 批量启用某些插件
func LoadEnableList(selected []string) {
	for _, pluginName := range selected {
		Enable(PluginType(pluginName))
	}
}

// 批量禁用某些插件
func LoadDisabledList(selected []string) {
	for _, pluginName := range selected {
		Disable(PluginType(pluginName))
	}
}

// 批量卸载某些插件
func LoadUnloadedList(selected []string) {
	for _, pluginName := range selected {
		Unload(PluginType(pluginName))
	}
}

// 获取全部注册的插件列表
func PluginsList() []PluginInfo {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	var list []PluginInfo
	for pluginType, ifaceMap := range registry.Plugins {
		var ifaceTypes []reflect.Type
		for ifaceType := range ifaceMap {
			ifaceTypes = append(ifaceTypes, ifaceType)
		}

		list = append(list, PluginInfo{
			PluginType: pluginType,
			Interfaces: ifaceTypes,
			Status:     registry.Status[pluginType],
		})
	}
	return list
}
