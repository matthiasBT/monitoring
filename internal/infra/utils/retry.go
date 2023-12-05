// Package utils provides utility functions and types for error handling and
// retry logic, particularly with PostgreSQL connections.
package utils

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

// Retrier is a struct that encapsulates retry logic parameters and a logger.
// It is used to perform a function multiple times with delays, logging each attempt.
type Retrier struct {
	Logger           logging.ILogger // Logger for logging retry attempts and errors
	Attempts         int             // Number of retry attempts
	IntervalFirst    time.Duration   // Initial interval between retries
	IntervalIncrease time.Duration   // Incremental increase of the interval per retry
}

// ErrRetryAborted is the error returned when a retry operation is aborted due to context cancellation.
var ErrRetryAborted = errors.New("retry aborted by context")

// retriableErrorsPostgreSQL defines a set of PostgreSQL error codes that are considered retriable.
var retriableErrorsPostgreSQL = map[string]bool{
	pgerrcode.ConnectionException:                           true,
	pgerrcode.ConnectionDoesNotExist:                        true,
	pgerrcode.ConnectionFailure:                             true,
	pgerrcode.SQLClientUnableToEstablishSQLConnection:       true,
	pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection: true,
	pgerrcode.TransactionResolutionUnknown:                  true,
	pgerrcode.ProtocolViolation:                             true,
}

// RetryChecked attempts to execute a function multiple times (based on Retrier settings)
// until it succeeds or a non-retriable error is encountered.
// It takes a context, a function to execute, and a function to check if an error is retriable.
func (r *Retrier) RetryChecked(ctx context.Context, f func() (any, error), checkError func(error) bool) (any, error) {
	var result any
	var err error
	var errChain []error

	interval := r.IntervalFirst
	for i := 0; i <= r.Attempts; i++ {
		if i > 1 {
			r.Logger.Infof("Starting retry %d of %d\n", i, r.Attempts)
		}
		result, err = f()
		if err != nil && checkError(err) && i != r.Attempts {
			r.Logger.Infof("Retriable error: %s. Repeat after: %s\n", err.Error(), interval)
			errChain = append(errChain, err)
			if sleepErr := sleepInContext(ctx, interval); sleepErr != nil {
				errChain = append(errChain, sleepErr)
				return nil, errors.Join(errChain...)
			}
			interval += r.IntervalIncrease
			continue
		} else if err != nil {
			r.Logger.Errorf("Non-retriable error: %s\n", err.Error())
			errChain = append(errChain, err)
			return result, errors.Join(errChain...)
		}
		r.Logger.Infoln("Success")
		return result, nil
	}
	if err := errors.Join(errChain...); err != nil {
		r.Logger.Errorf("All attempts failed\n")
		return nil, err
	}
	r.Logger.Infoln("Success")
	return result, nil
}

// CheckConnectionError checks if an error is a retriable PostgreSQL error or a network error.
// It is used as an error checker function for RetryChecked.
func CheckConnectionError(err error) bool {
	var pgErr *pgconn.PgError
	var netErr *net.OpError
	return errors.As(err, &pgErr) && retriableErrorsPostgreSQL[pgErr.Code] || errors.As(err, &netErr)
}

// sleepInContext pauses execution for the specified duration or until the context is cancelled.
// Returns ErrRetryAborted if the context is cancelled before the timeout.
func sleepInContext(ctx context.Context, timeout time.Duration) error {
	select {
	case <-ctx.Done():
		return ErrRetryAborted
	case <-time.After(timeout):
		return nil
	}
}
