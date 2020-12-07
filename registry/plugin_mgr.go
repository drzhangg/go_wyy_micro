package registry

import (
	"context"
	"fmt"
	"sync"
)

var (
	pluginMgr = &PluginMgr{
		plugins: make(map[string]Registry),
	}
)

type PluginMgr struct {
	plugins map[string]Registry
	lock    sync.Mutex
}

func (p *PluginMgr) registerPlugin(plugin Registry) (err error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	//判断是否注册过
	_, ok := p.plugins[plugin.Name()]
	if ok {
		err = fmt.Errorf("duplicate registry plugin") //重复注册插件
		return
	}

	p.plugins[plugin.Name()] = plugin
	return
}

// 初始化插件
func (p *PluginMgr) initRegistry(ctx context.Context, name string, opts ...Option) (registry Registry, err error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	plugin, ok := p.plugins[name]
	if !ok {
		err = fmt.Errorf("plugin %s not exists", name)
		return
	}

	registry = plugin
	err = plugin.Init(ctx, opts...)
	return
}

func RegisterPlugin(registry Registry) (err error) {
	return pluginMgr.registerPlugin(registry)
}

func InitRegistry(ctx context.Context, name string, opts ...Option) (registry Registry, err error) {
	return pluginMgr.initRegistry(ctx, name, opts...)
}
