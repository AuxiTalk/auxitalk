# AGENTS.md

This repository is the AuxiTalk Core runtime.

AuxiTalk is an event-driven automation runtime for connecting people, AI agents, tools, plugins, communication channels, terminals, APIs, and dashboards through safe, observable workflows.

## Primary context

Read these files first:

1. `docs/product/vision.md` — product purpose and use cases
2. `README.md` — project overview
3. `docs/workflow-engine.md` — workflow model
4. `docs/plugins/protocol-draft.md` — plugin protocol
5. `docs/plugins/supervisor.md` — plugin supervisor
6. `docs/architecture/core.md` — core architecture

## Required commands

Before finishing code changes, run:

```sh
go test ./...
```

If Go files were edited, run:

```sh
gofmt -w <changed-go-files>
```

Do not commit unless explicitly asked by the user.
Do not push unless explicitly asked by the user.

## Current architecture

Important packages:

- `cmd/auxitalkd` — daemon entrypoint
- `cmd/auxitalkctl` — CLI placeholder
- `internal/runtime` — runtime lifecycle and integration layer
- `internal/events` — internal event bus
- `internal/actions` — action gate and action store
- `internal/workflows` — workflow engine, registry, mock executor
- `internal/plugins` — plugin manifest and registry
- `internal/plugins/supervisor` — external plugin process supervision and JSON-RPC
- `internal/capabilities` — capability routing
- `internal/sessions` — session state
- `internal/context` — context building
- `pkg/types` — public domain types
- `pkg/protocol` — JSON-RPC protocol types

## Product direction

Do not treat AuxiTalk as only a chatbot.

It should support:

- chat/message automation;
- terminal copilots;
- AI agent orchestration;
- event/webhook automation;
- human-in-the-loop approvals;
- dashboard and CLI control surfaces;
- local-first workflow automation.

## Design principles

- Event-driven core.
- Plugin-first integrations.
- Safe by default.
- Local-first when possible.
- Observable and auditable.
- Multi-interface.
- AI-assisted, not AI-only.

## Safety rules

- Never introduce code that logs secrets.
- Never store API keys in committed files.
- Risky actions must go through action requests/gates.
- Real side effects must be explicit, testable, and policy-controlled.
- Prefer dry-run/mock executors until a real executor is approved.
- Plugin stdout is reserved for JSON-RPC messages.
- Plugin stderr is for human-readable logs.

## Coding conventions

- Keep the core small and readable.
- Prefer simple Go over clever abstractions.
- Add tests for new behavior.
- Keep public contracts in `pkg/types` or `pkg/protocol`.
- Keep runtime implementation details under `internal`.
- Do not add comments unless they clarify non-obvious behavior.
- Prefer small, incremental commits.

## Common workflows

### Adding a new core type

1. Add type in `pkg/types`.
2. Add validation method if needed.
3. Add unit tests.
4. Update docs if it changes product behavior.
5. Run `go test ./...`.

### Adding runtime behavior

1. Inspect `internal/runtime/runtime.go`.
2. Check event/action/workflow interactions.
3. Add tests in the closest package.
4. Avoid side effects unless explicitly approved.
5. Run `go test ./...`.

### Adding plugin protocol behavior

1. Inspect `pkg/protocol` and `internal/plugins/supervisor`.
2. Preserve JSON-RPC 2.0 line-delimited stdio.
3. Keep timeout and payload limits.
4. Add fake plugin tests.
5. Run `go test ./...`.

## Known priorities

- Wire workflow engine into runtime event handling.
- Load workflows from config.
- Improve dashboard with real runtime data.
- Test WhatsApp plugin through core.
- Add SQLite persistence for events/actions/workflows.
