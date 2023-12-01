package utils

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

var ErrSome = errors.New("some error")

type UnreliableCallCounter struct {
	called         int
	stopFailuresAt int
}

func (cc *UnreliableCallCounter) do() (any, error) {
	cc.called += 1
	if cc.called == cc.stopFailuresAt {
		time.Sleep(100 * time.Millisecond)
		return 100500, nil
	}
	return nil, ErrSome
}

func TestRetrier_RetryChecked(t *testing.T) {
	type args struct {
		timeout    time.Duration
		f          func() (any, error)
		checkError func(error) bool
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr error
	}{
		{
			name: "success_first_try",
			args: args{
				f:          func() (any, error) { return 100500, nil },
				checkError: nil,
			},
			want:    100500,
			wantErr: nil,
		},
		{
			name: "success_second_try",
			args: args{
				f:          (&UnreliableCallCounter{stopFailuresAt: 2}).do,
				checkError: func(err error) bool { return errors.Is(err, ErrSome) },
			},
			want:    100500,
			wantErr: nil,
		},
		{
			name: "failure_on_unchecked_error",
			args: args{
				f:          (&UnreliableCallCounter{stopFailuresAt: 2}).do,
				checkError: func(err error) bool { return !errors.Is(err, ErrSome) },
			},
			want:    nil,
			wantErr: ErrSome,
		},
		{
			name: "number_of_attempts_exceeded",
			args: args{
				f:          (&UnreliableCallCounter{stopFailuresAt: 5}).do,
				checkError: func(err error) bool { return errors.Is(err, ErrSome) },
			},
			want:    nil,
			wantErr: ErrSome,
		},
		{
			name: "timeout",
			args: args{
				f:          (&UnreliableCallCounter{stopFailuresAt: 5}).do,
				checkError: func(err error) bool { return errors.Is(err, ErrSome) },
				timeout:    1 * time.Millisecond,
			},
			want:    nil,
			wantErr: ErrRetryAborted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Retrier{
				Attempts:         2,
				IntervalFirst:    100 * time.Millisecond,
				IntervalIncrease: 100 * time.Millisecond,
				Logger:           logging.SetupLogger(),
			}
			ctx := context.Background()
			if tt.args.timeout != 0 {
				ctxt, cancel := context.WithTimeout(ctx, tt.args.timeout)
				ctx = ctxt
				defer cancel()
			}
			got, err := r.RetryChecked(ctx, tt.args.f, tt.args.checkError)
			if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("RetryChecked() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetryChecked() got = %v, want %v", got, tt.want)
			}
		})
	}
}
