//go:build gotools

// Package gotools lists external build and test tools.
// These tools will appear in the `go.mod` file, but will not be a part of the build.
// They will be also excluded from the binaries as the "gotools" tag is not used.
package gotools

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/jstemmer/go-junit-report"
	_ "github.com/rakyll/gotest"
)
