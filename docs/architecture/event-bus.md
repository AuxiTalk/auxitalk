# Event Bus

AuxiTalk uses an internal event bus to connect core modules without hard dependencies between them.

Implementation lives in `internal/events`.

## Responsibilities

- Publish validated runtime events.
- Subscribe handlers by event type.
- Subscribe wildcard handlers to all events.
- Respect context cancellation.
- Enforce per-handler timeout.
- Return handler errors to the publisher.
- Keep optional in-memory event history.

## Event validation

Every published event must pass `types.Event.Validate()`.

Required fields:

- `id`
- `type`
- `source`
- `createdAt`

## Subscription modes

Typed subscription:

```go
bus.Subscribe("message.received", handler)
```

Wildcard subscription:

```go
bus.Subscribe("*", handler)
```

Empty event type also means wildcard subscription.

## Handler timeout

Every handler receives a context with timeout.

If a handler does not return before the timeout, publish returns an error.

## History

The bus can keep the last N events in memory.

```go
bus := events.New(events.Options{
    HandlerTimeout: 5 * time.Second,
    HistoryLimit: 100,
})
```

If `HistoryLimit` is zero or negative, history is disabled.

## Design notes

The event bus is intentionally internal. Plugins do not receive direct access to it.

Plugins emit events through the JSON-RPC protocol, and the core validates and publishes those events internally.
