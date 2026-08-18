# Conventions — File Formats, Paths & Rules

Read this before doing any file operations across all phases.

---

## Directory Structure

```
.claude/
├── prds/
│   └── <feature-name>.md          # Product requirement documents
├── epics/
│   ├── <feature-name>/
│   │   ├── epic.md                # Technical epic
│   │   ├── <N>.md                 # Task files (named by GitHub issue number after sync)
│   │   ├── <N>-analysis.md        # Parallel work stream analysis
│   │   ├── github-mapping.md      # Issue number → URL mapping
│   │   ├── execution-status.md    # Active agents tracker
│   │   └── updates/
│   │       └── <issue_N>/
│   │           ├── stream-A.md    # Per-agent progress
│   │           ├── progress.md    # Overall issue progress
│   │           └── execution.md  # Execution state
│   └── archived/
│       └── <feature-name>/        # Completed epics
└── context/                       # Project context docs (separate system)
```

---

## Frontmatter Schemas

### PRD (.claude/prds/<name>.md)
```yaml
---
name: <feature-name>        # kebab-case, matches filename
description: <one-liner>    # used in lists and summaries
status: backlog | active | completed
created: <ISO 8601>         # date -u +"%Y-%m-%dT%H:%M:%SZ"
---
```

### Epic (.claude/epics/<name>/epic.md)
```yaml
---
name: <feature-name>
status: backlog | in-progress | completed
created: <ISO 8601>
updated: <ISO 8601>
progress: 0%                # recalculated when tasks close
prd: .claude/prds/<name>.md
github: https://github.com/<owner>/<repo>/issues/<N>  # set on sync
---
```

### Task (.claude/epics/<name>/<N>.md)
```yaml
---
name: <Task Title>
status: open | in-progress | closed
created: <ISO 8601>
updated: <ISO 8601>
github: https://github.com/<owner>/<repo>/issues/<N>  # set on sync
depends_on: []              # issue numbers this must wait for
parallel: true              # can run concurrently with non-conflicting tasks
conflicts_with: []          # issue numbers that touch the same files
---
```

### Progress (.claude/epics/<name>/updates/<N>/progress.md)
```yaml
---
issue: <N>
started: <ISO 8601>
last_sync: <ISO 8601>
completion: 0%
---
```

---

## TDD Rule (Red → Green)

Every implementation task in CCPM follows strict test-driven development. This is not optional and is not delegated to a separate "testing" task — it is how each unit of implementation work gets built.

1. **Red** — before writing any implementation code, write a test that encodes one acceptance criterion and confirm it fails (and fails for the expected reason, not a typo/setup error).
2. **Green** — write the minimum implementation code needed to make that test pass. Run the full relevant test suite and confirm it passes with no other tests broken.
3. **Refactor** — once green, clean up implementation and test code while keeping the suite passing.

Repeat per acceptance criterion / test case until the task's Test Plan is fully covered.

**Where this shows up:**
- **Epics** define a `Test Strategy` (frameworks, test types, coverage expectations) before decomposition — see `plan.md`.
- **Tasks** define a `Test Plan` (concrete test cases derived from acceptance criteria) at creation time, before `Technical Details` — see `structure.md`.
- **Execution** requires agents to show the red step (failing test) before the green step (passing implementation) for each change — see `execute.md`.

**Hard rules:**
- Never write implementation code for a behavior that doesn't yet have a failing test.
- Never check "Tests written and passing" on a task without having run them red first.
- A task with no Test Plan is not ready to build — go back and write one.

---

## Code Review Rule (per-task `/code-review`)

Every task's implementation is verified with the `/code-review` skill right after the agent finishes it, before the stream is marked completed. This per-stream review is not optional, and it does not substitute for the epic-close multi-agent review below — both run, at different points in the lifecycle.

