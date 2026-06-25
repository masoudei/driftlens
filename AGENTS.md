# AGENTS.md

# DriftLens Agent Operating System

Version: 1.0

---

# Mission

DriftLens is a GitOps Intelligence Platform.

The goal is not merely detecting drift.

The goal is understanding drift.

DriftLens correlates:

* Git Commits
* CI/CD Pipelines
* Terraform Runs
* OpenTofu Runs
* ArgoCD Events
* Flux Events
* Kubernetes Audit Events
* OpenTelemetry Traces

into a unified causality graph that explains:

* Who changed production
* What changed
* Why it changed
* What was affected
* How risky it is
* How it should be fixed

Agents must optimize for explainability, causality, and operational intelligence.

---

# Required Reading Order

Before making any code changes, agents MUST read:

1. docs/vision.md
2. docs/market-gap.md
3. docs/design-principles.md
4. docs/architecture.md
5. docs/roadmap.md
6. docs/implementation-plan.md

These documents define project intent.

Code must follow docs.

Docs are the source of truth.

---

# Core Philosophy

DriftLens is not:

* Another GitOps controller
* Another dashboard
* Another Kubernetes UI
* Another Terraform wrapper

DriftLens is:

Infrastructure Change Intelligence.

Every feature should answer one or more:

* Why?
* Who?
* When?
* What?
* Impact?
* Risk?

If a feature does not improve those answers, reconsider implementing it.

---

# Architecture Rules

Always favor:

Event → Correlation → Graph → Insight

over:

Event → Alert

DriftLens should explain events, not merely report them.

All major components should be loosely coupled.

Preferred architecture:

Collector
→ Event Bus
→ Correlator
→ Graph Engine
→ Risk Engine
→ API
→ UI

Avoid tightly coupled services.

---

# Graphify First Rule

DriftLens uses Graphify thinking.

Everything is represented as nodes and relationships.

Examples:

Commit
└── CAUSED ──► Pipeline

Pipeline
└── TRIGGERED ──► TerraformRun

TerraformRun
└── DEPLOYED ──► ArgoSync

ArgoSync
└── MODIFIED ──► Deployment

Deployment
└── DRIFTED ──► DriftEvent

DriftEvent
└── IMPACTS ──► Service

Before introducing new storage models ask:

Can this be represented as a graph?

Prefer graph modeling whenever possible.

---

# Engineering Rules

## Rule 1

Keep modules small.

Target:

* <500 LOC per file
* <2000 LOC per package

---

## Rule 2

Favor composition over inheritance.

---

## Rule 3

Avoid hidden magic.

Explicit behavior preferred.

---

## Rule 4

Public interfaces first.

Design contracts before implementations.

---

## Rule 5

Every exported package must contain:

* package documentation
* examples
* tests

---

# Security Rules

Security is mandatory.

Never trade security for convenience.

## Secrets

Never:

* hardcode credentials
* hardcode tokens
* commit secrets

Use:

* environment variables
* secret managers
* Kubernetes secrets

---

## Auditability

All critical actions must be auditable.

Examples:

* user login
* settings changes
* risk score changes
* drift acknowledgement

---

## Multi-Tenancy

Design future features assuming:

* multiple organizations
* multiple teams
* multiple clusters

---

## Least Privilege

Collectors must request minimum RBAC permissions.

Avoid cluster-admin.

---

## Supply Chain

Require:

* dependency scanning
* SBOM generation
* image signing

---

# Testing Rules

Every new feature requires:

## Unit Tests

Required.

## Integration Tests

Preferred.

## End-to-End Tests

Required for major workflows.

Examples:

Git Commit
→ Argo Sync
→ Drift
→ Timeline

must be testable.

---

# Documentation Rules

Documentation is part of the product.

Every significant change should evaluate:

Do docs need updates?

Possible updates:

* architecture.md
* roadmap.md
* implementation-plan.md
* design-principles.md

Agents should never leave documentation stale.

---

# Git Workflow

DriftLens follows professional open-source Git practices.

## Branching Strategy

```
main          — production-ready, protected, linear history
├── feat/*    — new features (e.g. feat/graph-engine)
├── fix/*     — bug fixes (e.g. fix/pagination-offset)
├── refactor/*— code restructuring (e.g. refactor/event-bus)
├── docs/*    — documentation only (e.g. docs/api-reference)
├── test/*    — test additions/fixes (e.g. test/graph-traversal)
└── chore/*   — tooling, CI, deps (e.g. chore/update-deps)
```

Rules:

- Branch off `main`.
- Keep branches short-lived. Open a PR early.
- Delete branch after merge.
- Never commit directly to `main`.

## Commit Message Convention

Every commit must follow:

```
<type>(<scope>): <imperative subject>
<BLANK LINE>
<body (optional)>
<BLANK LINE>
<footer (optional)>
```

Types:

| Type       | Usage                              |
|------------|------------------------------------|
| `feat`     | New feature                        |
| `fix`      | Bug fix                            |
| `refactor` | Code restructure (no behavior change) |
| `docs`     | Documentation changes              |
| `test`     | Adding or fixing tests             |
| `build`    | Build system, CI, dependencies     |
| `chore`    | Tooling, config, maintenance       |
| `perf`     | Performance improvement            |
| `style`    | Formatting (no logic change)       |

