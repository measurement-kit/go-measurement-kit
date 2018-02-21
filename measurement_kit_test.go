package mk

import (
	"fmt"
	"testing"
)

func TestNdtRun(t *testing.T) {
	nt := Nettest{
		Name: "Ndt",
	}
	nt.On("log", func(event interface{}) {
		fmt.Println("Got log event", event)
	})
	nt.On("status.update.*", func(event interface{}) {
		fmt.Println("Got status update event", event)
	})
	if err := nt.Run(); err != nil {
		t.Fatalf("Got error: %s", err)
	}
}
