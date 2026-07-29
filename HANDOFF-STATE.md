# HANDOFF-STATE.md — the baton

**Read this first. Keep it true at all times.** This is the single source of truth for where we are and what to do next. Update it continuously (a usage limit may cut you off without warning).

---

## Current status
- **Phase:** 1 (MVP). First 40% milestone **paid**; build authorized.
- **Milestone:** A — Foundations & catalogue (see `05-BUILD-ROADMAP.md`).
- **Step:** A.1 — repo not yet scaffolded. **Nothing built yet.**
- **Last worked by:** Claude · 2026-07-29 (created handoff package + collaboration kit).
- **Repo / branch:** _none yet — to be created in A.1._
- **Last commit:** _n/a._

## ▶ Next action (the ONE thing to do next)
Scaffold the repo per `03-ARCHITECTURE.md`: Nuxt 3 `web/` + Go `api/` + PostgreSQL, the folder layout, and a `.env.example` matching `06-SETUP-AND-ENV.md`. Get both dev servers running and talking (a trivial `GET /health` from Go, fetched by Nuxt). **Copy `AGENTS.md`, `CLAUDE.md`, `HANDOFF-STATE.md`, `WORKLOG.md` to the new repo root.** Commit. Then stop and update this file — do not proceed past A.1 in one go.

## Done & verified
- _(nothing yet — planning only)_
- Handoff docs `00`–`06`, `catalogue/`, `design-source/`, `reference-docs/` are complete and current.

## In progress (exact state)
- _(nothing in progress)_

## How to run / verify current state
- No code yet. To view the design target: open files in `design-source/` in a browser (see `design-source/HOW-TO-VIEW.md`).
- Local dev commands (once scaffolded) are in `06-SETUP-AND-ENV.md`.

## Blockers / waiting on the user (client deliverables)
- Nomba credentials (sandbox first), WhatsApp business number, delivery-fee amounts per zone, product photography, customer-care page copy, domain DNS access. None block A.1 (scaffold) — but chase them; they gate later steps.

## Decisions / deviations from the numbered docs
- _(none yet — record any here with a one-line reason and flag the user)_

## Reminders (non-negotiable)
Money in kobo (int64, no floats) · `paid` only after server-side verify · guest checkout only in phase 1 · mobile-first · ghost admin stays unlinked/unmentioned to client · commits as `mikechibuzor`, no `Co-Authored-By`.
