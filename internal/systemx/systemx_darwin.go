// Package systemx contains system extensions.
package systemx

// CABundlePath returns the CA bundle path to use. If the provided path
// is not empty, that is used. Otherwise, we'll use macOS default.
func CABundlePath(path string) string {
	if path == "" {
		return "/etc/ssl/cert.pem" // default CA bundle path on macOS
	}
	return path
}
