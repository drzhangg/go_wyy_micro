package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go_wyy_micro/registry"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxServiceNum          = 0
	MaxSyncServiceInterval = time.Second * 10
)

//etcd注册信息
type EtcdRegistry struct {
	options     *registry.Options
	client      *clientv3.Client       //etcd客户端
	serviceChan chan *registry.Service //节点信息

	value              atomic.Value
	lock               sync.Mutex
	registryServiceMap map[string]*RegistryService
}

//map获取节点信息
type AllServiceInfo struct {
	serviceMap map[string]*registry.Service
}

//etcd租约信息
type RegistryService struct {
	id          clientv3.LeaseID                        //租约id
	service     *registry.Service                       //
	registered  bool                                    //是否注册
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse //续租chan
}

var (
	etcdRegistry *EtcdRegistry = &EtcdRegistry{
		serviceChan:        make(chan *registry.Service, MaxServiceNum),
		registryServiceMap: make(map[string]*RegistryService, MaxServiceNum),
	}
)

func init() {
	//allServiceInfo := &AllServiceInfo{
	//	serviceMap: make(map[string]*registry.Service, MaxServiceNum),
	//}
	//
	//etcdRegistry.value.Store(allServiceInfo)
	registry.RegisterPlugin(etcdRegistry)
	go etcdRegistry.run()
}

/*
//插件的名字
	Name() string
	//初始化插件
	Init(ctx context.Context, opt ...Options)
	//服务注册
	Register(ctx context.Context, service *Service) (err error)
	//服务反注册
	Unregister(ctx context.Context, service *Service) (err error)
	//服务发现：通过服务的名字获取服务的位置信息（ip和port列表）
	GetService(ctx context.Context, name string) (service *Service, err error)
*/
func (e *EtcdRegistry) Name() string {
	return "etcd"
}

func (e *EtcdRegistry) Init(ctx context.Context, opts ...registry.Option) (err error) {
	e.options = &registry.Options{}
	for _, opt := range opts {
		opt(e.options)
	}

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.options.Addrs,
		DialTimeout: e.options.Timeout,
	})

	if err != nil {
		err = fmt.Errorf("init etcd failed, err:%v", err)
		return
	}
	return
}

func (e *EtcdRegistry) Register(ctx context.Context, service *registry.Service) (err error) {
	select {
	case e.serviceChan <- service:
	default:
		err = fmt.Errorf("register chan is full")
		return
	}
	return
}

func (e *EtcdRegistry) Unregister(ctx context.Context, service *registry.Service) (err error) {
	return
}

//获取当前需要注册的服务
func (e *EtcdRegistry) run() {
	//ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		case service := <-e.serviceChan:
			//判断服务是否存在，存在说明已经注册过
			registryService, ok := e.registryServiceMap[service.Name]
			if ok {
				for _, node := range service.Nodes {
					registryService.service.Nodes = append(registryService.service.Nodes, node)
				}
				registryService.registered = false
				break
			}
			//如果不存在，将服务放到map里
			registryService = &RegistryService{
				service: service,
			}
			e.registryServiceMap[service.Name] = registryService
		//case <-ticker.C:
		//	e.syncServiceFromEtcd()
		default:
			e.registerOrKeepAlive()
			time.Sleep(time.Millisecond * 5000)
		}
	}
}

//注册或者续约
func (e *EtcdRegistry) registerOrKeepAlive() {
	for _, registryService := range e.registryServiceMap {
		if registryService.registered {
			e.keepAlive(registryService)
			continue
		}
		e.registerService(registryService)
	}
}

//保持心跳
func (e *EtcdRegistry) keepAlive(registryService *RegistryService) {
	select {
	case resp := <-registryService.keepAliveCh:
		if resp == nil {
			registryService.registered = false
			return
		}
	}
	return
}

//注册服务 在etcd中进行put
func (e *EtcdRegistry) registerService(registryService *RegistryService) {
	//续租时间
	resp, err := e.client.Grant(context.TODO(), e.options.HeartBeat)
	if err != nil {
		return
	}

	registryService.id = resp.ID

	for _, node := range registryService.service.Nodes {
		tmp := &registry.Service{
			Name:  registryService.service.Name,
			Nodes: []*registry.Node{node},
		}

		data, err := json.Marshal(tmp)
		if err != nil {
			continue
		}

		key := e.serviceNodePath(tmp)
		fmt.Printf("register key:%s\n", key)
		_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
		if err != nil {
			continue
		}

		ch, err := e.client.KeepAlive(context.TODO(), resp.ID)
		if err != nil {
			continue
		}
		registryService.keepAliveCh = ch
		registryService.registered = true
	}
}

func (e *EtcdRegistry) serviceNodePath(service *registry.Service) string {
	nodeIP := fmt.Sprintf("%s:%d", service.Nodes[0].IP, service.Nodes[0].Port)
	return path.Join(e.options.RegistryPath, service.Name, nodeIP)
}

func (e *EtcdRegistry) servicePath(name string) string {
	return path.Join(e.options.RegistryPath, name)
}

func (e *EtcdRegistry) getServiceFromCache(ctx context.Context, name string) (service *registry.Service, ok bool) {
	allServiceInfo := e.value.Load().(*AllServiceInfo)
	//一般情况下，都会从缓存中读取
	service, ok = allServiceInfo.serviceMap[name]
	return
}

func (e *EtcdRegistry) syncServiceFromEtcd() {
	var allserviceInfoNew = &AllServiceInfo{serviceMap: make(map[string]*registry.Service, MaxServiceNum)}

	ctx := context.TODO()
	allServiceInfo := e.value.Load().(*AllServiceInfo)

	//对于缓存的每一个服务，都需要从etcd中进行更新
	for _, service := range allServiceInfo.serviceMap {
		key := e.servicePath(service.Name)
		resp, err := e.client.Get(ctx, key, clientv3.WithPrefix())
		if err != nil {
			allserviceInfoNew.serviceMap[service.Name] = service
			continue
		}
		serviceNew := &registry.Service{
			Name: service.Name,
		}

		for _, kv := range resp.Kvs {
			value := kv.Value
			var tmpService registry.Service
			err = json.Unmarshal(value, &tmpService)
			if err != nil {
				fmt.Printf("unmarshal failed, err:%v value:%s", err, string(value))
				return
			}

			for _, node := range tmpService.Nodes {
				serviceNew.Nodes = append(serviceNew.Nodes, node)
			}
		}
		allserviceInfoNew.serviceMap[serviceNew.Name] = serviceNew
	}
	e.value.Store(allserviceInfoNew)
}

//func (e *EtcdRegistry) GetService(ctx context.Context, name string) (service *registry.Service, err error) {
//	//一般情况下，都会从缓存中读取
//
//}
