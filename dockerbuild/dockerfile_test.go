// Package dockerbuild contains static assertion tests for the repository's
// Docker-related files (Dockerfile, .dockerignore, docker-compose.yml,
// CI workflow). These read the files as plain text so they run without a
// Docker daemon, unlike an actual `docker build`.
package dockerbuild

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func readRepoFile(t *testing.T, relPath string) string {
	t.Helper()
	content, err := os.ReadFile("../" + relPath)
	if err != nil {
		t.Fatalf("read %s: %v", relPath, err)
	}
	return string(content)
}

func TestDockerfileIsMultiStage(t *testing.T) {
	content := readRepoFile(t, "Dockerfile")
	fromRe := regexp.MustCompile(`(?m)^FROM\s+`)
	matches := fromRe.FindAllString(content, -1)
	if len(matches) < 2 {
		t.Fatalf("expected Dockerfile to have at least 2 FROM instructions (multi-stage build), got %d", len(matches))
	}
}

func TestDockerfileUsesCGoDisabled(t *testing.T) {
	content := readRepoFile(t, "Dockerfile")
	if !strings.Contains(content, "CGO_ENABLED=0") {
		t.Fatal("expected Dockerfile to set CGO_ENABLED=0 for a static build")
	}
}

func TestDockerfileDoesNotHardcodeToken(t *testing.T) {
	content := readRepoFile(t, "Dockerfile")
	if strings.Contains(content, "EVEPING_DISCORD_TOKEN") {
		t.Fatal("Dockerfile must not reference EVEPING_DISCORD_TOKEN; it must be injected at container runtime only")
	}
}

func TestDockerfileRunsAsNonRoot(t *testing.T) {
	content := readRepoFile(t, "Dockerfile")
	userRe := regexp.MustCompile(`(?m)^USER\s+\S+`)
	if !userRe.MatchString(content) {
		t.Fatal("expected Dockerfile to switch to a non-root user with a USER instruction")
	}
}

func TestDockerfileInstallsCACertificates(t *testing.T) {
	content := readRepoFile(t, "Dockerfile")
	if !strings.Contains(content, "ca-certificates") {
		t.Fatal("expected Dockerfile to install ca-certificates in the runtime stage, otherwise TLS calls to the Discord API fail")
	}
}

func TestDockerignoreExcludesGitAndClaude(t *testing.T) {
	content := readRepoFile(t, ".dockerignore")
	for _, entry := range []string{".git", ".claude"} {
		if !strings.Contains(content, entry) {
			t.Fatalf("expected .dockerignore to exclude %q", entry)
		}
	}
}
