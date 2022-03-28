package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"

	pb "github.com/xwi88/plugin-examples/internal/pb/go"
)

// input params
var (
	defaultEtcdPoints = []string{
		"http://127.0.0.1:2379",
	}
	Endpoints        = flag.String("endpoints", strings.Join(defaultEtcdPoints, ","), "etcd endpoints")
	NameRegister     = flag.String("registerName", "register-plugin-examples-rd-v1.0-grpc", "register name")
	ServiceKeyPrefix = flag.String("serviceKey", "/services/xwi88:plugin-examples/rd/v1.0/grpc", "service key")
	Port             = flag.Int("port", 0, "listening port, 0: use random port")
	NodeName         = flag.String("nodeName", "", "server name")
	Interval         = flag.Duration("interval", time.Second*10, "register interval")
)

const moduleName = "rd-register-server"

var (
	pid         int
	ppid        int
	localIP     string
	port        int
	endpoints   []string
	serviceAddr string

	registerKey      string
	registerName     string
	registerVal      string
	registerInterval time.Duration
	nodeName         string
)

// server is used to implement hello.GreeterServer.
type server struct {
	pb.GreeterServer
	pid  int
	port int
}

// SayHello implements hello.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("[%s] pid[%v] Received: %v\n", moduleName, pid, in.String())

	return &pb.HelloReply{
		NodeName: nodeName,
		Message:  "server reply message: " + in.GetMessage(),
		Ip:       serviceAddr}, nil
}

func main() {
	var err error
	err = initServerParams()
	if err != nil {
		log.Fatalf("[%s] initParam failed, err:%v", moduleName, err)
	}

	initRegisterConfig()

	err = registerServer(registerConfig, nil, etcdConfig)
	if err != nil {
		log.Fatalf("[%s] registerServer failed, err:%v", moduleName, err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("[%s] failed to listen: %v", moduleName, err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{pid: pid, port: port})

	quit := make(chan struct{})
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-ch
		err = closeRegister()
		if err != nil {
			log.Fatalf("[%s] failed to RegisterClose: %s", moduleName, err.Error())
		}
		s.GracefulStop()
		log.Printf("[%s] [pid=%v] sig[%v] killed", moduleName, pid, sig)
		quit <- struct{}{}
	}()

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("[%s] failed to serve: %s", moduleName, err.Error())
		}
	}()

	err = runRegister()
	if err != nil {
		log.Fatalf("[%s] runRegister err: %s", moduleName, err.Error())
	}
	<-quit
	log.Printf("[%s] exit", moduleName)
}

func exit() error {
	return syscall.Kill(syscall.Getpid(), syscall.SIGKILL)
}
