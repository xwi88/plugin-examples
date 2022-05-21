package collector

import (
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/v8fg/collectd"
)

// Collect freq, shall not less than 5 seconds

var nodeCollector *NodeCollector
var onceCh sync.Once

type NodeCollector struct {
	container bool
	done      chan struct{}
	pid       int32
}

var (
	MyCpuMetrics = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "my_cpu"}, []string{"tag"})
	MyMemMetrics = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "my_mem"}, []string{"tag"})
)

// Run start node collector to collect some metrics
func Run(pid int32, ticker time.Duration) {
	container := collectd.RunningInDockerContainerPid(pid)
	nodeCollector = &NodeCollector{done: make(chan struct{}), container: container, pid: pid}
	log.Printf("nodeCollector run with pid: %v, in container: %v", pid, container)

	go func() {
		ticker := time.NewTicker(ticker)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updateMyMetrics()
			case <-nodeCollector.done:
				return
			}
		}
	}()
}

func Close() {
	onceCh.Do(func() {
		nodeCollector.done <- struct{}{}
	})
}

func updateMyMetrics() {
	cif := collectd.CpuInfo(nodeCollector.container, nodeCollector.pid)
	mif := collectd.MemInfo(nodeCollector.container, nodeCollector.pid)

	log.Printf("sendMyMetrics cpu info: %+v", cif)
	log.Printf("sendMyMetrics mem info: %+v", mif)

	MyCpuMetrics.With(prometheus.Labels{"tag": "count"}).Set(cif.Count)
	MyCpuMetrics.With(prometheus.Labels{"tag": "used"}).Set(cif.Used)
	MyCpuMetrics.With(prometheus.Labels{"tag": "percent"}).Set(cif.Percent)
	MyCpuMetrics.With(prometheus.Labels{"tag": "container"}).Set(float64(cif.Container))

	MyMemMetrics.With(prometheus.Labels{"tag": "total"}).Set(mif.Total)
	MyMemMetrics.With(prometheus.Labels{"tag": "used"}).Set(mif.Used)
	MyMemMetrics.With(prometheus.Labels{"tag": "percent"}).Set(mif.Percent)
	MyMemMetrics.With(prometheus.Labels{"tag": "container"}).Set(float64(mif.Container))
}
