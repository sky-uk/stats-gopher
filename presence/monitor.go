package presence

import "time"

type monitor struct {
	c                chan *Timeout
	kill             chan bool
	timer            *time.Timer
	timeout          time.Duration
	start            time.Time
	lastNotification time.Time
	cancelled        bool
}

func newMonitor(pulseTimeout time.Duration) (m *monitor) {
	start := time.Now()

	m = &monitor{
		c:                make(chan *Timeout, 1),
		kill:             make(chan bool, 1),
		timer:            time.NewTimer(pulseTimeout),
		timeout:          pulseTimeout,
		start:            start,
		lastNotification: start,
	}

	go m.run()

	return
}

func (m *monitor) pulse() {
	m.timer.Reset(m.timeout)
	m.lastNotification = time.Now()
}

func (m *monitor) cancel() {
	defer func() {
		recover()
	}()

	m.cancelled = true
	m.kill <- true
}

func (m *monitor) run() {
	m.timer.Reset(m.timeout)

	select {
	case <-m.kill:
	case <-m.timer.C:
		m.timeoutNow()
	}

	m.timer.Stop()
	close(m.c)
	close(m.kill)
}

func (m *monitor) timeoutNow() {
	end := time.Now()
	m.c <- &Timeout{
		Start:            m.start,
		LastNotification: m.lastNotification,
		End:              end,
		Duration:         m.lastNotification.Sub(m.start),
		Wait:             m.timeout,
	}
}
