package mk

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var handleMu sync.Mutex
var handleMap = make(map[string][]interface{})

func fire(s string, e Event) error {
	parts := strings.Split(s, ".")
	handles := make([]interface{}, 0)

	// Given a event key "foo.bar.tar", we are looking for either the exact handle
	// match "foo.bar.tar", "*" or foo.bar.*, foo.*
	handleNames := []string{
		s,
		"*",
	}
	for i := 1; i < len(parts); i++ {
		hn := fmt.Sprintf("%s.*", strings.Join(parts[0:i], "."))
		handleNames = append(handleNames, hn)
	}

	handleMu.Lock()

	for _, hn := range handleNames {
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
		f.Call(args)
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
