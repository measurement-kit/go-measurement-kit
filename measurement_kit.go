package mk

/*
#include <measurement_kit/ffi.h>
*/
// #cgo darwin,amd64 LDFLAGS: -L/usr/local/opt/openssl/lib
// #cgo darwin,amd64 LDFLAGS: /usr/local/lib/libmeasurement_kit.a
// #cgo darwin,amd64 LDFLAGS: /usr/local/opt/libevent/lib/libevent_core.a
// #cgo darwin,amd64 LDFLAGS: /usr/local/opt/libevent/lib/libevent_extra.a
// #cgo darwin,amd64 LDFLAGS: /usr/local/opt/libevent/lib/libevent_openssl.a
// #cgo darwin,amd64 LDFLAGS: /usr/local/opt/libevent/lib/libevent_pthreads.a
// #cgo darwin,amd64 LDFLAGS: /usr/local/opt/libmaxminddb/lib/libmaxminddb.a
// #cgo darwin,amd64 LDFLAGS: /usr/local/opt/openssl/lib/libssl.a
// #cgo darwin,amd64 LDFLAGS: /usr/local/opt/openssl/lib/libcrypto.a
// #cgo darwin,amd64 LDFLAGS: -lcurl
//
// #cgo windows LDFLAGS: -static
// #cgo windows,amd64 CFLAGS: -I/usr/local/opt/mingw-w64-measurement-kit/include/
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-measurement-kit/lib/libmeasurement_kit.a
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-libmaxminddb/lib/libmaxminddb.a
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-curl/lib/libcurl.a
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-libevent/lib/libevent_openssl.a
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-libressl/lib/libssl.a
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-libressl/lib/libcrypto.a
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-libevent/lib/libevent_core.a
// #cgo windows,amd64 LDFLAGS: /usr/local/opt/mingw-w64-libevent/lib/libevent_extra.a
// #cgo windows,amd64 LDFLAGS: -lws2_32
//
// #cgo linux,amd64 LDFLAGS: -static
// #cgo linux,amd64 LDFLAGS: /usr/local/lib/libmeasurement_kit.a
// #cgo linux,amd64 LDFLAGS: /usr/local/lib/libmaxminddb.a
// #cgo linux,amd64 LDFLAGS: /usr/local/lib/libcurl.a
// #cgo linux,amd64 LDFLAGS: /usr/lib/libevent_openssl.a
// #cgo linux,amd64 LDFLAGS: /usr/lib/libssl.a
// #cgo linux,amd64 LDFLAGS: /usr/lib/libcrypto.a
// #cgo linux,amd64 LDFLAGS: /usr/lib/libevent_core.a
// #cgo linux,amd64 LDFLAGS: /usr/lib/libevent_extra.a
// #cgo linux,amd64 LDFLAGS: /usr/lib/libevent_pthreads.a
// #cgo linux,amd64 LDFLAGS: /lib/libz.a
import "C"

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"
)

// NettestOptions are the options to be passed to a particular nettest
type NettestOptions struct {
	IncludeIP         bool
	IncludeASN        bool
	IncludeCountry    bool
	ProbeCC           string
	ProbeASN          string
	ProbeIP           string
	DisableBouncer    bool
	DisableCollector  bool
	DisableReportFile bool
	RandomizeInput    bool
	SoftwareName      string
	SoftwareVersion   string
	Inputs            []string
	InputFilepaths    []string
	BouncerBaseURL    string
	CollectorBaseURL  string
	MaxRuntime        float32

	GeoIPCountryPath string
	GeoIPASNPath     string
	OutputPath       string
	CaBundlePath     string
	LogLevel         string
	Annotations      map[string]string
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

	LogLevel         string  `json:"log_level"`
	Percentage       float64 `json:"percentage"`
	Message          string  `json:"message"`
	ProbeASN         string  `json:"probe_asn"`
	ProbeCC          string  `json:"probe_cc"`
	ProbeNetworkName string  `json:"probe_network_name"`
	ProbeIP          string  `json:"probe_ip"`
	ReportID         string  `json:"report_id"`
	DownloadedKB     float64 `json:"downloaded_kb"`
	UploadedKB       float64 `json:"uploaded_kb"`
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
