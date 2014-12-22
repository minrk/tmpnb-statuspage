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

type TMPNB struct {
	StatsEndpoint string `envconfig:"STATS_ENDPOINT"`
}

func (s StatusPage) init() {
	err := envconfig.Process("STATUS_PAGE", &s)

	if err != nil {
		log.Fatalf("Unable to process status page env: %v\n", err)
	}
}

func (t TMPNB) init() {
	err := envconfig.Process("TMPNB", &t)
	if err != nil {
		log.Fatalf("Unable to process tmpnb env: %v\n", err)
	}

	resp, err := http.Get(t.StatsEndpoint)

	if err != nil {
		log.Fatalf("Unable to reach endpoint initially: %v\n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read from endpoint: %v\n", err)
	}

	var usage Usage
	err = json.Unmarshal(body, &usage)
	if err != nil {
		log.Fatalf("Unable to parse JSON body from endpoint: %v\n", err)
	}

	log.Println(usage)
}

func main() {
	var statusPage StatusPage
	var tmpnb TMPNB

	statusPage.init()
	tmpnb.init()

}
