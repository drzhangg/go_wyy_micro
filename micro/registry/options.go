package registry

import "time"

type Options struct {
	Addrs        []string
	Timeout      time.Duration
	RegistryPath string
	HeartBeat    int64
}

type Option func(options *Options)

func WithAddr(addrs []string) Option {
	return func(options *Options) {
		options.Addrs = addrs
	}
}

func WithRegistryPath(path string) Option {
	return func(options *Options) {
		options.RegistryPath = path
	}
}

func WithHeartBeat(heartBeat int64) Option {
	return func(options *Options) {
		options.HeartBeat = heartBeat
	}
}
