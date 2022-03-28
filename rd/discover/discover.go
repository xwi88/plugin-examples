package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	clientV3 "go.etcd.io/etcd/client/v3"

	"github.com/v8fg/rd"

	"github.com/xwi88/plugin-examples/internal/utils"
)

var (
	discoverConfig rd.DiscoverConfig
	etcdConfig     clientV3.Config
)

func initClientParams() (err error) {
	flag.Parse()

	pid = os.Getpid()
	ppid = os.Getppid()
	localIP = utils.LocalIP()
	nodeName = *NodeName
	endpoints = strings.Split(*Endpoints, ",")

	nameDiscover = *NameDiscover
	scheme = *Scheme
	serviceKeyPrefix = *ServiceKeyPrefix
	interval = *Interval

	if len(nodeName) == 0 {
		nodeName = fmt.Sprintf("[pid=%v] [%s]", pid, moduleName)
	}

	etcdConfig = clientV3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 6,
	}

	log.Printf("[%s] init pid: %v, ppid: %v, ip:%s, interval:%v, nameDiscover: %s, scheme: %s, serviceKeyPrefix: %s, "+
		"etcd endPoints:%v\n", moduleName, pid, ppid, localIP, interval, nameDiscover, scheme, serviceKeyPrefix, endpoints)
	return err
}

func initDiscoverConfig() {
	addressesParser := func(key string, data []byte) (addr string, err error) {
		// fmt.Printf("pid[%v], ppid[%v], addressesParser consume, key:%s, data: %s\n", pid, ppid, key, data)
		return string(data), err
	}

	messagesHandler := func(messages <-chan string) {
		for message := range messages {
			log.Printf("pid[%v], messages consume, count:%v, content:%s\n", pid, len(messages), message)
		}
	}
	errorsHandler := func(errors <-chan error) {
		for errMsg := range errors {
			log.Printf("pid[%v], errors consume, count:%v, content:%s\n", pid, len(errors), errMsg)
		}
	}
	// logger = log.New(log.Writer(), fmt.Sprintf("[%v] ", moduleName), log.LstdFlags|log.Lshortfile)
	logger := log.New(log.Writer(), fmt.Sprintf("[%v] ", moduleName), log.LstdFlags)
	logger = nil

	discoverConfig = rd.DiscoverConfig{
		CommonConfig: rd.CommonConfig{
			ChannelBufferSize: 64,
			ErrorsHandler:     errorsHandler,
			MessagesHandler:   messagesHandler,
			Logger:            logger,
		},
		AddressesParser: addressesParser,
		ReturnResolve:   false,
	}

	// register init
	discoverConfig.Scheme = scheme
	discoverConfig.Name = nameDiscover
	discoverConfig.Service = serviceKeyPrefix
	discoverConfig.Return.Errors = true
	discoverConfig.Return.Messages = true
	// discoverConfig.Logger = logger
}

func registerDiscover(discoverConfig rd.DiscoverConfig, client *clientV3.Client, etcdConfig clientV3.Config) error {
	err := rd.DiscoverEtcd(&discoverConfig, client, &etcdConfig)
	return err
}

func runDiscover() error {
	errs := rd.DiscoverRun()
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

func runDiscoverWithName(name, scheme, service string) error {
	errs := rd.DiscoverRunWithParam(name, scheme, service)
	if len(errs) == 0 {
		return nil
	}
	return errs[0]
}

func closeDiscover() {
	rd.DiscoverClose()
}

func closeDiscoverWithName(name, scheme, service string) {
	rd.DiscoverCloseWithParam(name, scheme, service)
}
