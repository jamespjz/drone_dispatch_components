// 插件管理中心

package plugin

import (
	"gitee.com/jamespi/drone_dispatch/service"
	"reflect"
	"sync"
)

type PluginType string

const (
	DJIDock2Plugin PluginType = "dji_dock2" // DJI Dock2插件
	FH2Plugin      PluginType = "fh2"       // 司空2插件
	DJIPilotPlugin PluginType = "dji_pilot" // DJI Pilot插件
)

type Registry struct {
	// DroneAdapters is a map of drone types to their respective adapters.
	mu            sync.RWMutex // 并发安全
	PluginFactory map[PluginType]map[reflect.Type]func() interface{}
	Status        map[PluginType]service.PluginStatus // 插件状态
}

// registry 全局单例模式- 使用工厂模式
var registry = &Registry{
	PluginFactory: make(map[PluginType]map[reflect.Type]func() interface{}),
	Status:        make(map[PluginType]service.PluginStatus),
}

// RegisterPlugin 注册插件.
func RegisterPlugin(pluginType PluginType, ifaceType reflect.Type, Build func() interface{}) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if registry.PluginFactory[pluginType] == nil {
		registry.PluginFactory[pluginType] = make(map[reflect.Type]func() interface{})
	}

	registry.PluginFactory[pluginType][ifaceType] = Build
	registry.Status[pluginType] = service.PluginRegistered
}

// Enable 启用插件
func Enable(pluginType PluginType) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	if _, ok := registry.PluginFactory[pluginType]; ok {
		registry.Status[pluginType] = service.PluginEnabled
	}
}

// Disable 禁用插件
func Disable(pluginType PluginType) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, ok := registry.PluginFactory[pluginType]; ok {
		registry.Status[pluginType] = service.PluginDisabled
	}
}

// Unload 卸载插件
func Unload(pluginType PluginType) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	if _, ok := registry.PluginFactory[pluginType]; ok {
		registry.Status[pluginType] = service.PluginUnloaded
	}
	delete(registry.PluginFactory, pluginType)
}

// Get 获取启用状态下的适配器
// 判断该插件是否存在并且已返回目标T接口类型的插件实例
func Get[T interface{}](pluginType PluginType) (T, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	// 泛型参数T的类型实例化
	var plugin T
	ifaceType := reflect.TypeOf(&plugin).Elem()
	if factory, ok := registry.PluginFactory[pluginType][ifaceType]; ok {
		impl := factory()
		return impl.(T), true
	}

	return plugin, false
}
