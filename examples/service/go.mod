module github.com/nexmoinc/gosrvlib-sample-service

go 1.14

replace github.com/nexmoinc/gosrvlib => ../..

require (
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/mock v1.4.3
	github.com/nexmoinc/gosrvlib v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.0.0
	go.uber.org/zap v1.15.0
)
