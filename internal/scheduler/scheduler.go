// Package scheduler runs a callback on a fixed interval for the lifetime of
// the process.
package scheduler

import "time"

// Scheduler invokes run every interval until Stop is called. The interval
// is injected so production code can use 24 hours while tests use a much
// shorter duration.
type Scheduler struct {
	interval time.Duration
	run      func()
	stop     chan struct{}
}

// New creates a Scheduler that calls run once per interval once Start is
// called.
func New(interval time.Duration, run func()) *Scheduler {
	return &Scheduler{interval: interval, run: run, stop: make(chan struct{})}
}

// Start blocks, calling run() immediately and then again every interval,
// until Stop is called. The immediate first call matters in production: on
// every restart (crash, redeploy, host reboot) it re-evaluates "tomorrow"
// right away instead of leaving a up-to-interval gap (24h) during which
// events starting tomorrow would never be checked, since eveping keeps no
// persisted state to catch up from later.
func (s *Scheduler) Start() {
	select {
	case <-s.stop:
		return
	default:
		s.run()
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.run()
		case <-s.stop:
			return
		}
	}
}

// Stop ends the Start loop. It must be called at most once.
func (s *Scheduler) Stop() {
	close(s.stop)
}
