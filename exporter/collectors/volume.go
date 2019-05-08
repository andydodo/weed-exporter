package collectors

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net"
	"os"
)

type VolumeCollector struct {
	Path     []string
	VolumeUp prometheus.Gauge
	//Todo: monitor volume down
	//VolumeDown prometheus.Gauge
	VolumeDown *prometheus.GaugeVec
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
				Namespace: "SeaWeedfs",
				Name:      "VolumeUp",
				Help:      "Seaweedfs volume Up",
			}),
		VolumeDown: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "SeaWeedfs",
				Name:      "VolumeDown",
				Help:      "Seaweedfs volume Down",
			},
			[]string{"volumeNode"},
		),
	}
}

func (c *VolumeCollector) collect() error {
	c.VolumeUp.Set(float64(len(c.Path)))
	//Todo: monitor volume down
	//c.VolumeDown.WithLabelValues("").Set(float64(0))
	for _, path := range c.Path {
		_, err := net.Dial("tcp", path)

		if err != nil {
			c.VolumeUp.Dec()
			//Todo: monitor volume down
			c.VolumeDown.WithLabelValues(path).Set(float64(0))
		} else {
			//Todo: monitor volume down
			c.VolumeDown.WithLabelValues(path).Set(float64(1))
		}
	}
	return nil
}

func (c *VolumeCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		c.VolumeUp,
		c.VolumeDown,
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
