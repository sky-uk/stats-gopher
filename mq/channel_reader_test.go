package mq

import (
	"fmt"
	"testing"
	"time"
)

func TestChannelReader(t *testing.T) {
	maxSize := 10
	ch := make(chan interface{}, 128)
	r := NewChannelReader(ch, 10)

	var chunk []interface{}
	go r.Read()

	// initially the output should block because there is no data
	select {
	case <-r.Out:
		t.Fatalf("Read() did not block when no data was available")
	default:
		break
	}

	// now put some data in
	ch <- 1
	ch <- 2
	ch <- 3
	ch <- 4
	ch <- 5

	// wait for the read goroutine
	timeout := time.After(1 * time.Millisecond)

	// some data should be available
	select {
	case chunk = <-r.Out:
		// check the right nuber of elements are available
		if len(chunk) != 5 {
			t.Fatal("the reader did not output a 5 element chunk of all available data")
		}

		// check the input data and out data match
		for i, e := range chunk {
			actual := fmt.Sprintf("%v", e)
			expected := fmt.Sprintf("%v", i+1)
			if actual != expected {
				t.Fatalf("element %d in the chunk was %s when it should have been %s", i, actual, expected)
			}
		}
	case <-timeout:
		t.Fatalf("Read() did return available data")
	}

	// the input should have been drained, nothing is available
	select {
	case <-r.Out:
		t.Fatalf("Read() did not block when no data was available")
	default:
		break
	}

	// put some more data in
	ch <- 1
	ch <- 1
	ch <- 1

	// wait for the read goroutine
	timeout = time.After(1 * time.Millisecond)

	// more data should be available
	select {
	case chunk = <-r.Out:
		if len(chunk) != 3 {
			t.Fatalf("the reader did not output a 3 element chunk of all available data: %v", chunk)
		}
	case <-timeout:
		t.Fatalf("Read() did return available data")
	}

	// overfill the input relative to the max size
	overfill := 5
	for i := -overfill; i < maxSize; i++ {
		ch <- i
	}

	if len(<-r.Out) != maxSize {
		t.Fatal("The first chunk of output should have been the maximum size")
	}

	if len(<-r.Out) != overfill {
		t.Fatal("The second chunk of output should have been the size of the over-fill")
	}
}
