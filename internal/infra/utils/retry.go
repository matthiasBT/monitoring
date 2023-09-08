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

type Retrier struct {
	Attempts         int
	IntervalFirst    time.Duration
	IntervalIncrease time.Duration
	Logger           logging.ILogger
}

var ErrRetryAborted = errors.New("retry aborted by context")

var retriableErrorsPostgreSQL = map[string]bool{
	pgerrcode.ConnectionException:                           true,
	pgerrcode.ConnectionDoesNotExist:                        true,
	pgerrcode.ConnectionFailure:                             true,
	pgerrcode.SQLClientUnableToEstablishSQLConnection:       true,
	pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection: true,
	pgerrcode.TransactionResolutionUnknown:                  true,
	pgerrcode.ProtocolViolation:                             true,
}

func (r *Retrier) RetryChecked(ctx context.Context, f func() (any, error), checkError func(error) bool) (any, error) {
	var result any
	var errChain []error

	interval := r.IntervalFirst
	for i := 1; i <= r.Attempts; i++ {
		if i > 1 { // prevent spamming
			r.Logger.Infof("Starting attempt %d of %d\n", i, r.Attempts)
		}
		result, err := f()
		if err != nil && checkError(err) {
			r.Logger.Infof("Retriable error: %s\n", err.Error())
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

func CheckConnectionError(err error) bool {
	var pgErr *pgconn.PgError
	var netErr *net.OpError
	return errors.As(err, &pgErr) && retriableErrorsPostgreSQL[pgErr.Code] || errors.As(err, &netErr)
}

func sleepInContext(ctx context.Context, timeout time.Duration) error {
	select {
	case <-ctx.Done():
		return ErrRetryAborted
	case <-time.After(timeout):
		return nil
	}
}
