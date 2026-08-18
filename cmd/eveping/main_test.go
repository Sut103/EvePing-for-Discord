package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/Sut103/EvePing-for-Discord/internal/batch"
)

func TestLoadToken_MissingEnv_ReturnsError(t *testing.T) {
	getenv := func(key string) string { return "" }

	_, err := loadToken(getenv)

	if err == nil {
		t.Fatal("expected an error when the token env var is unset, got nil")
	}
}

func TestLoadToken_PresentEnv_ReturnsToken(t *testing.T) {
	getenv := func(key string) string {
		if key == tokenEnvVar {
			return "secret-token"
		}
		return ""
	}

	token, err := loadToken(getenv)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "secret-token" {
		t.Fatalf("token = %q, want %q", token, "secret-token")
	}
}

func TestLogBatchResult_IncludesCounts(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	result := batch.BatchResult{TargetEvents: 3, SentSuccess: 5, SentFailure: 2}

	logBatchResult(logger, time.Now(), result)

	out := buf.String()
	for _, want := range []string{"target_events=3", "sent_success=5", "sent_failure=2"} {
		if !strings.Contains(out, want) {
			t.Fatalf("log output %q does not contain %q", out, want)
		}
	}
}
