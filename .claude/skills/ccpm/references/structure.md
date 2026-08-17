# Structure — Break Down an Epic

This phase converts a technical epic into concrete, numbered task files with dependency and parallelization metadata.

---

## Epic Decomposition

**Trigger**: User wants to break an epic into actionable tasks.

### Preflight
- Verify `.claude/epics/<name>/epic.md` exists with valid frontmatter.
- If numbered task files (001.md, 002.md...) already exist in the epic directory, list them and confirm deletion before recreating.
- If epic status is "completed", warn the user before proceeding.

### Process

Read the epic fully. Analyze for parallelism — which pieces of work can happen simultaneously without file conflicts?

**Task types to consider:**
- Setup: environment, scaffolding, dependencies
- Data: models, schemas, migrations
- API: endpoints, services, integration
- UI: components, pages, styling
- Docs: README, API docs, changelogs

Every task above that produces behavior (Data/API/UI, etc.) carries its own tests written TDD-style (red → green, see `references/conventions.md`) — testing is not a separate task type. Only use a standalone "Tests" task for suites that genuinely span multiple implementation tasks (e.g. a cross-feature e2e flow) and can't be attached to a single task's Test Plan.

**Parallelization strategy by epic size:**
- Small (<5 tasks): create sequentially
- Medium (5–10 tasks): batch into 2–3 groups, spawn parallel Task agents
- Large (>10 tasks): analyze dependencies first, launch parallel agents (max 5 concurrent), create dependent tasks after prerequisites

For parallel creation, use the Task tool:
```yaml
Task:
  description: "Create task files batch N"
  subagent_type: "general-purpose"
  prompt: |
    Create task files for epic: <name>
    Tasks to create: [list 3-4 tasks]
    Save to: .claude/epics/<name>/001.md, 002.md, etc.
    Follow the task file format exactly.
    Return: list of files created.
```

### Task File Format

```markdown
---
name: <Task Title>
status: open
created: <run: date -u +"%Y-%m-%dT%H:%M:%SZ">
updated: <same as created>
github: (will be set on sync)
depends_on: []
parallel: true
conflicts_with: []
---

# Task: <Task Title>

## Description

## Acceptance Criteria
- [ ]

## Test Plan
- [ ] <failing test to write first (red)> — covers acceptance criterion above
- [ ] <failing test to write first (red)> — covers acceptance criterion above

## Technical Details

## Dependencies

## Effort Estimate
- Size: XS/S/M/L/XL
- Hours: N

## Definition of Done
- [ ] Failing tests written first (red), one per Test Plan item, and confirmed to fail for the right reason
- [ ] Implementation makes all Test Plan tests pass (green)
- [ ] Full relevant test suite passes — no other tests broken
- [ ] Code reviewed
```

Every acceptance criterion must map to at least one Test Plan entry — a task with acceptance criteria but no corresponding test case is incomplete. See `references/conventions.md` for the TDD Rule this enforces.

**Numbering**: sequential 001.md, 002.md, etc. Tasks are renamed to GitHub issue numbers after sync — do not hard-code dependencies by filename, use the `depends_on` array.

### After Creating All Tasks

Append a summary to the epic file:

```markdown
## Tasks Created
- [ ] 001.md - <Title> (parallel: true/false)
- [ ] 002.md - <Title> (parallel: true/false)

Total tasks: N
Parallel tasks: N
Sequential tasks: N
Estimated total effort: N hours
```

**After completion**: Confirm "✅ Created N tasks for epic: <name>" and suggest: "Ready to push to GitHub? Say: sync the <name> epic"

---

## Dependency Rules
- `depends_on` lists task numbers that must complete before this task can start.
- `parallel: true` means the task can run concurrently with others it doesn't conflict with.
- `conflicts_with` lists tasks that touch the same files — these cannot run in parallel.
- Circular dependencies are an error — check before finalizing.
