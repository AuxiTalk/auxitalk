# Core Types

AuxiTalk uses shared domain types to keep the core and plugins aligned.

The public Go definitions live in `pkg/types`.

## Event

Represents something that happened in the runtime.

Required fields:

- `id`
- `type`
- `source`
- `createdAt`

Common examples:

```txt
message.received
message.sent
session.updated
suggestion.created
action.requested
feedback.received
```

## Session

Represents one active or historical conversation.

Required fields:

- `id`
- `channel`
- `createdAt`
- `updatedAt`

A session may contain participants, messages, metadata, and state.

## Message

Represents a normalized conversation message.

Required fields:

- `id`
- `sessionId`
- `authorId`
- `direction`
- `source`
- `createdAt`

Supported directions:

```txt
inbound
outbound
```

## Suggestion

Represents a runtime recommendation.

Supported actions:

```txt
respond
wait
ask_user
ignore
```

For `respond`, `text` is required.

Confidence must be between `0` and `1`.

## Action request

Represents a request to perform an operation.

Risk levels:

```txt
low
medium
high
```

Statuses:

```txt
requested
allowed
confirmed
denied
executed
failed
```

High-risk actions, such as sending a message, must go through the Action Gate.

## Capability

Represents something a plugin can provide.

Examples:

```txt
conversation.observe
message.send
ai.complete
memory.query
memory.write
ui.suggestion.display
```

## Plugin manifest

Represents plugin metadata and declared permissions/capabilities.

Required fields:

- `id`
- `name`
- `version`
- `runtime`
- `entry`
- `kind`

Supported kinds:

```txt
input
output
ai
memory
ui
policy
tool
profile
```
