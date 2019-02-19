// Package ldflags contains variables set on
// compiling the binary
//   Authors: Ringo Hoffmann
package ldflags

var (
	// AppVersion is the tag version
	AppVersion = ""
	// AppCommit is the last commit hash
	AppCommit = ""
	// GoVersion is the version of the go compiler
	GoVersion = ""
	// Release is "TRUE" if this is a release build
	Release = ""
)
