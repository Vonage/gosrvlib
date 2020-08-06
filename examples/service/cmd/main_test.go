package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/nexmoinc/gosrvlib-sample-service/internal/cli"
	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestProgramVersion(t *testing.T) {
	os.Args = []string{cli.AppName, "version"}
	out, err := testutil.CaptureOutput(func() {
		main()
	})
	require.NoError(t, err)

	match, err := regexp.MatchString("^[\\d]+\\.[\\d]+\\.[\\d]+[\\s]*$", out)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !match {
		t.Errorf("The expected version has not been returned")
	}
}
