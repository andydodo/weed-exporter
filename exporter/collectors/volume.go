package collectors

import (
	"fmt"
	"net"
	"os"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type VolumeCollector struct {
	Path     []string
	VolumeUp prometheus.Gauge
}

func NewVolumeCollector(path []string) *VolumeCollector {
	if len(path) == 0 {
		fmt.Println("there is no Seaweedfs volume found")
		os.Exit(1)
	}
	return &VolumeCollector{
		Path: path,
		VolumeUp: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "Seaweedfs",
				Name:      "VolumeUp",
				Help:      "Seaweedfs volume Up",
			}),
	}
}

func (c *VolumeCollector) collect() error {
	c.VolumeUp.Set(float64(len(c.Path)))
	for _, path := range c.Path {
		_, err := net.Dial("tcp", path)

		if err != nil {
			c.VolumeUp.Dec()
			fmt.Println("dial volume error")
		}
	}
	return nil
}

func (c *VolumeCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		c.VolumeUp,
	}
}

func (c *VolumeCollector) Describe(ch chan<- *prometheus.Desc) {
	//c.VolumeUp.Describe(ch)
	for _, metrics := range c.collectorList() {
		metrics.Describe(ch)
	}
}

func (c *VolumeCollector) Collect(ch chan<- prometheus.Metric) {
	//c.VolumeUp.Collect(ch)
	if err := c.collect(); err != nil {
		log.Fatalf("failed collecting cluster usage metrics: %s", err)
		return
	}
	for _, metrics := range c.collectorList() {
		metrics.Collect(ch)
	}
}

/* func DoExporter(host , addr  string) {
	var path []string
	parts := []string{":8080", ":8081", ":8082", ":8083", ":8084", ":8085", ":8086", ":8087", ":8088", ":8089", ":8090", ":8091"}

	for _, part := range parts {
		ip := host
		ip += part
		path = append(path, ip)
	}

	prometheus.MustRegister(NewVolumeCollector(path))

	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("cannot start Seaweedfs exporter: %s", err)
	}

} */
