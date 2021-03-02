package sqlconn

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	defaultConnMaxIdle          = 2                // Maximum number of idle connections (0 = unlimited)
	defaultConnMaxLifetime      = time.Duration(0) // Maximum connection lifetime in seconds (0 = unlimited reuse)
	defaultConnMaxOpen          = 2                // Maximum number of open connections (0 = unlimited connections)
	defaultConnectMaxRetry      = 1                // Number of maximum connection retries
	defaultConnectRetryInterval = 3 * time.Second  // Connection retry time in seconds
)

func defaultConfig(driver, dsn string) *config {
	return &config{
		quoteIDFunc:          defaultQuoteID,
		quoteValueFunc:       defaultQuoteValue,
		checkConnectionFunc:  checkConnection,
		sqlOpenFunc:          sql.Open,
		connectFunc:          connectWithBackoff,
		connMaxLifetime:      defaultConnMaxLifetime,
		connectRetryInterval: defaultConnectRetryInterval,
		connectMaxRetry:      defaultConnectMaxRetry,
		connMaxIdle:          defaultConnMaxIdle,
		connMaxOpen:          defaultConnMaxOpen,
		driver:               driver,
		dsn:                  dsn,
	}
}

type config struct {
	quoteIDFunc          SQLQuoteFunc
	quoteValueFunc       SQLQuoteFunc
	checkConnectionFunc  CheckConnectionFunc
	sqlOpenFunc          SQLOpenFunc
	connectFunc          ConnectFunc
	connMaxLifetime      time.Duration
	connectRetryInterval time.Duration
	connectMaxRetry      int
	connMaxIdle          int
	connMaxOpen          int
	driver               string
	dsn                  string
}

// nolint:gocyclo
func (c *config) validate() error {
	if strings.TrimSpace(c.driver) == "" {
		return fmt.Errorf("database driver must be set")
	}

	if strings.TrimSpace(c.dsn) == "" {
		return fmt.Errorf("database DSN must be set")
	}

	if c.connectMaxRetry < 1 {
		return fmt.Errorf("database connect max retry must be greater than 0")
	}

	if c.connectRetryInterval < 1*time.Second {
		return fmt.Errorf("database connect retry interval must be greater than 1s")
	}

	if c.connectFunc == nil {
		return fmt.Errorf("database connect function must be set")
	}

	if c.checkConnectionFunc == nil {
		return fmt.Errorf("check connection function must be set")
	}

	if c.sqlOpenFunc == nil {
		return fmt.Errorf("sql open function must be set")
	}

	if c.connMaxIdle < 1 {
		return fmt.Errorf("database pool max idle connections must be greater than 0")
	}

	if c.quoteIDFunc == nil {
		return fmt.Errorf("the QuoteID function must be set")
	}

	if c.quoteValueFunc == nil {
		return fmt.Errorf("the QuoteValue function must be set")
	}

	return nil
}

// defaultQuoteID is the QuoteID default function for mysql-like databases.
func defaultQuoteID(s string) string {
	if s == "" {
		return s
	}

	parts := strings.Split(s, ".")

	for k, v := range parts {
		parts[k] = "`" + strings.ReplaceAll(escape(v), "`", "``") + "`"
	}

	return strings.Join(parts, ".")
}

// defaultQuoteValue is the QuoteValue default function for mysql-like databases.
func defaultQuoteValue(s string) string {
	if s == "" {
		return s
	}

	return "'" + strings.ReplaceAll(escape(s), "'", "''") + "'"
}

func escape(s string) string {
	dest := make([]byte, 0, 2*len(s))

	var escape byte

	for i := 0; i < len(s); i++ {
		c := s[i]

		escape = 0

		switch c {
		case 0:
			escape = '0'
		case '\n':
			escape = 'n'
		case '\r':
			escape = 'r'
		case '\\':
			escape = '\\'
		case '\032':
			escape = 'Z'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}
