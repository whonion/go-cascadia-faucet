package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Proxy struct {
	ID             string   `json:"_id"`
	IP             string   `json:"ip"`
	AnonymityLevel string   `json:"anonymityLevel"`
	ASN            string   `json:"asn"`
	City           string   `json:"city"`
	Country        string   `json:"country"`
	CreatedAt      string   `json:"created_at"`
	Google         bool     `json:"google"`
	ISP            string   `json:"isp"`
	LastChecked    int64    `json:"lastChecked"`
	Latency        float64  `json:"latency"`
	Org            string   `json:"org"`
	Port           string   `json:"port"`
	Protocols      []string `json:"protocols"`
	Region         string   `json:"region"`
	ResponseTime   int      `json:"responseTime"`
	Speed          int      `json:"speed"`
	UpdatedAt      string   `json:"updated_at"`
	WorkingPercent float64  `json:"workingPercent"`
	Uptime         float64  `json:"upTime"`
	UptimeSuccess  int      `json:"upTimeSuccessCount"`
	UptimeTry      int      `json:"upTimeTryCount"`
}

func main() {
	// Open the JSON file
	file, err := os.Open("Free_Proxy_List.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Read the file contents into a byte slice
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the JSON into a slice of Proxy objects
	var proxies []Proxy
	err = json.Unmarshal(bytes, &proxies)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a file to write the formatted data to
	outfile, err := os.Create("proxy.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outfile.Close()

	// Write the formatted data to the output file
	for _, p := range proxies {
		for _, protocol := range p.Protocols {
			protocol = strings.ToLower(protocol)
			fmt.Fprintf(outfile, "%s://%s:%s\n", protocol, p.IP, p.Port)
		}
	}
}
