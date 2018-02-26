package mk

/*
#include <measurement_kit/ffi.h>
*/
// #cgo LDFLAGS: -lmeasurement_kit
import "C"

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

// NettestOptions are the options to be passed to a particular nettest
type NettestOptions struct {
	IncludeIP        int    `json:"save_real_probe_ip"`
	IncludeASN       int    `json:"save_real_probe_asn"`
	IncludeCountry   int    `json:"save_real_probe_cc"`
	DisableCollector int    `json:"no_collector"`
	SoftwareName     string `json:"software_name"`
	SoftwareVersion  string `json:"software_version"`

	GeoIPCountryPath string `json:"geoip_country_path"`
	GeoASNPath       string `json:"geoip_asn_path"`
	OutputPath       string `json:"output_filepath"`
	CaBundlePath     string `json:"net/ca_bundle_path"`
}

// Nettest is a wrapper for running a particular nettest
type Nettest struct {
	Name           string
	Options        NettestOptions
	DisabledEvents []string
}

// On will register an event handler
// The possible events are:
// - log
// - status.queued
// - status.started
// - status.report_created
// - status.geoip_lookup
// - status.progress
// - status.update.performance
// - status.update.websites
// - status.end
// - failure.measurement
// - failure.report_submission
// - entry
// It is possible to register events with wildcards.
// For example On("status.*", ...) will fire on status.queued, status.started, ...
func (nt *Nettest) On(s string, v interface{}) error {
	handleMu.Lock()
	defer handleMu.Unlock()

	if reflect.ValueOf(v).Type().Kind() != reflect.Func {
		return errors.New("handler is not a function")
	}
	return addHandler(s, v)
}

// Event is an event fired from measurement_kit
type Event struct {

	// Is the key for the event. See On for the list of possible events.
	Key string `json:"type"` // XXX rename this

	// Contains the value for the fired event
	Value map[string]interface{} `json:"value"`
}

// Run will run the test inside
func (nt *Nettest) Run() error {
	td := taskData{
		Type:      nt.Name,
		Verbosity: "DEBUG2",
		Options:   nt.Options,
	}
	if nt.DisabledEvents != nil {
		td.DisabledEvents = nt.DisabledEvents
	} else {
		td.DisabledEvents = make([]string, 0)
	}

	tdBytes, err := json.Marshal(td)
	if err != nil {
		return err
	}
	fmt.Printf("eventStr: %s\n", string(tdBytes))

	pTaskData := (*C.char)(unsafe.Pointer(&tdBytes[0]))
	task := C.mk_task_start(pTaskData)
	defer C.mk_task_destroy(task)
	if task == nil {
		return errors.New("Got a null task data from mk_task_start")
	}

	for {
		isDone := C.mk_task_is_done(task)
		if isDone == 1 {
			break
		}

		event := C.mk_task_wait_for_next_event(task)
		defer C.mk_event_destroy(event)
		if event == nil {
			return errors.New("Got a null event")
		}

		eventStr := C.GoString(C.mk_event_serialize(event))
		fmt.Printf("eventStr: %s\n", eventStr)

		var e Event
		if err := json.Unmarshal([]byte(eventStr), &e); err != nil {
			return err
		}
		fire(e.Key, e)
	}
	return nil
}
