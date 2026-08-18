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
	// Accepts either the short form (`build: .`) or the long form
	// (`build:` with a nested `context: .`), so reformatting the compose
	// file to add build args or an explicit dockerfile path doesn't
	// spuriously break this check.
	shortForm := regexp.MustCompile(`(?m)^\s*build:\s*\.\s*$`)
	longForm := regexp.MustCompile(`(?m)^\s*build:\s*$[\s\S]*?^\s*context:\s*\.\s*$`)
	if !shortForm.MatchString(content) && !longForm.MatchString(content) {
		t.Fatal("expected docker-compose.yml to build from the repository root (build: . or build.context: .)")
	}
}
