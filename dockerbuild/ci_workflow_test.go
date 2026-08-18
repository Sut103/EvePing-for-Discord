package dockerbuild

import (
	"strings"
	"testing"
)

func TestCIWorkflowRunsDockerBuild(t *testing.T) {
	content := readRepoFile(t, ".github/workflows/ci.yml")
	if !strings.Contains(content, "docker build") {
		t.Fatal("expected .github/workflows/ci.yml to run `docker build` to verify the Dockerfile builds")
	}
}

func TestCIWorkflowDoesNotPushImage(t *testing.T) {
	content := readRepoFile(t, ".github/workflows/ci.yml")
	for _, forbidden := range []string{"docker push", "docker/login-action", "docker/build-push-action"} {
		if strings.Contains(content, forbidden) {
			t.Fatalf(".github/workflows/ci.yml must not push images to a registry; found %q", forbidden)
		}
	}
}
