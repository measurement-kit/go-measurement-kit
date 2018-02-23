package mk

import "encoding/json"

// Entry represents a measurement entry
type Entry struct {
	//InputHashes []interface{} `json:"input_hashes"`
	Annotations          map[string]interface{} `json:"annotations"`
	DataFormatVersion    string                 `json:"data_format_version"`
	ID                   string                 `json:"id"`
	Input                string                 `json:"input"`
	MeasurementStartTime string                 `json:"measurement_start_time"`
	Options              []string               `json:"options"`
	ProbeAsn             string                 `json:"probe_asn"`
	ProbeCc              string                 `json:"probe_cc"`
	ProbeCity            string                 `json:"probe_city"`
	ProbeIP              string                 `json:"probe_ip"`
	ReportID             string                 `json:"report_id"`
	SoftwareName         string                 `json:"software_name"`
	SoftwareVersion      string                 `json:"software_version"`
	TestHelpers          map[string]interface{} `json:"test_helpers"`
	TestKeys             map[string]interface{} `json:"test_keys"`
	TestName             string                 `json:"test_name"`
	TestRuntime          float32                `json:"test_runtime"`
	TestStartTime        string                 `json:"test_start_time"`
	TestVersion          string                 `json:"test_version"`
}

// ParseEntry into a Entry struct
func ParseEntry(s string, e *Entry) error {
	return json.Unmarshal([]byte(s), &e)
}
