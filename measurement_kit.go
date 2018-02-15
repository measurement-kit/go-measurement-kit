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
	"strings"
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
	Name           string
	Options        NettestOptions
	DisabledEvents []string
}

type taskData struct {
	DisabledEvents []string       `json:"disabled_events"`
	Type           string         `json:"type"`
	Verbosity      string         `json:"verbosity"`
	Options        NettestOptions `json:"options"`
}

var handleMu sync.Mutex
var handleMap = make(map[string][]interface{})

// Event is an event fired from measurement_kit
type Event struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func fire(s string, e Event) error {
	parts := strings.Split(s, ".")
	handles := make([]interface{}, 0)

	handleMu.Lock()
	// Add the literal match
	hs, ok := handleMap[s]
	if ok {
		handles = append(handles, hs...)
	}

	// Look for wildcards such as foo.bar.*
	for i := 1; i < len(parts); i++ {
		hn := fmt.Sprintf("%s.*", strings.Join(parts[0:i], "."))
		hs, ok := handleMap[hn]
		if ok {
			handles = append(handles, hs...)
		}
	}
	handleMu.Unlock()

	for _, handle := range handles {
		f := reflect.ValueOf(handle)
		args := make([]reflect.Value, 1)
		args[0] = reflect.ValueOf(e)

		// XXX should I do this call inside of a goroutine?
		values := f.Call(args)
		return values[0].Interface().(error)
	}
	return nil
}

func addHandler(s string, v interface{}) error {
	if _, ok := handleMap[s]; !ok {
		handleMap[s] = make([]interface{}, 0)
	}
	handleMap[s] = append(handleMap[s], v)
	return nil
}

// On will register an event handler
func (nt *Nettest) On(s string, v interface{}) error {
	handleMu.Lock()
	defer handleMu.Unlock()

	if reflect.ValueOf(v).Type().Kind() != reflect.Func {
		return errors.New("handler is not a function")
	}
	return addHandler(s, v)
}

// Run will run the test inside
func (nt *Nettest) Run() error {
	td := taskData{
		DisabledEvents: nt.DisabledEvents,
		Type:           nt.Name,
		Verbosity:      "INFO",
		Options:        nt.Options,
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

		var e Event
		if err := json.Unmarshal([]byte(eventStr), &e); err != nil {
			return err
		}
		fire(e.Key, e)
	}
	C.mk_task_destroy(task)
	return nil
}
