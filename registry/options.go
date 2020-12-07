package registry

import "time"

type Options struct {
	Addrs        []string      //etcd地址
	Timeout      time.Duration //超时时间
	RegistryPath string        //注册地址
	HeartBeat    int64         //心跳时间
}

// 选项模式
type Option func(opts *Options)

func WithAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func WithRegistryPath(registryPath string) Option {
	return func(opts *Options) {
		opts.RegistryPath = registryPath
	}
}

func WithHeartBeat(hearbeat int64) Option {
	return func(opts *Options) {
		opts.HeartBeat = hearbeat
	}
}
