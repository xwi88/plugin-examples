package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/atomic"
	"google.golang.org/grpc"

	pb "github.com/xwi88/plugin-examples/internal/pb/go"
)

const moduleName = "rd-discover-client"

// input params
var (
	defaultEtcdPoints = []string{
		"http://127.0.0.1:2379",
	}
	Endpoints        = flag.String("endpoints", strings.Join(defaultEtcdPoints, ","), "etcd endpoints")
	NameDiscover     = flag.String("discoverName", "discover-plugin-examples-rd-v1.0-grpc", "discover name")
	Scheme           = flag.String("scheme", "services", "service scheme, default services")
	ServiceKeyPrefix = flag.String("serviceKey", "xwi88:plugin-examples/rd/v1.0/grpc", "service key")
	NodeName         = flag.String("nodeName", "", "server name")
	Interval         = flag.Duration("interval", time.Second*10, "interval to call server")
)

var (
	pid       int
	ppid      int
	localIP   string
	nodeName  string
	endpoints []string

	scheme           string
	nameDiscover     string
	serviceKeyPrefix string
	interval         time.Duration
)

func main() {
	flag.Parse()

	err := initClientParams()
	if err != nil {
		log.Fatalf("[%s] initClientParam err:%v", moduleName, err)
	}

	initDiscoverConfig()
	err = registerDiscover(discoverConfig, nil, etcdConfig)
	if err != nil {
		log.Fatalf("[%s] registerDiscover err:%v", moduleName, err)
	}

	err = runDiscover()
	if err != nil {
		log.Fatalf("[%s] runDiscover err:%v", moduleName, err)
	}

	quit := make(chan struct{})
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-ch
		closeDiscover()
		log.Printf("[%s] [pid=%v] sig[%v] killed", moduleName, pid, sig)
		quit <- struct{}{}
	}()

	err = runGRPCClient(quit)
	if err != nil {
		log.Fatalf("[%s] runGRPCClient err: %s", moduleName, err.Error())
	}

	log.Printf("[%s] exit", moduleName)
}

func runGRPCClient(quit <-chan struct{}) error {
	count := atomic.NewUint64(0)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	newWithDefaultServiceConfig := grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`)
	// _, err := grpc.DialContext(context.TODO(), scheme+"://authority/"+service, newWithDefaultServiceConfig, grpc.WithBlock(), grpc.WithInsecure())
	conn, err := grpc.DialContext(context.TODO(), scheme+"://authority/"+serviceKeyPrefix, newWithDefaultServiceConfig, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[%s] runDiscover err:%v", moduleName, err)
	}

loop:
	for {
		if conn == nil {
			conn, err = grpc.DialContext(context.TODO(), scheme+"://authority/"+serviceKeyPrefix, newWithDefaultServiceConfig, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("[%s] runDiscover err:%v", moduleName, err)
			}
		}

		select {
		case <-ticker.C:
			client := pb.NewGreeterClient(conn)
			resp, err := client.SayHello(context.Background(),
				&pb.HelloRequest{
					Ip:       localIP,
					NodeName: nodeName,
					Name:     nameDiscover,
					Command:  "",
					Message:  fmt.Sprintf("client request %v", time.Now().Format("2006-01-02T15:04:05.000+0800")),
				},
			)
			if err == nil {
				log.Printf("[%s] pid[%v] receive response %s\n", moduleName, pid, resp.String())
			} else {
				log.Printf("[%s] call server error:%s\n", moduleName, err)
			}
		case <-quit:
			log.Printf("exit runGRPCClient loop")
			break loop
		}
		count.Add(1)
		log.Printf("loop running:%v", count.Load())
	}
	return errors.New("exit runGRPCClient loop")
}
