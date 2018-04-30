package main

import (
	"fmt"

	"github.com/measurement-kit/go-measurement-kit"
)

func main() {
	nt := mk.NewNettest("WebConnectivity")
	nt.Options = mk.NettestOptions{
		CaBundlePath:     "cert.pem",
		IncludeIP:        false,
		IncludeASN:       true,
		IncludeCountry:   true,
		DisableCollector: false,
		GeoIPCountryPath: "",
		GeoIPASNPath:     "",
		OutputPath:       "/tmp/ooniprobe-report.jsonl",
		LogLevel:         "INFO",
		Inputs:           []string{"https://ooni.torproject.org/"},
	}

	nt.On("log", func(event interface{}) {
		fmt.Println("[log]", event)
	})
	nt.On("status.update.*", func(event interface{}) {
		fmt.Println("[status.update.*]", event)
	})
	nt.On("failure.*", func(event interface{}) {
		fmt.Println("[failure.*]", event)
	})
	nt.On("measurement", func(event mk.Event) {
		fmt.Println("[measurement]", event)
	})
	nt.On("*", func(event interface{}) {
		fmt.Println("[*]", event)
	})
	if err := nt.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err)
	} else {
		fmt.Println("Done")
	}
}
