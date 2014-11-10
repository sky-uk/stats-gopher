package presence

import (
	"runtime"
	"testing"
	"time"
)

func TestSessionTimeout(t *testing.T) {
	session := newSession("key1")

	session.monitor("fishing", time.Minute)
	session.monitor("rowing", time.Second)

	rowing := session.monitors["rowing"]

	t0 := time.Now()

	rowing.timeoutNow()

	fishing := session.monitors["fishing"]

	runtime.Gosched()

	if !fishing.cancelled {
		t.Fatalf("other monitors should have been cancelled")
	}

	select {
	case timeout := <-session.c:
		if timeout.Key != "key1" {
			t.Fatalf("end event key should be the session key")
		}
		if timeout.Code != "rowing" {
			t.Fatalf("end event code should be the monitor name")
		}
		if timeout.Start != rowing.start {
			t.Fatalf("start event time should be the same as the timeout")
		}
		if timeout.LastNotification != rowing.lastNotification {
			t.Fatalf("lastNotification event time should be the same as the timeout")
		}
		if timeout.End.Unix() < t0.Unix() {
			t.Fatalf("end event time was before the timeout occurred")
		}
		if timeout.End.Unix() > time.Now().Unix() {
			t.Fatalf("end event time was after the timeout occurred")
		}
	default:
		t.Fatalf("the monitor timeout did not result in an end event")
	}
}

func TestSessionPulse(t *testing.T) {
	session := newSession("key")

	session.monitor("fishing", time.Hour)
	session.monitor("rowing", time.Minute)

	rowing := session.monitors["rowing"]

	t0 := time.Now()

	session.pulse("rowing")

	t1 := time.Now()
	tLastNotificationUnix := rowing.lastNotification.Unix()

	if tLastNotificationUnix < t0.Unix() {
		t.Fatalf("the last notification time was before the pulse")
	}

	if tLastNotificationUnix > t1.Unix() {
		t.Fatalf("the last notification time was after the pulse")
	}
}