1. Once all Test Plan items are green, run `/code-review` against the stream's diff.
2. **Critical findings** (broken build, a Test Plan item that doesn't actually pass, security vulnerability, data loss) — fix immediately, re-run the relevant tests, and re-review before proceeding. Do not mark the stream completed with an open critical finding.
3. **Non-critical findings** (style, minor simplification, low-confidence suggestions) — record them in the stream's progress file for later triage; they do not block marking the stream completed.

**Where this shows up:**
- **Execution** — each agent runs `/code-review` immediately after its Test Plan goes green and before marking status: completed — see `execute.md`.

---

## Epic Close Review Rule (multi-agent `/code-review`)

Once every task in an epic is closed and before the epic is merged (`Merging an Epic` in `sync.md`), the full epic diff gets one more pass — a multi-agent code review, not just a rerun of the per-stream reviews above. This is not optional and is not skipped even if every per-stream review already passed clean.

1. Launch multiple agents in parallel, each independently running the `/code-review` skill against the full epic diff (the epic worktree/branch diff vs. `main`), not just a single task's slice of it.
2. Consolidate the agents' findings into one report — dedupe overlapping items, and keep the higher severity when reports disagree — split into critical and non-critical findings.
3. **Report the consolidated results to the user** before proceeding to merge.
4. **Critical findings** — fix immediately, re-run affected tests, and re-review before merging. Do not merge an epic with an open critical finding.
5. **Non-critical findings** — record them (e.g. in the epic's progress file) for later triage; they do not block the merge.

**Where this shows up:**
- **Sync** — runs before the merge steps in "Merging an Epic," with the consolidated results reported to the user — see `sync.md`.

---

## Datetime Rule

Always get real current datetime from the system — never use placeholder text:
```bash
date -u +"%Y-%m-%dT%H:%M:%SZ"
```

---

## Frontmatter Update Pattern

When updating a single frontmatter field in an existing file:
```bash
sed -i.bak "/^<field>:/c\\<field>: <value>" <file>
rm <file>.bak
```

When stripping frontmatter to get body content for GitHub:
```bash
sed '1,/^---$/d; 1,/^---$/d' <file> > /tmp/body.md
```

---

## GitHub Operations

### Repository Safety Check (run before any write operation)
```bash
remote_url=$(git remote get-url origin 2>/dev/null || echo "")
if [[ "$remote_url" == *"automazeio/ccpm"* ]]; then
  echo "❌ Cannot write to the CCPM template repository."
  echo "Update remote: git remote set-url origin https://github.com/YOUR/REPO.git"
  exit 1
fi
REPO=$(echo "$remote_url" | sed 's|.*github.com[:/]||' | sed 's|\.git$||')
```

### Authentication
Don't pre-check authentication. Run the `gh` command and handle failure:
```bash
gh <command> || echo "❌ GitHub CLI failed. Run: gh auth login"
```

### Getting Issue Numbers
```bash
# From a task file's github field:
grep 'github:' <file> | grep -oE '[0-9]+$'
```

---

## Git / Worktree Conventions

- One branch per epic: `epic/<name>`
- Worktrees live at `../epic-<name>/` (sibling to project root)
- Always start branches from an up-to-date main:
  ```bash
  git checkout main && git pull origin main
  git worktree add ../epic-<name> -b epic/<name>
  ```
- Commit format inside epics: `Issue #<N>: <description>`
- Never use `--force` in any git operation
- **Verified working in Claude Code on the Web**: `git worktree add`/`remove`, and committing inside the resulting worktree directory, function normally in a Claude Code on the Web (remote/cloud) session — confirmed by creating, writing/committing in, and removing a worktree in that environment. Do not assume worktrees are unavailable there; use the epic worktree flow as-is regardless of whether the session is local or web/remote.

---

## Naming Conventions

- Feature names: kebab-case, lowercase, letters/numbers/hyphens, starts with a letter
- Task files before sync: `001.md`, `002.md`, ... (sequential)
- Task files after sync: renamed to GitHub issue number (e.g., `1234.md`)
- Labels applied on sync: `epic`, `epic:<name>`, `feature` (for epics); `task`, `epic:<name>` (for tasks)

---

## Epic Progress Calculation

```bash
total=$(ls .claude/epics/<name>/[0-9]*.md 2>/dev/null | wc -l)
closed=$(grep -l '^status: closed' .claude/epics/<name>/[0-9]*.md 2>/dev/null | wc -l)
progress=$((closed * 100 / total))
```

Update epic frontmatter when any task closes.
