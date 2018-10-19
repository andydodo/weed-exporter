package collectors

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type MasterCollector struct {
	Path      string
	MasterUp  prometheus.Gauge
	Volumes   *prometheus.GaugeVec
	Max       prometheus.Gauge
	Free      prometheus.Gauge
	Size      *prometheus.GaugeVec
	FileCount *prometheus.GaugeVec
}

func NewMasterCollector(path string) *MasterCollector {
	if path == "" {
		log.Fatalf("there is no path found")
		os.Exit(1)
	}
	return &MasterCollector{
		Path: path,
		MasterUp: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "Seaweedfs",
				Name:      "MasterUp",
				Help:      "Seaweedfs master Up",
			}),
		Volumes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "SeaWeedfs",
				Name:      "Volume",
				Help:      "Volume server have volumes",
			},
			[]string{"rack", "volume"},
		),
		Max: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "SeaWeedfs",
				Name:      "Max",
				Help:      "Max volumes",
			},
		),
		Free: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "SeaWeedfs",
				Name:      "Free",
				Help:      "Free volumes",
			},
		),
		Size: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "SeaWeedfs",
				Name:      "Size",
				Help:      "VolumeId size information",
			},
			[]string{"rack", "server", "volumeid"},
		),
		FileCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "SeaWeedfs",
				Name:      "FileCount",
				Help:      "VolumeId filecount information",
			},
			[]string{"rack", "server", "volumeid"},
		),
	}
}

func (c *MasterCollector) collect() error {
	c.MasterUp.Set(float64(1))
	_, err := net.Dial("tcp", c.Path)

	if err != nil {
		c.MasterUp.Set(float64(0))
		fmt.Println("dial master error")
	}

	client := http.Client{
		Timeout: time.Duration(3) * time.Second,
	}

	resp, err := client.Get("http://" + c.Path + "/vol/status?pretty=y")
	if err != nil {
		fmt.Println("curl api error")
		return nil
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body error")
		return nil
	}

	var weedInfo *Info
	err = json.Unmarshal(data, &weedInfo)
	if err != nil {
		fmt.Println("json convert to struct error")
		return err
	}

	max := weedInfo.Volumes.Max
	c.Max.Set(float64(max))
	free := weedInfo.Volumes.Free
	c.Free.Set(float64(free))

	for rack, volumes := range weedInfo.Volumes.DataCenters.DefaultDataCenter {
		for volume, vids := range volumes {
			c.Volumes.WithLabelValues(rack, volume).Set(float64(len(vids)))
			for _, v := range vids {
				//Todo: size == mb
				size := v.Size / 1024 / 1024
				id := strconv.FormatInt(v.ID, 10)
				c.Size.WithLabelValues(rack, volume, id).Set(float64(size))
				c.FileCount.WithLabelValues(rack, volume, id).Set(float64(v.FileCount))
			}
		}
	}
	return nil
}

func (c *MasterCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		c.MasterUp,
		c.Volumes,
		c.Max,
		c.Free,
		c.Size,
		c.FileCount,
	}
}

func (c *MasterCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.collectorList() {
		metric.Describe(ch)
	}
}

func (c *MasterCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(); err != nil {
		log.Fatalf("failed collecting cluster usage metrics: %s", err)
		return
	}
	for _, metric := range c.collectorList() {
		metric.Collect(ch)
	}
}

/* func DoExporter(path, addr string) {

	prometheus.MustRegister(NewMasterCollector(path))

	http.Handle("/metrics", prometheus.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("cannot start weedStatus exporter: %s", err)
	}

} */
