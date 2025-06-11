// 插件生命周期管理器
package plugin

import "gitee.com/jamespi/lecheng-drone/service"

// 批量启用某些插件
func LoadEnableList(selected []string) {
	for _, pluginName := range selected {
		Enable(pluginName)
	}
}

// 批量禁用某些插件
func LoadDisabledList(selected []string) {
	for _, pluginName := range selected {
		Disable(pluginName)
	}
}

// 批量卸载某些插件
func LoadUnloadedList(selected []string) {
	for _, pluginName := range selected {
		Unload(pluginName)
	}
}

// 获取全部注册的插件列表
func PluginsList() []service.PluginMeta {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	var list []service.PluginMeta
	for _, meta := range registry.Plugins {
		list = append(list, *meta)
	}
	return list
}
