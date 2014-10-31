package retry

import (
	"fmt"
	"testing"
	"time"
)

type testTryer struct {
	tries int
}

func (tryer *testTryer) Try() error {
	tryer.tries++
	return fmt.Errorf("FAILED")
}

type testTicker struct {
	stopped bool
	c       chan time.Time
}

func (testTicker *testTicker) Stop() {
	testTicker.stopped = true
	close(testTicker.c)
}

func TestRetryWithFailingTryer(t *testing.T) {
	retry := NewRetry()
	testTicker := &testTicker{
		c: make(chan time.Time),
	}
	retry.ticker = testTicker
	retry.timeCh = testTicker.c

	errors := make(chan error, 128)

	// stash errors in a slice
	go func() {
		for e := range retry.Errors {
			errors <- e
		}
	}()

	tryer := &testTryer{}
	ok := retry.Execute(tryer)

	if ok {
		t.Fatalf("expected the failed try to return false")
	}

	if tryer.tries != 1 {
		t.Fatalf("expected tries to be 1 but was: %d", tryer.tries)
	}

	select {
	case <-errors:
	default:
		t.Fatalf("expected an error non were passed")
	}

	// test retries
	testTicker.c <- time.Now()
	timeout := time.After(1 * time.Millisecond)

	select {
	case <-errors:
		if tryer.tries != 2 {
			t.Fatalf("expected tries to be 2 but was: %d", tryer.tries)
		}
	case <-timeout:
		t.Fatalf("timed-ocut waiting for an error")
	}

	testTicker.c <- time.Now()
	timeout = time.After(1 * time.Millisecond)

	select {
	case <-errors:
		if tryer.tries != 3 {
			t.Fatalf("expected tries to be 3 but was: %d", tryer.tries)
		}
	case <-timeout:
		t.Fatalf("timed-out waiting for an error")
	}

	retry.Stop()

	if _, stillOpen := <-retry.Errors; stillOpen {
		t.Fatalf("Error channel should have been closed when the retry is stopped")
	}
}

type succeedingTestTryer struct {
	tries int
}

func (tryer *succeedingTestTryer) Try() error {
	tryer.tries++
	return nil
}

func TestRetryWithImmediatelySucceedingTryer(t *testing.T) {
	retry := NewRetry()
	testTicker := &testTicker{
		c: make(chan time.Time),
	}
	retry.ticker = testTicker
	retry.timeCh = testTicker.c

	errors := make(chan error, 128)

	// stash errors in a slice
	go func() {
		for e := range retry.Errors {
			errors <- e
		}
	}()

	tryer := &succeedingTestTryer{}
	ok := retry.Execute(tryer)

	if !ok {
		t.Fatalf("expected the successful try to return true")
	}

	if tryer.tries != 1 {
		t.Fatalf("expected tries to be 1 but was: %d", tryer.tries)
	}

	select {
	case <-errors:
		t.Fatalf("expected no errors but an error was passed")
	default:
	}

	if _, stillOpen := <-retry.Errors; stillOpen {
		t.Fatalf("Error channel should have been closed when the retry is stopped")
	}
}

type eventuallySucceedingTestTryer struct {
	tries int
}

func (tryer *eventuallySucceedingTestTryer) Try() error {
	tryer.tries++

	if tryer.tries >= 3 {
		return nil
	}

	return fmt.Errorf("failed")
}

func TestRetryWithEventuallySucceedingTryer(t *testing.T) {
	retry := NewRetry()
	testTicker := &testTicker{
		c: make(chan time.Time),
	}
	retry.ticker = testTicker
	retry.timeCh = testTicker.c

	errors := make(chan error, 128)

	// stash errors in a slice
	go func() {
		for e := range retry.Errors {
			errors <- e
		}
		close(errors)
	}()

	tryer := &eventuallySucceedingTestTryer{}
	retry.Execute(tryer)

	testTicker.c <- time.Now()
	testTicker.c <- time.Now()

	errorSlice := make([]error, 0, 3)
	for e := range errors {
		errorSlice = append(errorSlice, e)
	}

	if len(errorSlice) != 2 {
		t.Fatalf("expected there to be 2 errors but there were: %d", len(errorSlice))
	}

	if tryer.tries != 3 {
		t.Fatalf("expected tries to be 3 but was: %d", tryer.tries)
	}

	if _, stillOpen := <-retry.Errors; stillOpen {
		t.Fatalf("Error channel should have been closed when the retry is stopped")
	}
}

/*
  panic if Execute is called more than once
*/
