package presence

import "time"

// MonitorPool handles presence notifications
type MonitorPool struct {
	c        chan<- Timeout
	C        <-chan Timeout
	timeouts map[string]time.Duration
	sessions map[string]*session
}

// Notification describes the details of a notification of presence
type Notification struct {
	Key  string
	Code string
}

// NewMonitorPool creates a Monitor with default timeout values
// 45 minutes for a user event
// 30 seconds for a browser heartbeat
func NewMonitorPool(timeouts map[string]time.Duration) *MonitorPool {
	c := make(chan Timeout, 65536)

	return &MonitorPool{
		c:        c,
		C:        c,
		timeouts: timeouts,
		sessions: make(map[string]*session),
	}
}

// Notify the monitor pool of the presence of something
func (mp *MonitorPool) Notify(n *Notification) {
	mp.session(n.Key).pulse(n.Code)
}

func (mp *MonitorPool) session(key string) *session {
	if session, ok := mp.sessions[key]; ok {
		return session
	}

	session := newSession(key)

	for name, timeout := range mp.timeouts {
		session.monitor(name, timeout)
		go func() {
			timeout := <-session.c
			if timeout == nil {
				return
			}
			mp.c <- *timeout
			delete(mp.sessions, key)
		}()
	}

	mp.sessions[key] = session

	return session
}
