package mk

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type taskData struct {
	DisabledEvents []string       `json:"disabled_events"`
	Type           string         `json:"type"`
	Verbosity      string         `json:"verbosity"`
	Options        NettestOptions `json:"options"`
}

var handleMu sync.Mutex
var handleMap = make(map[string][]interface{})

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
		if len(values) > 0 {
			return values[0].Interface().(error)
		}
		return nil
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
