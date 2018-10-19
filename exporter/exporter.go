package exporter

import (
	"log"
	"net/http"
	"sync"

	"github.com/andy/weed-exporter/exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type WeedExporter struct {
	mu         sync.Mutex
	collectors []prometheus.Collector
}

var _ prometheus.Collector = &WeedExporter{}

func NewWeedExporter(server, path string) *WeedExporter {
	var exporter *WeedExporter
	var pathes []string
	parts := []string{":8001", ":8002", ":8003", ":8004", ":8005", ":8006", ":8007", ":8008", ":8009", ":8010",
		":8011", "8012", "8013", "8014", "8015", "8016", "8017", "8018", "8019", "8020", "8021", "8022"}

	for _, part := range parts {
		ip := path
		ip += part
		pathes = append(pathes, ip)
	}

	switch server {
	case "master":
		exporter = &WeedExporter{
			collectors: []prometheus.Collector{
				collectors.NewMasterCollector(path),
			},
		}
	case "volume":
		exporter = &WeedExporter{
			collectors: []prometheus.Collector{
				collectors.NewVolumeCollector(pathes),
			},
		}
	}
	return exporter
}

func (c *WeedExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, cc := range c.collectors {
		cc.Describe(ch)
	}
}

func (c *WeedExporter) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, cc := range c.collectors {
		cc.Collect(ch)
	}
}

func DoExporter(path, addr, server string) {

	prometheus.MustRegister(NewWeedExporter(server, path))

	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("cannot start weed exporter: %s", err)
	}

}
