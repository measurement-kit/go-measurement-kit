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
			IncludeIP:        1,
			IncludeASN:       1,
			IncludeCountry:   1,
			DisableCollector: 1,
			SoftwareName:     "ooniprobe-dev",
			SoftwareVersion:  "0.0.1",
			GeoIPCountryPath: "",
			GeoASNPath:       "",
			OutputPath:       "/tmp/ooniprobe.output",
		},
	}
	fmt.Println("Starting")
	nt.On("log", func(event interface{}) {
		fmt.Println("Got log event", event)
	})
	nt.On("status.update.*", func(event interface{}) {
		fmt.Println("Got status update event", event)
	})
	nt.On("failure.*", func(event interface{}) {
		fmt.Println("Got a failure event", event)
	})
	if err := nt.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err)
	} else {
		fmt.Println("Exiting")
	}
}
