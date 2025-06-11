// 插件管理中心

package plugin

import (
	"gitee.com/jamespi/lecheng-drone/service"
	"sync"
)

type Registry struct {
	// DroneAdapters is a map of drone types to their respective adapters.
	mu      sync.RWMutex // 并发安全
	Plugins map[string]*service.PluginMeta
}

// 实例化全局单例.
var registry = &Registry{
	Plugins: make(map[string]*service.PluginMeta),
}

// 注册插件.
func RegisterPlugin(pluginName string, Build func() service.DroneAdapter) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.Plugins[pluginName] = &service.PluginMeta{
		Name:   pluginName,
		Status: service.PluginRegistered,
		Build:  Build,
	}
}

// 启用插件
func Enable(pluginName string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if meta, ok := registry.Plugins[pluginName]; ok {
		meta.Status = service.PluginEnabled
	}
}

// 禁用插件
func Disable(pluginName string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if meta, ok := registry.Plugins[pluginName]; ok {
		meta.Status = service.PluginDisabled
	}
}

// 卸载插件
func Unload(pluginName string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if meta, ok := registry.Plugins[pluginName]; ok {
		meta.Status = service.PluginUnloaded
	}
	delete(registry.Plugins, pluginName)
}

// 获取启用状态下的适配器
func Get(pluginName string) service.DroneAdapter {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	if meta, ok := registry.Plugins[pluginName]; ok && meta.Status == service.PluginEnabled {
		return meta.Build()
	}
	return nil
}
