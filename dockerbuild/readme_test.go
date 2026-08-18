package dockerbuild

import (
	"strings"
	"testing"
)

func TestReadmeDocumentsDockerBuildAndRun(t *testing.T) {
	content := readRepoFile(t, "README.md")
	for _, want := range []string{"docker build", "docker run", "EVEPING_DISCORD_TOKEN"} {
		if !strings.Contains(content, want) {
			t.Fatalf("expected README.md to mention %q", want)
		}
	}
}

func TestReadmeDocumentsDockerCompose(t *testing.T) {
	content := readRepoFile(t, "README.md")
	if !strings.Contains(content, "docker compose") && !strings.Contains(content, "docker-compose") {
		t.Fatal("expected README.md to mention docker compose / docker-compose")
	}
}
