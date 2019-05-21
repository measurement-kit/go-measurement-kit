package mk

import "C"
import (
	"encoding/json"
	"unsafe"

	"github.com/measurement-kit/go-measurement-kit/internal/systemx"
)

// TaskData contains the measurement_kit task_data structure
type TaskData struct {
	DisabledEvents []string              `json:"disabled_events,omitempty"`
	Name           string                `json:"name"`
	LogLevel       string                `json:"log_level,omitempty"`
	OutputFilePath string                `json:"output_filepath,omitempty"`
	Options        measurementKitOptions `json:"options"`

	Inputs         []string          `json:"inputs,omitempty"`
	InputFilepaths []string          `json:"input_filepaths,omitempty"`
	Annotations    map[string]string `json:"annotations,omitempty"`
}

type measurementKitOptions struct {
	SaveRealProbeIP  bool    `json:"save_real_probe_ip,omitempty"`
	SaveRealProbeASN bool    `json:"save_real_probe_asn,omitempty"`
	SaveRealProbeCC  bool    `json:"save_real_probe_cc,omitempty"`
	NoBouncer        bool    `json:"no_bouncer,omitempty"`
	NoCollector      bool    `json:"no_collector,omitempty"`
	NoFileReport     bool    `json:"no_file_report,omitempty"`
	RandomizeInput   bool    `json:"randomize_input,omitempty"`
	SoftwareName     string  `json:"software_name,omitempty"`
	SoftwareVersion  string  `json:"software_version,omitempty"`
	ProbeCC          string  `json:"probe_cc,omitempty"`
	ProbeASN         string  `json:"probe_asn,omitempty"`
	ProbeIP          string  `json:"probe_ip,omitempty"`
	BouncerBaseURL   string  `json:"bouncer_base_url,omitempty"`
	CollectorBaseURL string  `json:"collector_base_url,omitempty"`
	GeoIPCountryPath string  `json:"geoip_country_path,omitempty"`
	GeoIPASNPath     string  `json:"geoip_asn_path,omitempty"`
	CaBundlePath     string  `json:"net/ca_bundle_path,omitempty"`
	MaxRuntime       float32 `json:"max_runtime,omitempty"`
}

// MakeTaskData returns a pointer to a TaskData structure
func MakeTaskData(nt *Nettest) (*TaskData, error) {
	logLevel := nt.Options.LogLevel
	if nt.Options.LogLevel == "" {
		logLevel = "INFO"
	}

	td := TaskData{
		Name:     nt.Name,
		LogLevel: logLevel,
		Options: measurementKitOptions{
			SaveRealProbeIP:  nt.Options.IncludeIP,
			SaveRealProbeASN: nt.Options.IncludeASN,
			SaveRealProbeCC:  nt.Options.IncludeCountry,
			NoBouncer:        nt.Options.DisableBouncer,
			NoCollector:      nt.Options.DisableCollector,
			NoFileReport:     nt.Options.DisableReportFile,
			RandomizeInput:   nt.Options.RandomizeInput,
			ProbeASN:         nt.Options.ProbeASN,
			ProbeCC:          nt.Options.ProbeCC,
			ProbeIP:          nt.Options.ProbeIP,
			CollectorBaseURL: nt.Options.CollectorBaseURL,
			BouncerBaseURL:   nt.Options.BouncerBaseURL,
			MaxRuntime:       nt.Options.MaxRuntime,

			SoftwareName:    nt.Options.SoftwareName,
			SoftwareVersion: nt.Options.SoftwareVersion,

			GeoIPCountryPath: nt.Options.GeoIPCountryPath,
			GeoIPASNPath:     nt.Options.GeoIPASNPath,
			CaBundlePath:     systemx.CABundlePath(nt.Options.CaBundlePath),
		},
		OutputFilePath: nt.Options.OutputPath,
		Annotations:    nt.Options.Annotations,
	}

	if len(nt.Options.Inputs) > 0 {
		td.Inputs = nt.Options.Inputs
	}
	if len(nt.Options.InputFilepaths) > 0 {
		td.InputFilepaths = nt.Options.InputFilepaths
	}

	if nt.DisabledEvents != nil {
		td.DisabledEvents = nt.DisabledEvents
	} else {
		td.DisabledEvents = make([]string, 0)
	}
	return &td, nil
}

// ToPointer converts a TaskData object to a C style pointer
func (td *TaskData) ToPointer() (*C.char, error) {
	tdBytes, err := json.Marshal(td)
	if err != nil {
		return nil, err
	}
	return (*C.char)(unsafe.Pointer(&tdBytes[0])), nil
}
