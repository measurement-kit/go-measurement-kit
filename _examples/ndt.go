package main

import (
	"fmt"

	"github.com/measurement-kit/go-measurement-kit"
)

func main() {
	nt := mk.Nettest{
		Name: "Ndt",
	}
	nt.On("log", func(event interface{}) {
		fmt.Println("Got log event", event)
	})
	nt.On("status.update.*", func(event interface{}) {
		fmt.Println("Got status update event", event)
	})
	if err := nt.Run(); err != nil {
		fmt.Printf("Got error: %s\n", err)
	}
}
