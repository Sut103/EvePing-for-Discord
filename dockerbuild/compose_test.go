package dockerbuild

import (
	"regexp"
	"strings"
	"testing"
)

func TestComposeFileReferencesTokenEnvVar(t *testing.T) {
	content := readRepoFile(t, "docker-compose.yml")
	if !strings.Contains(content, "EVEPING_DISCORD_TOKEN") {
		t.Fatal("expected docker-compose.yml to reference EVEPING_DISCORD_TOKEN via environment/env_file")
	}
	hardcoded := regexp.MustCompile(`EVEPING_DISCORD_TOKEN=[^$\s]`)
	if hardcoded.MatchString(content) {
		t.Fatal("EVEPING_DISCORD_TOKEN must not have a hardcoded value in docker-compose.yml; inject it via a variable")
	}
}

func TestComposeFileBuildsFromRepoDockerfile(t *testing.T) {
	content := readRepoFile(t, "docker-compose.yml")
	buildRe := regexp.MustCompile(`(?m)^\s*build:\s*\.\s*$`)
	if !buildRe.MatchString(content) {
		t.Fatal("expected docker-compose.yml to build from the repository root (build: .)")
	}
}
