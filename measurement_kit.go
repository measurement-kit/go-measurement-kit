package mk

/*
#include <measurement_kit/ffi.h>
*/
// #cgo LDFLAGS: -lmeasurement_kit
import "C"

import (
	"encoding/json"
	"errors"
	"sync"
	"unsafe"
)

// NettestOptions are the options to be passed to a particular nettest
type NettestOptions struct {
	IncludeIP        bool   `json:"save_real_probe_ip"`
	IncludeASN       bool   `json:"save_real_probe_asn"`
	IncludeCountry   bool   `json:"save_real_probe_cc"`
	DisableCollector bool   `json:"no_collector"`
	SoftwareName     string `json:"software_name"`
	SoftwareVersion  string `json:"software_version"`

	GeoIPCountryPath string `json:"geoip_country_path"`
	GeoASNPath       string `json:"geoip_asn_path"`
	OutputPath       string `json:"output_path"`
	CaBundlePath     string `json:"net/ca_bundle_path"`
}

// Nettest is a wrapper for running a particular nettest
type Nettest struct {
	Options NettestOptions
	done    chan bool
	err     error
}

type taskData struct {
	EnabledEvents []string       `json:"enabled_events"`
	Type          string         `json:"type"`
	Verbosity     string         `json:"verbosity"`
	Options       NettestOptions `json:"options"`
}

var handleLock sync.Mutex
var handleVals = make(map[int]func(interface{}))
var handleIndex int

func newHandle(v func(interface{})) int {
	handleLock.Lock()
	defer handleLock.Unlock()
	i := handleIndex
	handleIndex++
	handleVals[i] = v
	return i
}

func notifyEventHandlers(event string) {
	for i := 0; i < handleIndex; i++ {
		handleVals[i](event)
	}
}

// RegisterEventHandler will register an event handler
func (nt *Nettest) RegisterEventHandler(v func(interface{})) {
	newHandle(v)
}

// Start will start the test inside of a goroutine
func (nt *Nettest) Start(name string) (chan bool, error) {
	nt.done = make(chan bool, 1)

	td := taskData{
		EnabledEvents: []string{"LOG", "PERFORMANCE"},
		Type:          name,
		Verbosity:     "INFO",
		Options:       nt.Options,
	}
	tdBytes, err := json.Marshal(td)
	if err != nil {
		return nt.done, err
	}

	pTaskData := (*C.char)(unsafe.Pointer(&tdBytes[0]))
	task := C.mk_task_start(pTaskData)
	if task == nil {
		return nt.done, errors.New("Got a null task data from mk_task_start")
	}

	go func() {
		for {
			event := C.mk_task_wait_for_next_event(task)
			if event == nil {
				nt.err = errors.New("Got a null event")
				break
			}
			eventJSON := C.GoString(C.mk_event_serialize(event))
			C.mk_event_destroy(event)
			if eventJSON == "null" {
				break
			}
			notifyEventHandlers(eventJSON)
		}
		C.mk_task_destroy(task)
		nt.done <- true
	}()

	return nt.done, nil
}
