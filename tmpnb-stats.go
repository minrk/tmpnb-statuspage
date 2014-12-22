package main

// Expects these environment variables
// STATUS_PAGE_API_KEY
// STATUS_PAGE_PAGE_ID
// STATUS_PAGE_TMPNB_METRIC_ID

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Simple config for statuspage.io
// MetricIDs maps an environment variable to the metric ID
type StatusPage struct {
	APIKey        string `json:"apiKey" envconfig:"API_KEY"`
	PageID        string `json:"pageID" envconfig:"PAGE_ID"`
	TmpnbMetricID string `json:"tmpnbMetricID" envconfig:"TMPNB_METRIC_ID"`
}

// Adhering to the tmpnb-redirector "standard"
type Usage struct {
	Available int              `json:"available"`
	Version   string           `json:"version"`
	Capacity  int              `json:"capacity"`
	Hosts     map[string]Usage `json:"hosts"`
}

type TmpnbStats struct {
	StatsEndpoint string        `envconfig:"ENDPOINT"`
	Period        time.Duration `envconfig:"PERIOD"`
}

func (t TmpnbStats) usage() Usage {
	resp, err := http.Get(t.StatsEndpoint)

	if err != nil {
		log.Fatalf("Unable to reach endpoint: %v\n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read body from endpoint: %v\n", err)
	}

	var usage Usage
	err = json.Unmarshal(body, &usage)
	if err != nil {
		log.Fatalf("Unable to parse tmpnb JSON from endpoint: %v\n", err)
	}
	return usage
}

func (t TmpnbStats) percentAvailable() int {
	usage := t.usage()
	return (usage.Available * 100) / usage.Capacity
}

func main() {
	var statusPage StatusPage
	var tmpnb TmpnbStats

	if err := envconfig.Process("STATUS_PAGE", &statusPage); err != nil {
		log.Fatalf("Unable to process status page env: %v\n", err)
	}

	if err := envconfig.Process("TMPNB_STATS", &tmpnb); err != nil {
		log.Fatalf("Unable to process tmpnb env: %v\n", err)
	}

	ticker := time.NewTicker(time.Second * tmpnb.Period)

	go func() {
		for _ = range ticker.C {
			log.Printf("%v availability %v%%", tmpnb.StatsEndpoint, tmpnb.percentAvailable())
		}
	}()

	//Forever young
	select {}

}
