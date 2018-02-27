package main

import (
	"fmt"

	"github.com/measurement-kit/go-measurement-kit"
)

func main() {
	nt := mk.Nettest{
		Name: "Ndt",
		Options: mk.NettestOptions{
			CaBundlePath:     "/etc/ssl/cert.pem",
			IncludeIP:        false,
			IncludeASN:       true,
			IncludeCountry:   true,
			DisableCollector: true,
			SoftwareName:     "ooniprobe-dev",
			SoftwareVersion:  "0.0.1",
			GeoIPCountryPath: "",
			GeoIPASNPath:     "",
			OutputPath:       "/tmp/ooniprobe-report.jsonl",
			LogLevel:         "DEBUG2",
		},
	}
	nt.On("log", func(event interface{}) {
		fmt.Println("Got log event", event)
	})
	nt.On("status.update.*", func(event interface{}) {
		fmt.Println("Got status update event", event)
	})
	nt.On("failure.*", func(event interface{}) {
		fmt.Println("Got a failure event", event)
	})
	nt.On("measurement_entry", func(event interface{}) {
		fmt.Println("Got measurement_entry event", event)
	})
	nt.On("*", func(event interface{}) {
		fmt.Println("Catch all event", event)
	})
	if err := nt.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err)
	} else {
		fmt.Println("Exiting")
	}
}
