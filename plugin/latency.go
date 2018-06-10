package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var lat = make(map[string]float64)

type pvdata1 struct {
	read1  float64
	write1 float64
}

func getValues1(urlpassed string, query string) {
	res, err := http.Get(urlpassed)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	storeBefore1, err := getValue([]byte(body))
	if err != nil {
		panic(err.Error())
	}

	rand.Seed(time.Now().Unix())

	for _, result := range storeBefore1.Data.Result {
		lat[result.Metric.OpenebsPv], _ = strconv.ParseFloat(result.Value[1].(string), 32)
	}

	if query == "OpenEBS_read_latency" {
		read1ch <- lat
	} else if query == "OpenEBS_write_latency" {
		write1ch <- lat
	}
}

func getPVs1() map[string]pvdata1 {
	// Get request to url
	queries := []string{"OpenEBS_read_latency", "OpenEBS_write_latency"}
	for _, query := range queries {
		go getValues1(url+query, query)
	}
	var read1, write1 map[string]float64
	for i := 0; i < len(queries); i++ {
		select {
		case read1 = <-read1ch:
		case write1 = <-write1ch:
		}
	}
	lat := make(map[string]pvdata1)
	if len(read1) == len(write1) {
		for k1, v1 := range read1 {
			meta1, err := clientset.CoreV1().PersistentVolumes().Get(k1, metav1.GetOptions{})
			if err != nil {
				continue
			}
			lat[string(meta1.UID)] = pvdata1{
				read1:  v1,
				write1: write1[k1],
			}
		}
	}
	return lat
}

func (p *Plugin) updatePVs1() {
	m1 := getPVs1()
	if len(m1) > 0 {
		p.pvs1 = m1
	}
}

func (p *Plugin) getTopologyPv1(str string) string {
	return fmt.Sprintf("%s;<persistent_volume>", str)
}
