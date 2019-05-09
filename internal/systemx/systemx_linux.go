// Package systemx contains system extensions.
package systemx

import (
	"os"
)

// CABundlePath returns the CA bundle path to use. If the provided path
// is not empty, that is used. Otherwise, we'll use Linux defaults.
func CABundlePath(path string) string {
	if path == "" {
		// Copied from CURL's configure script
		availablePaths := []string{
			"/etc/ssl/certs/ca-certificates.crt",
			"/etc/pki/tls/certs/ca-bundle.crt",
			"/usr/share/ssl/certs/ca-bundle.crt",
			"/usr/local/share/certs/ca-root.crt",
			"/etc/ssl/cert.pem",
			"/usr/local/etc/openssl/cert.pem",
		}
		for _, possiblePath := range availablePaths {
			if _, err := os.Stat(possiblePath); err == nil {
				return possiblePath
			}
		}
		// FALLTHROUGH
	}
	return path
}
