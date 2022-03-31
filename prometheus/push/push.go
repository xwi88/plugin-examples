package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/xwi88/plugin-examples/collector"
)

const moduleName = "plugin-examples-prometheus-push"

var (
	url = "http://localhost:9091"
	job = "plugin-examples-prometheus"

	pid int

	tickerCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ticker_count_push",
		Help: "collect the counter data",
	}, []string{"date"})
)

func main() {
	quit := make(chan struct{})
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	pid = os.Getpid()

	go func() {
		sig := <-ch
		log.Printf("[%s] pid[%v] signal[%v] killed", moduleName, pid, sig)
		quit <- struct{}{}
	}()

	tickerCount(quit)
}

func tickerCount(quit chan struct{}) {
	tk := time.Second * 1
	ticker := time.NewTicker(tk)
	defer ticker.Stop()
	rand.Seed(time.Now().UnixNano())
	collector.Run(int32(pid), tk)

loop:
	for {
		select {
		case <-ticker.C:
			addNum := rand.Int63n(1000)
			generateTickerCounter(addNum)
			prometheusPushData()
			log.Printf("[%s] pid[%v] ticker rand num:%v", moduleName, pid, addNum)
		case <-quit:
			collector.Close()
			break loop
		default:
			addNum := rand.Int63n(100)
			generateTickerCounter(addNum)
			log.Printf("[%s] pid[%v] default rand num:%v", moduleName, pid, addNum)
			time.Sleep(time.Millisecond * 300)
		}
	}
	log.Printf("[%s] pid[%v] tickerCount break for loop", moduleName, pid)
}

func generateTickerCounter(num int64) {
	nowDate := time.Now().Format("20060102")
	tickerCounter.WithLabelValues(nowDate).Add(float64(num))
}

func prometheusPushData() {
	if err := push.New(url, job).Collector(tickerCounter).
		Collector(collector.MyMemMetrics).
		Collector(collector.MyCpuMetrics).Push(); err != nil {
		log.Printf("[%s] prometheusPushData err:%v", moduleName, err)
	}

}
