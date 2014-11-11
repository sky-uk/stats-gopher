package presence

import (
	"runtime"
	"testing"
	"time"
)

func TestMonitorPoolEndEvents(t *testing.T) {
	mp := NewMonitorPool(map[string]time.Duration{
		"heartbeat":   time.Millisecond * 3,
		"user-active": time.Millisecond * 2,
	})

	// test the session lookup
	s := mp.session("key1")

	if mp.session("key1") == nil {
		t.Fatalf("no session was returned")
	}

	if s.monitor("heartbeat", time.Millisecond*3).timeout != time.Millisecond*3 {
		t.Fatalf("session monitor timeout for 'heartbeat' was not set")
	}

	if s.monitor("user-active", time.Millisecond*2).timeout != time.Millisecond*2 {
		t.Fatalf("session monitor timeout for 'user' was not set")
	}

	if mp.session("key1") != s {
		t.Fatalf("existing session was not returned")
	}

	timeout := &Timeout{
		Start:            time.Now(),
		LastNotification: time.Now(),
		End:              time.Now(),
	}

	s.end("heartbeat", timeout)

	runtime.Gosched()

	select {
	case <-mp.C:
		runtime.Gosched()
		select {
		case <-mp.C:
			t.Fatalf("only one end event should have been sent")
		default:
		}
	default:
		t.Fatalf("no end event was sent")
	}
}

func TestSessionRemovalOnTimeout(t *testing.T) {
	mp := NewMonitorPool(map[string]time.Duration{
		"heartbeat": time.Millisecond * 3,
	})

	// test the session lookup
	s := mp.session("key1")

	s.end("heartbeat", &Timeout{
		Start:            time.Now(),
		LastNotification: time.Now(),
		End:              time.Now(),
	})

	runtime.Gosched()

	if _, ok := mp.sessions["key1"]; ok {
		t.Fatalf("session should have been removed when it ended")
	}
}
