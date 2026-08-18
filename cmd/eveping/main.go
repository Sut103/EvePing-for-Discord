// Command eveping is a stateless Discord bot that sends a DM reminder to
// every user interested in a guild's scheduled event, one day before it
// starts. It runs as a resident process: once per day it scans every guild
// the bot has joined, so there is no database or cache to keep in sync.
package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Sut103/EvePing-for-Discord/internal/batch"
	"github.com/Sut103/EvePing-for-Discord/internal/discordclient"
	"github.com/Sut103/EvePing-for-Discord/internal/scheduler"
)

const (
	tokenEnvVar   = "EVEPING_DISCORD_TOKEN"
	batchInterval = 24 * time.Hour
)

var errMissingToken = errors.New(tokenEnvVar + " environment variable is not set")

func loadToken(getenv func(string) string) (string, error) {
	token := getenv(tokenEnvVar)
	if token == "" {
		return "", errMissingToken
	}
	return token, nil
}

func logBatchResult(logger *log.Logger, start time.Time, result batch.BatchResult) {
	logger.Printf(
		"daily batch finished: duration=%s target_events=%d sent_success=%d sent_failure=%d error_count=%d",
		time.Since(start).Round(time.Millisecond),
		result.TargetEvents,
		result.SentSuccess,
		result.SentFailure,
		len(result.Errors),
	)
	for _, err := range result.Errors {
		logger.Printf("daily batch error: %v", err)
	}
}

func runDailyBatch(client discordclient.Client, logger *log.Logger) {
	logger.Println("daily batch starting")
	start := time.Now()
	result := batch.RunDailyBatch(client, time.Now())
	logBatchResult(logger, start, result)
}

func main() {
	logger := log.Default()

	token, err := loadToken(os.Getenv)
	if err != nil {
		logger.Fatal(err)
	}

	// discordgo.New sets Identify.Intents to IntentsAllWithoutPrivileged by
	// default, which already includes IntentGuildScheduledEvents — no
	// explicit intent configuration is needed here.
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Fatalf("create discord session: %v", err)
	}

	if err := session.Open(); err != nil {
		logger.Fatalf("open discord session: %v", err)
	}
	defer session.Close()

	client := discordclient.New(session)

	sched := scheduler.New(batchInterval, func() {
		runDailyBatch(client, logger)
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go waitAndStop(sched, sigCh, logger)

	sched.Start()
}
