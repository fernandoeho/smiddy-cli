# Smiddy CLI

**Spec-driven, autonomous development sprints — powered by Claude.**

Smiddy CLI is a CLI tool that turns a plain-text spec into working code through an iterative, multi-agent loop. You write what you want built. Smiddy sends it to an Architect agent that plans the work, an Implementation agent that writes the code, and a Review agent that checks the result — repeating until the spec is met or the iteration budget runs out.

```
smiddy init → smiddy setup → smiddy new → (edit specs.md) → smiddy run
```

---

## How it works

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  specs.md   │────▶│  Architect Agent │────▶│ Implementation  │
│  (you write)│     │  (plans the work)│     │     Agent       │
└─────────────┘     └──────────────────┘     │  (writes code)  │
                              ▲               └────────┬────────┘
                              │                        │
                    ┌─────────┴────────┐               │
                    │  Review Agent    │◀──────────────┘
                    │ COMPLETE / RETRY │
                    └──────────────────┘
```

Each sprint runs up to **5 iterations**. On every pass, the Architect reviews the implementation against your original spec and either marks it `COMPLETE` or refines the plan for the next iteration. All artifacts — plans, implementation reports, and review notes — are saved in `.smiddy/sprints/<n>/`.

---

## Requirements

- Go 1.21+
- An [Anthropic API key](https://console.anthropic.com/)

---

## Installation

### From source

```bash
git clone https://github.com/fernandoeho/smiddy.git
cd smiddy
go build -o smiddy .
# move to your PATH
mv smiddy /usr/local/bin/smiddy
```

### Set your API key

```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

Add that line to your shell profile (`~/.zshrc`, `~/.bashrc`) to make it permanent.

---

## Quickstart

```bash
# 1. Initialize a new Smiddy project in your repo
smiddy init

# 2. Describe your project (saves .smiddy/project-goals.md)
smiddy setup

# 3. Create a sprint and write the spec
smiddy new
# → opens .smiddy/sprints/1/specs.md — fill in what you want built

# 4. Run the sprint
smiddy run

# 5. Check progress
smiddy status
```

---

## Commands

| Command | Description |
|---|---|
| `smiddy init` | Initialize a `.smiddy/` workspace in the current directory |
| `smiddy setup` | Interactive onboarding — define project vision, audience, and constraints |
| `smiddy new` | Create a new sprint and generate a blank `specs.md` |
| `smiddy run` | Execute the current sprint through the full agent loop |
| `smiddy status` | Print the current sprint's state and iteration count |
| `smiddy clean` | Delete old sprint directories (keeps the latest) |

---

## Writing a spec

After `smiddy new`, edit `.smiddy/sprints/<n>/specs.md`. A good spec is concrete:

```markdown
# Sprint 3

## Goal
Add a REST endpoint POST /users that creates a new user in the database.

## Acceptance criteria
- Validates that email and name are present
- Returns 201 with the created user object on success
- Returns 422 with a JSON error body on validation failure
- Writes an integration test covering both cases
```

The more precise your spec, the fewer iterations the agent needs.

---

## Project workspace

```
.smiddy/
├── project-goals.md        # Vision, audience, constraints (smiddy setup)
├── project-map.md          # Architecture and key dependencies (edit manually)
└── sprints/
    └── 1/
        ├── specs.md        # Your input — what to build
        ├── architect-plan.md
        ├── implementation-1.md
        ├── review-1.md
        └── status.md
```

`project-goals.md` and `project-map.md` are injected into every sprint prompt as context. Keep them up to date as your project evolves.

---

## Agents

Smiddy uses three Claude-powered roles on every sprint:

**Architect** — reads your spec and project context, breaks the work into numbered tasks, produces a structured plan.

**Implementation Agent** — takes the architect's plan and writes the code, reporting every file touched.

**Review Agent** — compares the implementation report against the original spec. If everything is covered, it emits `COMPLETE`. Otherwise it refines the plan and hands it back to the implementation agent for another pass.

---

## License

[MIT](LICENSE)
