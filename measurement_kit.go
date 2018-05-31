package mk

/*
#include <measurement_kit/ffi.h>
*/
// #cgo darwin LDFLAGS: -lc++
// #cgo darwin,amd64 CFLAGS: -I${SRCDIR}/libs/macos/measurement-kit/0.9.0-alpha.2-1/x86_64/include/
// #cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/macos/measurement-kit/0.9.0-alpha.2-1/x86_64/lib/libmeasurement_kit.a
// #cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/macos/geoip-api-c/1.6.12-3/x86_64/lib/libGeoIP.a
// #cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/macos/libressl/2.6.4-4/x86_64/lib/libssl.a
// #cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/macos/libressl/2.6.4-4/x86_64/lib/libcrypto.a
// #cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/macos/libevent/2.1.8-4/x86_64/lib/libevent.a
// #cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/macos/libevent/2.1.8-4/x86_64/lib/libevent_openssl.a
// #cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/macos/libevent/2.1.8-4/x86_64/lib/libevent_pthreads.a
//
// #cgo windows LDFLAGS: -static
// #cgo windows,amd64 CFLAGS: -I${SRCDIR}/libs/mingw/measurement-kit/0.9.0-alpha.2-1/x86_64/include/
// #cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/mingw/measurement-kit/0.9.0-alpha.2-1/x86_64/lib/libmeasurement_kit.a
// #cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/mingw/geoip-api-c/1.6.12-3/x86_64/lib/libGeoIP.a
// #cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/mingw/libressl/2.6.4-4/x86_64/lib/libssl.a
// #cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/mingw/libressl/2.6.4-4/x86_64/lib/libcrypto.a
// #cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/mingw/libevent/2.1.8-4/x86_64/lib/libevent.a
// #cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/mingw/libevent/2.1.8-4/x86_64/lib/libevent_openssl.a
// #cgo windows,amd64 LDFLAGS: -lws2_32 -lstdc++
import "C"

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

// NettestOptions are the options to be passed to a particular nettest
type NettestOptions struct {
	IncludeIP        bool
	IncludeASN       bool
	IncludeCountry   bool
	ProbeCC          string
	ProbeASN         string
	ProbeIP          string
	DisableCollector bool
	SoftwareName     string
	SoftwareVersion  string
	Inputs           []string
	InputFilepaths   []string

	GeoIPCountryPath string
	GeoIPASNPath     string
	OutputPath       string
	CaBundlePath     string
	LogLevel         string
}

// NewNettest creates a new nettest instance
func NewNettest(name string) *Nettest {
	handleMap := make(map[string][]interface{})
	return &Nettest{
		Name:      name,
		handleMap: handleMap,
	}
}

// Nettest is a wrapper for running a particular nettest
type Nettest struct {
	Name           string
	Options        NettestOptions
	DisabledEvents []string

	handleMu  sync.Mutex
	handleMap map[string][]interface{}
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
// - measurement
// It is possible to register events with wildcards.
// For example On("status.*", ...) will fire on status.queued, status.started, ...
func (nt *Nettest) On(s string, v interface{}) error {
	nt.handleMu.Lock()
	defer nt.handleMu.Unlock()

	if reflect.ValueOf(v).Type().Kind() != reflect.Func {
		return errors.New("handler is not a function")
	}
	return nt.addHandler(s, v)
}

// EventValue are all the possible value keys
type EventValue struct {
	Idx     int64  `json:"idx"`
	JSONStr string `json:"json_str"`
	Failure string `json:"failure"`
	Input   string `json:"input"`

	LogLevel   string  `json:"log_level"`
	Percentage float64 `json:"percentage"`
	Message    string  `json:"message"`
	ProbeASN   string  `json:"probe_asn"`
	ProbeCC    string  `json:"probe_cc"`
	ProbeIP    string  `json:"probe_ip"`
	ReportID   string  `json:"report_id"`
}

// Event is an event fired from measurement_kit
type Event struct {

	// Is the key for the event. See On for the list of possible events.
	Key string `json:"key"`

	// Contains the value for the fired event
	Value EventValue `json:"value"`
}

// Run will run the test inside
func (nt *Nettest) Run() error {
	taskData, err := MakeTaskData(nt)
	if err != nil {
		return err
	}

	tdp, err := taskData.ToPointer()
	if err != nil {
		return err
	}

	task := C.mk_task_start(tdp)
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

		var e Event
		if err := json.Unmarshal([]byte(eventStr), &e); err != nil {
			return err
		}
		nt.fire(e.Key, e)
	}
	return nil
}
