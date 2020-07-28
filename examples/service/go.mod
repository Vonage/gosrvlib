module github.com/nexmoinc/gosrvlib-sample-service

go 1.14

replace github.com/nexmoinc/gosrvlib => ../..

require (
	github.com/caarlos0/env/v6 v6.2.2
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/mock v1.4.3
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/nexmoinc/gosrvlib v0.0.0-00010101000000-000000000000
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/spf13/afero v1.3.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.15.0
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
)
