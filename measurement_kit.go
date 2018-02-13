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
	Name    string
	Options NettestOptions
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

func notifyEventHandlers(event interface{}) {
	for i := 0; i < handleIndex; i++ {
		handleVals[i](event)
	}
}

// RegisterEventHandler will register an event handler
func (nt *Nettest) RegisterEventHandler(v func(interface{})) {
	newHandle(v)
}

var allEventTypes = []string{
	//	"QUEUED",
	//	"STARTED",
	"LOG",
	//	"CONFIGURED",
	//	"PROGRESS",
	"PERFORMANCE",
	//	"MEASUREMENT_ERROR",
	//	"REPORT_SUBMISSION_ERROR",
	//	"RESULT",
	//	"END",
}

// Run will run the test inside
func (nt *Nettest) Run() error {
	td := taskData{
		EnabledEvents: allEventTypes,
		Type:          nt.Name,
		Verbosity:     "INFO",
		Options:       nt.Options,
	}
	tdBytes, err := json.Marshal(td)
	if err != nil {
		return err
	}

	pTaskData := (*C.char)(unsafe.Pointer(&tdBytes[0]))
	task := C.mk_task_start(pTaskData)
	if task == nil {
		return errors.New("Got a null task data from mk_task_start")
	}

	for {
		event := C.mk_task_wait_for_next_event(task)
		if event == nil {
			return errors.New("Got a null event")
		}
		eventStr := C.GoString(C.mk_event_serialize(event))
		C.mk_event_destroy(event)
		if eventStr == "null" {
			break
		}

		var eventJSON map[string]interface{}
		if err := json.Unmarshal([]byte(eventStr), &eventJSON); err != nil {
			return err
		}
		notifyEventHandlers(eventJSON)
	}
	C.mk_task_destroy(task)
	return nil
}
