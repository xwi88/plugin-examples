package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/v8fg/rd"
	"github.com/v8fg/rd/config"
	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/xwi88/plugin-examples/internal/utils"
)

var (
	registerConfig config.RegisterConfig
	etcdConfig     clientV3.Config
)

func initServerParams() (err error) {
	flag.Parse()

	port = *Port
	if port == 0 {
		port, err = utils.FreePort("tcp")
		if err != nil {
			log.Fatalf("[%s] no free port:%v", moduleName, err)
			return err
		}
	}
	pid = os.Getpid()
	ppid = os.Getppid()
	localIP = utils.LocalIP()
	nodeName = *NodeName
	endpoints = strings.Split(*Endpoints, ",")
	serviceAddr = fmt.Sprintf("%s:%v", localIP, port)

	registerKey = fmt.Sprintf("%s://%v:%v", *ServiceKeyPrefix, localIP, port)
	registerName = *NameRegister
	registerVal = serviceAddr
	registerInterval = *Interval

	if len(nodeName) == 0 {
		nodeName = fmt.Sprintf("%v-%v", registerName, port)
	}

	log.Printf("[%s] init pid: %v, ppid: %v, ip:%v, serviceAddr: %v, registerName: %v, registerInterval: %v, etcd endPoints:%v\n",
		moduleName, pid, ppid, localIP, serviceAddr, registerName, registerInterval, endpoints)
	return nil
}

func initRegisterConfig() {
	messagesHandler := func(messages <-chan string) {
		for message := range messages {
			log.Printf("pid[%v], ppid[%v], messages consume, count:%v, content:%v\n", pid, ppid, len(messages), message)
		}
	}
	errorsHandler := func(errors <-chan error) {
		for errMsg := range errors {
			log.Printf("pid[%v], ppid[%v], errors consume, count:%v, content:%v\n", pid, ppid, len(errors), errMsg)
		}
	}
	// logger = log.New(log.Writer(), fmt.Sprintf("[%v] ", moduleName), log.LstdFlags|log.Lshortfile)
	logger := log.New(log.Writer(), fmt.Sprintf("[%v] ", moduleName), log.LstdFlags)
	logger = nil

	registerConfig = config.RegisterConfig{
		TTL: time.Second * 10,
		CommonConfig: config.CommonConfig{
			ChannelBufferSize: 64,
			ErrorsHandler:     errorsHandler,
			MessagesHandler:   messagesHandler,
			Logger:            logger,
		},
		MaxLoopTry: 16,
		MutableVal: true,
		// KeepAliveMode:     1,
	}

	// register init
	registerConfig.Name = registerName
	registerConfig.Key = registerKey
	registerConfig.Val = registerVal
	registerConfig.Return.Errors = true
	// registerConfig.Return.Messages = true
	registerConfig.KeepAlive.Interval = time.Second * 8
	registerConfig.KeepAlive.Mode = 1
	// registerConfig.Logger = logger

}

func registerServer(registerConfig config.RegisterConfig, client *clientV3.Client, etcdConfig clientV3.Config) error {
	return rd.RegisterEtcd(&registerConfig, nil, &clientV3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
	})
}

func runRegister() error {
	errs := rd.RegisterRun()
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

func runRegisterWithName(name string) error {
	errs := rd.RegisterRunWithParam(name)
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

func closeRegister() error {
	errs := rd.RegisterClose()
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

func closeRegisterWithName(name string) error {
	errs := rd.RegisterCloseWithParam(name)
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}
