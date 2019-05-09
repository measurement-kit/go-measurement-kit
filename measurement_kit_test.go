package mk

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func readFirstLine(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return []byte(line), nil
}

func TestTelegram(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "go-mk-testing")
	if err != nil {
		panic(err)
	}
	outputPath := path.Join(tmpdir, "telegram-result.jsonl")

	nt := NewNettest("Telegram")
	nt.Options = NettestOptions{
		DisableCollector: true,
		GeoIPASNPath:     "testdata/asn.mmdb",
		GeoIPCountryPath: "testdata/country.mmdb",
		CaBundlePath:     "testdata/ca-bundle.pem",
		OutputPath:       outputPath,
	}
	nt.On("log", func(e Event) {
		loglevel := e.Value.LogLevel
		msg := e.Value.Message
		if loglevel == "ERROR" || loglevel == "WARNING" {
			t.Errorf("%s level log message '%s'", loglevel, msg)
		}
	})
	nt.Run()

	line, err := readFirstLine(outputPath)
	if err != nil {
		t.Errorf("failed to open json file: %v", err)
	}
	var msmt map[string]interface{}
	err = json.Unmarshal(line, &msmt)
	if err != nil {
		t.Errorf("failed to parse json file: %v", err)
	}
	_, exist := msmt["test_runtime"]
	if !exist {
		t.Errorf("did not find the test_runtime key")
	}
}
