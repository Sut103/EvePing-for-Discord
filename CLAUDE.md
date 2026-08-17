# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language

Always respond to the user in Japanese (日本語). This applies to all chat output — explanations, summaries, and questions — regardless of the language the user writes in. Code, identifiers, commit messages, and file contents are unaffected by this rule and should follow normal engineering conventions.

## Project state

This repository is currently a greenfield project — it contains no application source code yet, only a LICENSE (MIT) and project-management tooling. There is no README, package manifest, build system, linter, or test suite to reference. Do not assume any particular language, framework, or directory layout until one has actually been established in the repo; check what exists before writing commands or code.

The project name ("EvePing for Discord") implies a Discord bot/integration related to EVE Online, but no implementation has been started.

## Project management workflow (ccpm)

This repo has the `ccpm` skill installed (`.claude/skills/ccpm/`), a spec-driven delivery workflow: **PRD → Epic → GitHub Issues → parallel agents → shipped code**. Use it for anything in the software delivery lifecycle — writing/planning features, decomposing epics into tasks, syncing work to GitHub issues, starting work on an issue, checking status/standup, or closing out issues/epics.

Key points from `.claude/skills/ccpm/SKILL.md` and `references/conventions.md`:

- **Requirements live in files, not heads.** Every feature starts as a PRD, becomes a technical epic, decomposes into GitHub issues, and is executed by parallel agents with full traceability.
- **TDD is enforced across the whole lifecycle** (red → green), not treated as a separate phase:
  - Plan: every epic defines a concrete Test Strategy before decomposition.
  - Structure: every task defines a Test Plan (one test case per acceptance criterion) before Technical Details.
  - Execute: every agent writes the failing test before the passing implementation, for each Test Plan item — never the reverse.
- **Script-first rule**: for deterministic, read-only status queries, run the provided bash scripts in `.claude/skills/ccpm/references/scripts/` instead of reasoning it out manually — e.g. `status.sh`, `standup.sh`, `epic-list.sh`, `epic-show.sh <name>`, `epic-status.sh <name>`, `prd-list.sh`, `prd-status.sh`, `search.sh <query>`, `in-progress.sh`, `next.sh`, `blocked.sh`, `validate.sh`. Reserve LLM reasoning for work that actually needs it: writing PRDs, analyzing parallelism, launching agents, synthesizing updates.
- The five phases (Plan, Structure, Sync, Execute, Track) each have a dedicated reference doc under `.claude/skills/ccpm/references/` — read the relevant one before acting in that phase.

## Once application code exists

When source files, a package manifest, or a test suite are added to this repo, update this CLAUDE.md with the actual build/lint/test commands and a description of the real architecture — do not leave this section aspirational.
