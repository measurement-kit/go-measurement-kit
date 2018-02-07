package mk

import (
	"fmt"
	"testing"
)

func TestNettestStart(t *testing.T) {
	nt := Nettest{}
	nt.RegisterEventHandler(func(event interface{}) {
		fmt.Println(event)
	})
	c, err := nt.Start("Ndt")
	if err != nil {
		fmt.Println("err")
	}
	<-c
}
