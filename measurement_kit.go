package main

/*
#include <measurement_kit/ffi.h>
*/
// #cgo LDFLAGS: -lmeasurement_kit
import "C"

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

type TaskData struct {
	EnabledEvents []string `json:"enabled_events"`
	Type          string   `json:"type"`
	Verbosity     string   `json:"verbosity"`
}

func main() {
	taskData := TaskData{
		EnabledEvents: []string{"LOG", "PERFORMANCE"},
		Type:          "Ndt",
		Verbosity:     "INFO",
	}
	taskDataBytes, err := json.Marshal(taskData)
	fmt.Printf("%s\n", taskDataBytes)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	task := C.mk_task_start((*C.char)(unsafe.Pointer(&taskDataBytes[0])))
	if task == nil {
		fmt.Println("err")
	}
	done := 0
	for done == 0 {
		event := C.mk_task_wait_for_next_event(task)
		if event == nil {
			panic("event is null")
		}
		eventJSON := C.GoString(C.mk_event_serialize(event))
		fmt.Println(eventJSON)
		if eventJSON == "null" {
			fmt.Println("Got null event")
			done = 1
		}
		C.mk_event_destroy(event)
	}
	C.mk_task_destroy(task)
}