Scopes (examples): `graph`, `api`, `ui`, `collector`, `engine`, `cli`, `docs`

Rules:

- **Atomic commits**: each commit is one logical change. Do not mix unrelated changes (e.g. a bug fix + a refactor + a docs update) in a single commit. Split into separate commits.
- Subject is lowercase, imperative, no period.
- Maximum 72 characters for subject.
- Body explains *what* and *why*, not *how*.
- Footer for breaking changes (`BREAKING CHANGE:`) or issue refs (`Closes #123`).
- Never use `-m` for multi-line commits — write full message.

Good examples:

```
feat(graph): add DFS traversal for causality chains

Implements depth-first traversal from any node to discover
full causality paths. Supports max-depth limiting to prevent
infinite loops in cyclic graphs.

Closes #42
```

```
fix(api): return 404 for unknown drift IDs

Previously returned empty 200 which broke client error handling.

Closes #57
```

Bad examples:

```
fix stuff                            ← vague, no scope
feat(api): added new endpoint        ← past tense, not imperative
chore: fixed lint                    ← should be style or fix
```

## Pull Request Workflow

When the user says "create PR" or "open PR":

1. Push the current branch to origin
2. Ensure the branch name matches the convention above
3. Use `gh pr create` with:

   - Title matching the commit convention (type(scope): subject)
   - Body containing:
     - **What** — summary of changes
     - **Why** — motivation / linked issue
     - **How** — high-level approach
     - **Checklist** — from DriftLens Quality Gate below
4. Request review from relevant maintainers

## Commit And Push Workflow

Do not commit or push automatically after each change. Wait for the user to explicitly say "commit" or "commit and push".

When the user says:

"commit"

or

"commit and push"

the agent must perform the following review loop before committing:

1. Review changed files with `git status` and `git diff`
2. Review AGENTS.md
3. Review docs folder
4. Review README.md — does the feature change the architecture, graph model, tech stack, or getting started flow?
5. Determine whether:

   * architecture changed
   * roadmap changed
   * implementation changed
   * design principles changed
6. Update documentation if required — including README.md if the change affects users, architecture, or setup
7. Update AGENTS.md if a new rule emerged
8. **Run tests** — write required tests for any new feature or change, then run the full test suite and confirm all tests pass:

   `go test ./...`

9. **Build and run locally** — verify the application compiles cleanly and, where applicable, runs without error:

   `go build ./...`

10. Stage only intended files — never stage secrets, binaries, or generated files
11. Write commit message following the convention above
12. Commit
13. Push if requested

Never commit stale documentation intentionally.

Never commit failing tests.

Never commit code that hasn't been built and verified locally.

## Push Rules

- Always verify remote is reachable before pushing: `git remote -v`
- Use `git push origin <branch>` — never force push unless explicitly told
- If `main` is ahead, rebase first: `git pull --rebase origin main`
- Never push secrets, tokens, or credentials
- If CI is configured, wait for CI to pass before merging

---

# Self Improvement Loop

After completing any significant task:

Run Reflection Pass.

Questions:

1. What was implemented?
2. Why was it implemented?
3. What architecture changed?
4. What assumptions were introduced?
5. What documentation should change?
6. What future work became obvious?
7. Should AGENTS.md evolve?

If yes:

Update docs.

Update AGENTS.md.

Then proceed.

---

# Recursive Learning Loop

AGENTS.md is a living document.

Every task cycle follows: **Learn → Implement → Test → Think → Evolve**

After each completed task, the agent MUST:

1. **Learn** — Analyze what was built. What patterns emerged? What was hard?
2. **Implement** — Did the code match the docs? Does the architecture doc need updating?
3. **Test** — Did tests pass? Did lint pass? Was the quality gate satisfied?
4. **Think** — Run the Reflection Pass (see Self Improvement Loop above)
5. **Evolve** — Update AGENTS.md with any new rule, convention, or pattern discovered during the task

This loop is recursive.

After updating AGENTS.md, future agent sessions inherit the new rules.

Over time, AGENTS.md becomes the cumulative brain of the project.

Examples of what to add to AGENTS.md:

- A naming convention that was ambiguous but is now settled
- A CLI flag pattern that should be used consistently
- A test pattern that worked well
- A Go package layout decision
- An API design rule discovered during implementation

If the task reveals nothing new worth encoding, skip the update.

But prefer to err on the side of documenting.

---

# Architecture Review Loop

For every feature ask:

Can this support:

* 10 clusters?
* 100 clusters?
* 1000 clusters?

If not:

Document limitations.

Do not silently create scaling debt.

---

# DriftLens Quality Gate

Before marking work complete:

Checklist:

[ ] Tests pass

[ ] Lint passes

[ ] Documentation updated

[ ] Architecture remains consistent

[ ] Security reviewed

[ ] No secrets introduced

[ ] Graph relationships documented

[ ] APIs documented

[ ] Observability added

[ ] Commit message prepared

---

# North Star

Every completed feature should move DriftLens closer to answering:

Who changed production?

Why did they change it?

What was impacted?

How risky was it?

How do we fix it?

If the answer quality improves, the feature is valuable.

If not, rethink the implementation.

End of AGENTS.md
