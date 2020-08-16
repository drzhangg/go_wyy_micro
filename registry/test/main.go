package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	config := clientv3.Config{
		DialTimeout: time.Second * 10,
	}
	client, err := clientv3.New(config)
	if err != nil {
		fmt.Errorf("clientv3.New err:%v", err)
	}
	fmt.Println(333)

	//putres, err := client.Put(context.TODO(), "name", "hello")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(putres.PrevKv.Value))

	getResp, err := client.Get(context.TODO(), "name1", clientv3.WithPrefix())
	if err != nil {
		fmt.Errorf("clientv3.Get err:%v", err)
	}
	fmt.Println(222)

	if len(getResp.Kvs) > 0 {
		fmt.Println(fmt.Println(getResp.Kvs))
	}else {
		_, err := client.Put(context.TODO(), "name1", "hello")
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(putres.PrevKv.Value))
	}

	for _, v := range getResp.Kvs {
		fmt.Println(111)
		fmt.Println(string(v.Value))
	}
}
