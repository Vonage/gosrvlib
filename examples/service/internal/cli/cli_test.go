package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/nexmoinc/gosrvlib/pkg/testutil"
	"github.com/stretchr/testify/require"
)

// // FIXME: refactor
//
// func TestCliBadParamError(t *testing.T) {
// 	for _, param := range badParamCases {
// 		os.Args = []string{ProgramName, param}
// 		cmd, err := cli()
// 		if err != nil {
// 			t.Errorf("Unexpected error: %v", err)
// 			return
// 		}
// 		if cmdtype := reflect.TypeOf(cmd).String(); cmdtype != "*cobra.Command" {
// 			t.Errorf("The expected type is '*cobra.Command', found: '%s'", cmdtype)
// 			return
// 		}
//
// 		old := os.Stderr // keep backup of the real stdout
// 		defer func() { os.Stderr = old }()
// 		os.Stderr = nil
//
// 		// execute the main function
// 		if err := cmd.Execute(); err == nil {
// 			t.Errorf("An error was expected")
// 		}
// 	}
// }
//
// func TestWrongParamError(t *testing.T) {
// 	os.Args = []string{ProgramName, "--unknown"}
// 	_, err := cli()
// 	if err == nil {
// 		t.Errorf("An error was expected")
// 		return
// 	}
// 	os.Args = []string{ProgramName, ""}
// }

// var badParamCases = []string{
// 	"=",
// 	"",
// 	"=../resources/test/etc/mysql_err",
// 	"--configDir=../resources/test/etc/db_err",
// 	"--configDir=../resources/test/etc/tls_err",
// }
//

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		osArgs     []string
		wantOutput func(t *testing.T, out string)
		wantErr    bool
	}{
		{
			name:       "call version subcommand",
			osArgs:     []string{AppName, "version"},
			wantErr:    false,
			wantOutput: matchTestVersion,
		},
		{
			name:       "",
			osArgs:     []string{AppName, "--unknown"},
			wantErr:    true,
			wantOutput: matchErrorOutput,
		},
		{
			name:       "",
			osArgs:     []string{AppName, "--logFormat"},
			wantErr:    true,
			wantOutput: matchErrorOutput,
		},
		{
			name:       "",
			osArgs:     []string{AppName, "--logLevel"},
			wantErr:    true,
			wantOutput: matchErrorOutput,
		},
		{
			name:    "",
			osArgs:  []string{AppName, "--logLevel", "INVALID"},
			wantErr: true,
		},
		{
			name:       "",
			osArgs:     []string{AppName, "--configDir"},
			wantErr:    true,
			wantOutput: matchErrorOutput,
		},
		{
			name:    "",
			osArgs:  []string{AppName, "-c", "../../resources/test/etc/srvxmplname/", "--logFormat", "invalid"},
			wantErr: true,
		},
		{
			name:    "",
			osArgs:  []string{AppName, "-c", "../../resources/test/etc/srvxmplname/", "--logLevel", "invalid"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cmd, _ := New("0.0.0-test", "0")
			require.NotNil(t, cmd)

			oldOsArgs := os.Args
			defer func() { os.Args = oldOsArgs }()
			os.Args = tt.osArgs

			// execute the main function
			var err error
			out := testutil.CaptureOutput(func() {
				err = cmd.Execute()
			})

			if tt.wantOutput != nil {
				t.Log(out)
				tt.wantOutput(t, out)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func matchErrorOutput(t *testing.T, out string) {
	if strings.HasPrefix(out, "Error:") {
		return
	}
	t.Errorf("An error message was expected")
}

func matchTestVersion(t *testing.T, out string) {
	if strings.HasPrefix(out, "0.0.0-test") {
		return
	}
	t.Errorf("A version number was expected")
}
