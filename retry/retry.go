package retry

import (
	"time"

	"github.com/cenkalti/backoff"
)

// Tryer is an operation which is likely to fail
type Tryer interface {
	Try() error
}

type ticker interface {
	Stop()
}

// Retry represents a
type Retry struct {
	Errors <-chan error
	errors chan<- error
	ticker ticker
	timeCh <-chan time.Time
}

// NewRetry with an exponential backoff. Errors are reported through a
// synchronous errors channel and _must_ be consumed in a separate goroutine.
// The Retry will be blocked until the previous error is consumed
func NewRetry() *Retry {
	errors := make(chan error)
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	return &Retry{
		Errors: errors,
		errors: errors,
		ticker: ticker,
		timeCh: ticker.C,
	}
}

// Execute tries initially, then if there is an error it uses a backoff ticker
// retrying a maximum of n attempts
// errors are passed into the return error channel which is closed after success
// or the final retry
func (retry *Retry) Execute(tryer Tryer) {
	if err := tryer.Try(); err == nil {
		retry.Stop()
	} else {
		retry.errors <- err
		go retry.retry(tryer)
	}
}

// Stop retrying the Tryer
func (retry *Retry) Stop() {
	close(retry.errors)
	retry.ticker.Stop()
}

func (retry *Retry) retry(tryer Tryer) {
	for {
		if _, stillOpen := <-retry.timeCh; !stillOpen {
			break
		}

		if err := tryer.Try(); err == nil {
			retry.Stop()
		} else {
			retry.errors <- err
		}
	}
}
