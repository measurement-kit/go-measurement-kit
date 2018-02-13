package mk

import (
	"fmt"
	"testing"
)

func TestNdtRun(t *testing.T) {
	nt := Nettest{
		Name: "Ndt",
	}
	nt.RegisterEventHandler(func(event interface{}) {
		fmt.Println("Got event", event)
	})
	if err := nt.Run(); err != nil {
		t.Fatalf("Got error: %s", err)
	}
}
