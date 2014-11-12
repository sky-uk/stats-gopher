package presence

import (
	"testing"
	"time"
)

type monitorClient struct {
	monitor  *monitor
	timeouts []*Timeout
}

func newMonitorClient(pulseTimeout time.Duration) *monitorClient {
	mc := &monitorClient{
		monitor:  newMonitor(pulseTimeout),
		timeouts: make([]*Timeout, 0, 100),
	}

	go func() {
		for {
			if e, ok := <-mc.monitor.c; ok {
				mc.timeouts = append(mc.timeouts, e)
			} else {
				break
			}
		}
	}()

	return mc
}

func TestMonitorTimeout(t *testing.T) {
	pulseTimeout := time.Millisecond * 10

	t0 := time.Now()

	mc := newMonitorClient(pulseTimeout)

	if mc.monitor.lastNotification != mc.monitor.start {
		t.Fatalf("lastNotification should initially be the start time")
	}

	t1 := time.Now()

	mc.monitor.pulse()

	time.Sleep(pulseTimeout / 2)

	t2 := time.Now()

	mc.monitor.pulse()

	if len(mc.timeouts) > 0 {
		t.Fatalf("EndEvents have been emitted but no heartbeats have stopped")
	}

	t3 := time.Now()

	time.Sleep(time.Duration(float64(pulseTimeout) * 2))

	if len(mc.timeouts) != 1 {
		t.Fatalf("EndEvents were not emitted correctly: expect 1 event but got %d events", len(mc.timeouts))
	}

	timeout := mc.timeouts[0]

	if timeout.Start.Unix() < t0.Unix() {
		t.Fatalf("the start time was set to a time before monitor was created")
	}

	if timeout.Start.Unix() > t1.Unix() {
		t.Fatalf("the start time was set to a time after the monitor was created")
	}

	if timeout.LastNotification.Unix() < t2.Unix() {
		t.Fatalf("the lastEvent time was set to a time before the lastEvent was received")
	}

	if timeout.LastNotification.Unix() > t3.Unix() {
		t.Fatalf("the lastEvent time was set to a time after the lastEvent was received")
	}

	tEndUnix := t3.Unix() + int64(pulseTimeout)

	if timeout.End.Unix() < t3.Unix() {
		t.Fatalf("the end time was set to a time before the monitor timed out")
	}

	if timeout.End.Unix() > tEndUnix {
		t.Fatalf("the end time was set to the time after the monitor timed out")
	}

	if timeout.Duration != timeout.LastNotification.Sub(timeout.Start) {
		t.Fatalf("the duration should be the difference between the start and last notification time")
	}

	if timeout.Wait != mc.monitor.timeout {
		t.Fatalf("the wait time for the timeout should be set to the timeout for the monitor")
	}
}

func TestMonitorCancel(t *testing.T) {
	pulseTimeout := time.Millisecond * 5

	mc := newMonitorClient(pulseTimeout)

	mc.monitor.pulse()

	mc.monitor.cancel()

	time.Sleep(pulseTimeout)
	time.Sleep(pulseTimeout)

	if len(mc.timeouts) > 0 {
		t.Fatalf("EndEvents have been emitted but the monitor was cancelled")
	}

	// must not cause a panic
	mc.monitor.cancel()
}
