package sqlconn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_config_validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     *config
		wantErr bool
	}{
		{
			name:    "fail with empty driver",
			cfg:     defaultConfig("", "user:pass@tcp(127.0.0.1:1234)/testdb"),
			wantErr: true,
		},
		{
			name:    "fail with empty DSN",
			cfg:     defaultConfig("sqldb", ""),
			wantErr: true,
		},
		{
			name: "fail with invalid max retry value",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.connectMaxRetry = 0
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "fail with invalid retry interval",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.connectRetryInterval = 0
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "fail with invalid connect function",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.connectFunc = nil
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "fail with invalid check connection function",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.checkConnectionFunc = nil
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "fail with invalid sql open function",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.sqlOpenFunc = nil
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "fail with invalid max idle",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.connMaxIdle = 0
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "fail with invalid quoteIDFunc function",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.quoteIDFunc = nil
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "fail with invalid quoteValueFunc function",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				cfg.quoteValueFunc = nil
				return cfg
			}(),
			wantErr: true,
		},
		{
			name: "succeed with no errors",
			cfg: func() *config {
				cfg := defaultConfig("sqldb", "user:pass@tcp(127.0.0.1:1234)/testdb")
				return cfg
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.cfg.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_defaultConfig(t *testing.T) {
	t.Parallel()

	cfg := defaultConfig("test_driver", "test_dsn")
	require.NotNil(t, cfg)
	require.Equal(t, "test_driver", cfg.driver)
	require.Equal(t, "test_dsn", cfg.dsn)
	require.NotNil(t, cfg.quoteIDFunc)
	require.NotNil(t, cfg.quoteValueFunc)
	require.NotNil(t, cfg.connectFunc)
	require.NotNil(t, cfg.checkConnectionFunc)
	require.NotNil(t, cfg.sqlOpenFunc)
	require.NotEqual(t, 0, cfg.connectMaxRetry)
	require.NotEqual(t, 0, cfg.connectRetryInterval)
	require.Equal(t, defaultConnMaxIdle, cfg.connMaxIdle)
	require.Equal(t, defaultConnMaxLifetime, cfg.connMaxLifetime)
	require.Equal(t, defaultConnMaxOpen, cfg.connMaxOpen)
}

func Test_defaultQuoteID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "empty string",
			id:   "",
			want: "",
		},
		{
			name: "single value",
			id:   "test1",
			want: "`test1`",
		},
		{
			name: "two parts",
			id:   "parent.child",
			want: "`parent`.`child`",
		},
		{
			name: "multiple parts",
			id:   "parent.child.name",
			want: "`parent`.`child`.`name`",
		},
		{
			name: "multiple parts with space",
			id:   "one two.three four",
			want: "`one two`.`three four`",
		},
		{
			name: "escape backtick",
			id:   "test`4",
			want: "`test``4`",
		},
		{
			name: "escape multiple backtick",
			id:   "test```4",
			want: "`test``````4`",
		},
		{
			name: "special characters",
			id:   "test" + string([]byte{'\'', '"', '`', 0, '\n', '\r', '\\', '\032'}),
			want: "`test'\"``\\0\\n\\r\\\\\\Z`",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := defaultQuoteID(tt.id)

			require.Equal(t, tt.want, got, "QuoteID() got = %v, want %v", got, tt.want)
		})
	}
}

func Test_defaultQuoteValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  string
		want string
	}{
		{
			name: "empty string",
			val:  "",
			want: "",
		},
		{
			name: "simple value",
			val:  "test1",
			want: "'test1'",
		},
		{
			name: "value with spaces",
			val:  "one two three",
			want: "'one two three'",
		},
		{
			name: "escape single quote",
			val:  "test'2",
			want: "'test''2'",
		},
		{
			name: "escape multiple quotes",
			val:  "test'''2",
			want: "'test''''''2'",
		},
		{
			name: "special characters",
			val:  "test" + string([]byte{'\'', '"', '`', 0, '\n', '\r', '\\', '\032'}),
			want: "'test''\"`\\0\\n\\r\\\\\\Z'",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := defaultQuoteValue(tt.val)

			require.Equal(t, tt.want, got, "QuoteID() got = %v, want %v", got, tt.want)
		})
	}
}
