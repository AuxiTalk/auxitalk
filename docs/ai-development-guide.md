# AI Development Guide

This guide explains how AI coding agents should collaborate on AuxiTalk Core.

## Goal

Help build AuxiTalk as a safe, event-driven automation runtime.

AI agents should preserve the long-term architecture:

```txt
Event -> Workflow -> Action Request -> Gate -> Executor / Plugin -> Result Event
```

## Standard workflow

1. Read the relevant docs and code.
2. Make a short plan.
3. Implement a small slice.
4. Add or update tests.
5. Run `gofmt` for changed Go files.
6. Run `go test ./...`.
7. Summarize changed files and results.
8. Commit only when asked.
9. Push only when explicitly asked.

## Investigation checklist

Before editing, identify:

- which package owns the behavior;
- which public type is affected;
- whether runtime integration is needed;
- whether plugin protocol changes are needed;
- whether the change can cause real side effects;
- which tests should cover the behavior.

## Safe implementation pattern

Prefer this order:

1. Types and validation.
2. In-memory implementation.
3. Mock/dry-run executor.
4. Unit tests.
5. Runtime integration.
6. Real side effects behind gates and explicit approval.

## Areas of caution

### Plugin supervisor

Files:

- `internal/plugins/supervisor/process.go`
- `internal/plugins/supervisor/supervisor.go`
- `pkg/protocol/jsonrpc.go`

Rules:

- Preserve line-delimited JSON-RPC over stdio.
- Do not write logs to plugin stdout.
- Keep pending request cleanup on timeout.
- Keep payload limits.
- Test with fake plugins.

### Runtime

Files:

- `internal/runtime/runtime.go`

Rules:

- Runtime coordinates modules; avoid putting all business logic here.
- Emit events for observable state changes.
- Keep shutdown graceful.
- Avoid blocking event processing indefinitely.

### Workflows

Files:

- `internal/workflows/*`
- `pkg/types/workflow.go`
- `pkg/types/action_execution.go`

Rules:

- Keep workflow rules deterministic at first.
- Prefer dry-run executors until real side effects are approved.
- Generated actions should go through the action gate/store.

### Actions

Files:

- `internal/actions/*`
- `pkg/types/action.go`

Rules:

- Risky operations must be represented as `ActionRequest`.
- Respect runtime modes: `dev`, `local`, `strict`.
- Keep approvals auditable.

## Documentation expectations

When behavior changes product direction or developer workflow, update one of:

- `README.md`
- `docs/product/vision.md`
- `docs/workflow-engine.md`
- package-specific docs under `docs/architecture` or `docs/plugins`

## Good task slices

Good:

- Add workflow config types and tests.
- Add a mock executor behavior.
- Wire event bus to workflow engine with tests.
- Add plugin status endpoint data structures.

Too broad:

- Build the full dashboard and all integrations.
- Add real shell execution without policy gates.
- Rewrite the plugin protocol without migration plan.

## Final response format

When finishing a task, report:

- summary;
- files changed;
- tests run;
- commit hash if committed;
- blockers if any.
