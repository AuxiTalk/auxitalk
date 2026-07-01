# Session Manager and Context Builder

## Session Manager (`internal/sessions`)

Responsibilities:

- Create, get, update, delete sessions.
- Append normalized messages to sessions.
- Keep `UpdatedAt` current.
- Provide list of active sessions.

## Context Builder (`internal/context`)

Responsibilities:

- Build compact conversation context from a session.
- Truncate older messages when exceeding limit.
- Format messages with direction labels (`User:` / `Assistant:`).

The context builder is intentionally simple. More advanced reduction strategies (summarization, embeddings, token counting) can be added later without changing the interface.
