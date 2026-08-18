package main

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Sut103/EvePing-for-Discord/internal/scheduler"
)

func TestWaitAndStop_SignalStopsScheduler(t *testing.T) {
	sched := scheduler.New(time.Hour, func() {})
	sigCh := make(chan os.Signal, 1)
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)

	done := make(chan struct{})
	go func() {
		sched.Start()
		close(done)
	}()

	go waitAndStop(sched, sigCh, logger)
	sigCh <- os.Interrupt

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not stop after signal")
	}
}
