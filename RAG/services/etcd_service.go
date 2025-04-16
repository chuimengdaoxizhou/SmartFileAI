package services

//
//import (
//
//	//pb "RAG/allgrpc/allproto"
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/coreos/etcd/clientv3"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/resolver"
//	"log"
//)
//
//// EtcdResolver implements gRPC resolver for etcd-based service discovery
//type EtcdResolver struct {
//	client      *clientv3.Client
//	serviceName string
//	target      string
//	cc          grpc.ClientConnInterface
//	updateChan  chan []resolver.Address
//	ctx         context.Context
//	cancel      context.CancelFunc
//}
//
//// Builder implements resolver.Builder
//func (r *EtcdResolver) Build(target resolver.Target, cc grpc.ClientConnInterface, opts resolver.BuildOptions) (resolver.Resolver, error) {
//	r.target = target.URL.Host
//	r.cc = cc
//	r.updateChan = make(chan []resolver.Address, 1)
//	r.ctx, r.cancel = context.WithCancel(context.Background())
//
//	//go r.watch()
//	r.updateChan <- r.discoverServices()
//	return r, nil
//}
//
//func (r *EtcdResolver) Scheme() string {
//	return "etcd"
//}
//
//func (r *EtcdResolver) ResolveNow(options resolver.ResolveNowOptions) {
//	r.updateChan <- r.discoverServices()
//}
//
//func (r *EtcdResolver) Close() {
//	r.cancel()
//	r.client.Close()
//}
//
////func (r *EtcdResolver) watch() {
////	watchChan := r.client.Watch(r.ctx, fmt.Sprintf("/services/%s/", r.serviceName), clientv3.WithPrefix())
////	for {
////		select {
////		case <-r.ctx.Done():
////			return
////		case wresp := <-watchChan:
////			if wresp.Err() != nil {
////				log.Printf("Watch error: %v", wresp.Err())
////				continue
////			}
////			r.updateChan <- r.discoverServices()
////		case addrs := <-r.updateChan:
////			r.cc.UpdateState(resolver.State{Addresses: addrs})
////		}
////	}
////}
//
//func (r *EtcdResolver) discoverServices() []resolver.Address {
//	resp, err := r.client.Get(r.ctx, fmt.Sprintf("/services/%s/", r.serviceName), clientv3.WithPrefix())
//	if err != nil {
//		log.Printf("Failed to discover services: %v", err)
//		return nil
//	}
//
//	var addrs []resolver.Address
//	for _, kv := range resp.Kvs {
//		var serviceInfo map[string]interface{}
//		if err := json.Unmarshal(kv.Value, &serviceInfo); err != nil {
//			log.Printf("Failed to unmarshal service info: %v", err)
//			continue
//		}
//		addr := fmt.Sprintf("%s:%v", serviceInfo["host"], serviceInfo["port"])
//		addrs = append(addrs, resolver.Address{Addr: addr})
//	}
//	log.Printf("Discovered services: %v", addrs)
//	return addrs
//}
//
////func mainn() {
////	// 初始化 etcd 客户端
////	cfg := clientv3.Config{
////		Endpoints:   []string{"localhost:2379"},
////		DialTimeout: 5 * time.Second,
////	}
////	client, err := clientv3.New(cfg)
////	if err != nil {
////		log.Fatalf("Failed to connect to etcd: %v", err)
////	}
////	defer client.Close()
////
////	// 注册 etcd resolver
////	resolver.Register(&EtcdResolver{
////		client:      client,
////		serviceName: "grpc-data-service",
////	})
////
////	// 连接 gRPC 服务
////	conn, err := grpc.Dial(
////		"etcd:///grpc-data-service",
////		grpc.WithInsecure(),
////		grpc.WithBlock(),
////		grpc.WithTimeout(10*time.Second),
////	)
////	if err != nil {
////		log.Fatalf("Failed to connect to gRPC service: %v", err)
////	}
////	defer conn.Close()
////
////	// 创建 gRPC 客户端
////	dataClient := pb.NewDataManagementClient(conn)
////
////	// 测试 getDatabyPrompt
////	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
////	defer cancel()
////	resp, err := dataClient.getDatabyPrompt(ctx, &pb.Request{Prompt: "test prompt"})
////	if err != nil {
////		log.Fatalf("getDatabyPrompt failed: %v", err)
////	}
////	log.Printf("getDatabyPrompt response: %s", resp.Answer)
////
////	// 测试 updatabypath
////	resp, err = dataClient.updatabypath(ctx, &pb.Request{Prompt: "/home/chenyun/下载/train1.json"})
////	if err != nil {
////		log.Fatalf("updatabypath failed: %v", err)
////	}
////	log.Printf("updatabypath response: %s", resp.Answer)
////}
