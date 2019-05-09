// Package systemx contains system extensions.
package systemx

import (
	"os"
)

// CABundlePath returns the CA bundle path to use. If the provided path
// is not empty, that is used. Otherwise, we'll see whether this path is
// provided as part of the MEASUREMENT_KIT_CA_BUNDLE_PATH env variable.
func CABundlePath(path string) string {
	if path == "" {
		return os.Getenv("MEASUREMENT_KIT_CA_BUNDLE_PATH")
	}
	return path
}
