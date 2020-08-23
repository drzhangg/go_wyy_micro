package etcd

import "go_wyy_micro/registry"

const (
	MaxServiceNum = 10
)

// etcd注册插件
type EtcdRegistry struct {
	options *registry.Options
}