package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/nexmoinc/gosrvlib/pkg/errtrace"
	"github.com/nexmoinc/gosrvlib/pkg/logging"
)

// ReleaseFunc is an alias for a release lock function.
type ReleaseFunc func() error

var (
	// ErrAcquireLockTimeout is an error when the lock is not acquired within the timeout.
	ErrAcquireLockTimeout = errors.New("acquire lock timeout")

	// ErrAcquireLockError is an error when the lock is not acquired.
	ErrAcquireLockError = errors.New("acquire lock error")
)

const (
	sqlGetLock     = "SELECT COALESCE(GET_LOCK(?, ?), 2)"
	sqlReleaseLock = "DO RELEASE_LOCK(?)"
)

type Locker struct {
	db *sql.DB
}

func New(db *sql.DB) *Locker {
	return &Locker{db: db}
}

// AcquireLock attempts to acquire a database lock.
func (l *Locker) AcquireLock(ctx context.Context, key string, timeout time.Duration) (ReleaseFunc, error) {
	conn, err := l.db.Conn(ctx)

	if err != nil {
		return nil, errtrace.Trace(err)
	}

	row := conn.QueryRowContext(ctx, sqlGetLock, key, int(timeout.Seconds()))

	var res int
	err = row.Scan(&res)

	if err != nil {
		return nil, errtrace.Trace(err)
	}

	releaseFunc := func() error {
		defer logging.Close(ctx, conn, "error closing lock connection")

		// background context used to ensure that release lock is always executed
		_, err := conn.ExecContext(context.Background(), sqlReleaseLock, key)
		if err != nil {
			return errtrace.Trace(err)
		}

		return nil
	}

	switch res {
	case 0:
		return nil, ErrAcquireLockTimeout
	case 1:
		return releaseFunc, nil
	default:
		return nil, ErrAcquireLockError
	}
}
