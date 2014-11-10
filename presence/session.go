package presence

import "time"

type session struct {
	key      string
	c        chan *Timeout
	monitors map[string]*monitor
}

func newSession(key string) *session {
	return &session{
		key:      key,
		c:        make(chan *Timeout, 1),
		monitors: make(map[string]*monitor),
	}
}

func (s *session) monitor(name string, timeout time.Duration) {
	m := newMonitor(timeout)
	s.monitors[name] = m
	go s.waitForTimeout(name, m)
}

func (s *session) pulse(name string) {
	s.monitors[name].pulse()
}

func (s *session) waitForTimeout(name string, m *monitor) {
	timeout := <-m.c

	// the channel could have closed
	if timeout == nil {
		return
	}

	s.end(name, timeout)
}

func (s *session) end(name string, timeout *Timeout) {
	for _, m := range s.monitors {
		m.cancel()
	}

	timeout.Key = s.key
	timeout.Code = name

	select {
	case s.c <- timeout:
	default:
		return
	}
}
