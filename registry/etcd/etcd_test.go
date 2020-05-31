package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go_wyy_micro/registry"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	registryInst, err := registry.InitRegistry(context.TODO(), "etcd",
		registry.WithAddrs([]string{"127.0.0.1:2379"}),
		registry.WithRegistryPath("/ibinarytree/koala/"),
		registry.WithTimeout(time.Second),
		registry.WithHeartBeat(5))
	if err != nil {
		t.Errorf("init registry failed, err:%v", err)
		return
	}

	service := &registry.Service{
		Name: "comment_service",
	}

	service.Nodes = append(service.Nodes, &registry.Node{
		IP:   "127.0.0.1",
		Port: 8801,
	}, &registry.Node{
		IP:   "127.0.0.2",
		Port: 8801,
	})
	registryInst.Register(context.TODO(), service)

	//go func() {
	//	time.Sleep(time.Second * 10)
	//	service.Nodes = append(service.Nodes, &registry.Node{
	//		IP:   "127.0.0.3",
	//		Port: 8801,
	//	},
	//		&registry.Node{
	//			IP:   "127.0.0.4",
	//			Port: 8801,
	//		},
	//	)
	//	registryInst.Register(context.TODO(), service)
	//}()
	for {
		//service, err := registryInst.GetService(context.TODO(), "comment_service")
		//if err != nil {
		//	t.Errorf("get service failed, err:%v", err)
		//	return
		//}
		//
		//for _, node := range service.Nodes {
		//	fmt.Printf("service:%s, node:%#v\n", service.Name, node)
		//}
		//fmt.Println()
		time.Sleep(time.Second * 1)
	}
}

func TestEtcd(t *testing.T) {
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 10,
	}
	client, err := clientv3.New(config)
	if err != nil {
		fmt.Errorf("clientv3.New err:%v", err)
	}

	getResp, err := client.Get(context.TODO(), "name", clientv3.WithPrefix())
	if err != nil {
		fmt.Errorf("clientv3.Get err:%v", err)
	}

	for _, v := range getResp.Kvs {
		fmt.Println(string(v.Value))
	}
}
