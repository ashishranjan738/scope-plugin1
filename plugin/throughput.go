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

var tho = make(map[string]float64)
var storeAfter2 *Iops

type pvdata2 struct {
	read2  float64
	write2 float64
}

func getValues2(urlpassed string, query string) {
	res, err := http.Get(urlpassed)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	storeBefore2, err := getValue([]byte(body))
	if err != nil {
		panic(err.Error())
	}
	rand.Seed(time.Now().Unix())

	for _, result := range storeBefore2.Data.Result {
		tho[result.Metric.OpenebsPv], _ = strconv.ParseFloat(result.Value[1].(string), 32)
		tho[result.Metric.OpenebsPv] = tho[result.Metric.OpenebsPv] * 1024
	}

	if query == "OpenEBS_read_block_count_per_second" {
		read2ch <- tho
	} else if query == "OpenEBS_write_block_count_per_second" {
		write2ch <- tho
	}
}

func getPVs2() map[string]pvdata2 {
	// Get request to url
	queries := []string{"OpenEBS_read_block_count_per_second", "OpenEBS_write_block_count_per_second"}
	for _, query := range queries {
		go getValues2(url+query, query)
	}
	var read2, write2 map[string]float64
	for i := 0; i < len(queries); i++ {
		select {
		case read2 = <-read2ch:
		case write2 = <-write2ch:
		}
	}
	tho := make(map[string]pvdata2)
	if len(read2) == len(write2) {
		for k2, v2 := range read2 {
			meta2, err := clientset.CoreV1().PersistentVolumes().Get(k2, metav1.GetOptions{})
			if err != nil {
				continue
			}
			tho[string(meta2.UID)] = pvdata2{
				read2:  v2,
				write2: write2[k2],
			}
		}
	}
	return tho
}

func (p *Plugin) updatePVs2() {
	m2 := getPVs2()
	if len(m2) > 0 {
		p.pvs2 = m2
	}
}

func (p *Plugin) getTopologyPv2(str string) string {
	return fmt.Sprintf("%s;<persistent_volume>", str)
}
