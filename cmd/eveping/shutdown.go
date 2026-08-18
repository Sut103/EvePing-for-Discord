package main

import (
	"log"
	"os"

	"github.com/Sut103/EvePing-for-Discord/internal/scheduler"
)

// waitAndStop blocks until a shutdown signal arrives on sigCh, then stops
// sched so Start() returns and main can run its deferred cleanup (closing
// the Discord session). sigCh is injected so this is testable without
// sending a real OS signal.
func waitAndStop(sched *scheduler.Scheduler, sigCh <-chan os.Signal, logger *log.Logger) {
	sig := <-sigCh
	logger.Printf("received signal %s, shutting down", sig)
	sched.Stop()
}
