package main

import (
	"encoding/json"
	"fmt"

	"github.com/measurement-kit/go-measurement-kit"
)

type NDTSimple struct {
	Download float32 `json:"download"`
	Upload   float32 `json:"upload"`
	Ping     float32 `json:"ping"`
}

type NDTAdvanced struct {
	AvgRTT             float32 `json:"avg_rtt"`
	CongrestionLimited float32 `json:"congestion_limited"`
	FastRetran         float32 `json:"fast_retran"`
	MaxRTT             float32 `json:"max_rtt"`
	MinRTT             float32 `json:"min_rtt"`
	MSS                int64   `json:"mss"`

	OutOfOrder      float32 `json:"out_of_order"`
	PacketLoss      float32 `json:"packet_loss"`
	ReceiverLimited float32 `json:"receiver_limited"`
	SenderLimited   float32 `json:"sender_limited"`
	Timeouts        int64   `json:"timeouts"`
}

type NDTTestKeys struct {
	Simple   NDTSimple   `json:"simple"`
	Advanced NDTAdvanced `json:"advanced"`
}
type NDTMeasurement struct {
	TestKeys NDTTestKeys `json:"test_keys"`
}

func main() {
	var ndtMeasurement NDTMeasurement

	nt := mk.NewNettest("Ndt")
	nt.Options = mk.NettestOptions{
			CaBundlePath:     "/etc/ssl/cert.pem",
			IncludeIP:        false,
			IncludeASN:       true,
			IncludeCountry:   true,
			DisableCollector: false,
			GeoIPCountryPath: "",
			GeoIPASNPath:     "",
			OutputPath:       "/tmp/ooniprobe-report.jsonl",
			LogLevel:         "INFO",
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
		jsonBlob := []byte(event.Value["json_str"].(string))
		json.Unmarshal(jsonBlob, &ndtMeasurement)
	})
	nt.On("*", func(event interface{}) {
		fmt.Println("[*]", event)
	})
	if err := nt.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err)
	} else {
		fmt.Println("Done")
	}

	fmt.Printf("Results\n")
	fmt.Printf("=======\n")
	fmt.Printf("Upload: %f\n", ndtMeasurement.TestKeys.Simple.Upload)
	fmt.Printf("Download: %f\n", ndtMeasurement.TestKeys.Simple.Download)
	fmt.Printf("Ping: %f\n", ndtMeasurement.TestKeys.Simple.Ping)
}
