# AGENTS.md — Collaboration protocol (Codex + Claude)

Two AI engineers build this project in alternation: **Codex** (ChatGPT) and **Claude** (Claude Code). Either may be cut off mid-task by a usage/token limit **without warning**. This protocol makes the work resumable by the other at any second. Codex reads this file automatically; Claude reads `CLAUDE.md` (which points here). **Both follow these rules.**

## The baton

- **`HANDOFF-STATE.md`** = the single source of truth for *where we are right now and what to do next*. It always reflects the current reality.
- **`WORKLOG.md`** = append-only history (newest entry on top). One short entry per work session or meaningful chunk.

## Every session — do this in order

1. **On start, READ FIRST:** `HANDOFF-STATE.md`, then the top few entries of `WORKLOG.md`. Then `git status` + `git log --oneline -10` to see the real state of the code. Trust the code + git over prose if they disagree, and fix the note.
2. **Sync the code:** if `HANDOFF-STATE.md` names a branch/commit, make sure you're on it and up to date before touching anything.
3. **Work** on the ONE "Next action" in `HANDOFF-STATE.md` (or the next roadmap step if that's done). Keep changes small.
4. **Commit often** — small, reviewable commits. The other agent inherits your *committed* work, so uncommitted work is invisible to them.
5. **Update the baton CONTINUOUSLY, not just at the end** — because a limit may cut you off before any "final" update. After each meaningful step: refresh `HANDOFF-STATE.md` and, when you finish a chunk, prepend a `WORKLOG.md` entry.

## Before you stop (or when you sense a limit approaching)

- Commit (or clearly note any intentionally-uncommitted state and why).
- Make `HANDOFF-STATE.md` true: what's DONE+verified, what's IN PROGRESS and exactly where, the SINGLE next action, blockers, and how to run/verify.
- Prepend a `WORKLOG.md` entry: `## <ISO datetime> · <Codex|Claude> · <one-line summary>` + 2–5 bullets (did / state / next / blockers).

## Rules of engagement

- **One next action at a time.** Don't fan out; leave the baton pointing at a single clear step.
- **Never assume the other's uncommitted work exists.** If the baton says "in progress, uncommitted," re-derive from the code.
- **Author commits as `mikechibuzor`. Never add a `Co-Authored-By` trailer** (neither Claude nor Codex).
- **Don't silently change decisions** locked in the numbered handoff docs (`00`–`06`). If you must deviate, note it in `HANDOFF-STATE.md` under "Decisions/deviations" and flag the user.
- **Ask the user before anything irreversible** (deploys, deleting data, spending money, sending real messages/emails, live payment charges).
- **Keep the ghost admin unlinked and unmentioned to the client** (see `01-PROJECT-BRIEF.md`).
- **Golden invariants:** money in kobo (int64, never floats); an order is `paid` only after server-side verify/verified webhook; mobile-first; guest checkout only in phase 1.
- **Frontend typing:** strict TypeScript at all times. Define shared types and runtime values in dedicated modules and import them wherever reused.
- **Frontend constants:** centralize domain statuses, API paths, configuration defaults, and reusable UI copy. Prefer `as const` + inferred literal unions; use enums when runtime enum semantics improve clarity.
- **Frontend styling:** use SCSS, not plain CSS. Keep shared tokens and mixins in `web/assets/scss/design-system/`, global element rules in `_globals.scss`, and page-specific rules in page partials or scoped Vue `<style lang="scss">` blocks.
- **Code walkthroughs:** add a one-line intent comment immediately above every function and computed property. Inside multi-step functions, add concise comments at meaningful step transitions so the implementation reads as a walkthrough; explain purpose, not syntax.
- **Code quality:** use descriptive names and optimize for simplicity, readability, scalability, and maintainability.

## Project map (what to read)

`00-README-START-HERE.md` → `01-PROJECT-BRIEF.md` → `02-DESIGN-SPEC.md` → `03-ARCHITECTURE.md` → `04-DATA-MODEL.md` → `05-BUILD-ROADMAP.md` → `06-SETUP-AND-ENV.md`. Data in `catalogue/`, visuals in `design-source/`, primary docs in `reference-docs/`.

> **When the code repo is scaffolded (Milestone A, step 1), copy these four files — `AGENTS.md`, `CLAUDE.md`, `HANDOFF-STATE.md`, `WORKLOG.md` — to the repo root** so the baton lives beside the code and both agents auto-load it there.
