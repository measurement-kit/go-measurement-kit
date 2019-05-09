// +build !darwin,!linux,!windows

// Package systemx contains system extensions.
package systemx

// CABundlePath returns the CA bundle path to use. Since we don't know the
// system in which we are, we cannot guess a default CA bundle. So, the user
// must really provide a valid CA bundle path in this case.
func CABundlePath(path string) string {
	return path
}
