# AuxiTalk Plugin WhatsApp — Planning Document

> Robust planning document for the official WhatsApp plugin using whatsmeow.

---

## 1. Goal

Create an official WhatsApp plugin for AuxiTalk using `go.mau.fi/whatsmeow`.

The plugin should connect to WhatsApp through QR Code login, receive messages in real time, normalize them into AuxiTalk events, and send messages only through approved core actions.

---

## 2. Repository

```txt
AuxiTalk/plugin-whatsapp
```

## 3. Language and library

- Language: Go
- WhatsApp client: `go.mau.fi/whatsmeow`
- Session storage: SQLite using whatsmeow store
- Protocol with AuxiTalk Core: JSON-RPC 2.0 over stdio

---

## 4. Plugin kind

This plugin is both input and output.

Initial manifest kind decision:

```txt
kind: "input"
```

Future core may support combined kinds or tags. Until then, output capabilities are declared through capabilities.

---

## 5. Main responsibilities

### Must do

- Connect to WhatsApp using QR Code.
- Persist session locally.
- Reconnect automatically when possible.
- Receive text messages in real time.
- Normalize received messages into AuxiTalk events.
- Emit QR/status events for dashboard.
- Send messages through `message.send` capability.
- Log only to stderr.
- Never send messages without a core capability/action path.

### Must not do initially

- Read screen.
- Use OCR.
- Automate browser.
- Send messages outside approved capability call.
- Implement AI logic.
- Store conversation memory beyond WhatsApp session state.

---

## 6. MVP scope

### Included in v0.1

- QR Code login.
- Session persistence.
- Connect/disconnect status events.
- Receive text messages.
- Emit `message.received`.
- Send text messages through `message.send`.
- Basic health check.
- README + README.pt-BR.
- Configuration docs.
- Security docs.

### Not included in v0.1

- Media messages.
- Reactions.
- Groups advanced handling.
- Message edit/delete sync.
- Contact sync UI.
- Multi-account support.
- Voice/audio transcription.
- End-to-end moderation.

---

## 7. Events emitted

### `whatsapp.qr`

Emitted when login QR is available.

Payload:

```json
{
  "code": "qr-code-string",
  "expiresAt": "timestamp"
}
```

### `whatsapp.connected`

Payload:

```json
{
  "user": "whatsapp-jid"
}
```

### `whatsapp.disconnected`

Payload:

```json
{
  "reason": "string"
}
```

### `message.received`

Payload:

```json
{
  "messageId": "string",
  "sessionId": "whatsapp:<chat-jid>",
  "channel": "whatsapp",
  "chatId": "jid",
  "senderId": "jid",
  "senderName": "optional",
  "text": "message text",
  "timestamp": "timestamp",
  "isGroup": false
}
```

### `conversation.updated`

Optional in MVP.

Payload:

```json
{
  "sessionId": "whatsapp:<chat-jid>",
  "channel": "whatsapp",
  "lastMessageId": "string"
}
```

---

## 8. Capabilities

### `message.send`

Sends a WhatsApp text message.

Input:

```json
{
  "chatId": "jid",
  "text": "message text"
}
```

Output:

```json
{
  "messageId": "string",
  "chatId": "jid",
  "sent": true
}
```

Validation:

- `chatId` required.
- `text` required.
- text must not be empty.

### Future capabilities

```txt
conversation.list
message.react
message.download_media
message.mark_read
whatsapp.status
```

---

## 9. Configuration

Environment variables:

```txt
WHATSAPP_DB_PATH=./whatsapp.db
WHATSAPP_DEVICE_NAME=AuxiTalk
WHATSAPP_AUTO_RECONNECT=true
WHATSAPP_QR_PRINT=true
WHATSAPP_LOG_LEVEL=info
```

Defaults should be safe and local.

---

## 10. Storage

The plugin stores only WhatsApp session/device data locally.

Default:

```txt
./whatsapp.db
```

This file must be ignored by git.

The plugin must not store conversation memory in this DB. Conversation memory belongs to memory plugins.

---

## 11. Security

Important notes:

- This uses WhatsApp Web protocol through whatsmeow.
- User must explicitly scan QR Code.
- Session DB gives access to the WhatsApp account and must be protected.
- Never commit session DB.
- Sending messages should only happen through `message.send` capability.
- Dashboard/core Action Gate should approve high-risk sends.
- Logs must not expose full message contents when privacy mode is enabled.

Future privacy setting:

```txt
WHATSAPP_PRIVACY_LOGS=true
```

When enabled, logs should avoid raw message text.

---

## 12. Runtime flow

```txt
plugin starts
  -> loads config
  -> opens WhatsApp session store
  -> checks existing session
  -> if no session: emits whatsapp.qr and prints QR
  -> connects
  -> listens for messages
  -> emits message.received
  -> handles capability.call(message.send)
```

---

## 13. JSON-RPC methods

Core -> plugin:

```txt
plugin.handshake
plugin.start
plugin.stop
plugin.health
capability.call
```

Plugin -> core:

```txt
event.emit
log.write
```

---

## 14. Repository structure

```txt
plugin-whatsapp/
  README.md
  README.pt-BR.md
  plugin.json
  go.mod
  go.sum
  .gitignore
  .env.example

  cmd/
    plugin/
      main.go

  internal/
    config/
      config.go
    plugin/
      rpc.go
      runtime.go
    whatsapp/
      client.go
      events.go
      types.go

  docs/
    configuration.md
    capabilities.md
    events.md
    security.md
```

---

## 15. Implementation phases

### W1 — Foundation

- Create repository.
- Add manifest.
- Add Go module.
- Add config loader.
- Add JSON-RPC runtime skeleton.

### W2 — WhatsApp connection

- Add whatsmeow.
- Add SQLite session store.
- Add QR login flow.
- Add connect/disconnect events.

### W3 — Receiving messages

- Listen for incoming messages.
- Normalize text messages.
- Emit `message.received`.

### W4 — Sending messages

- Implement `message.send` capability.
- Validate input.
- Return message ID.

### W5 — Docs and hardening

- Add README and PT-BR docs.
- Add security docs.
- Add tests where practical.
- Add manual test guide.

---

## 16. Acceptance criteria for v0.1

- Builds as a single Go binary.
- Starts as JSON-RPC plugin.
- Shows QR Code on first login.
- Persists session locally.
- Reconnects without QR after login.
- Emits `message.received` for text messages.
- Sends text through `message.send`.
- Has documentation for setup and security.

---

## 17. Open questions

- Should QR be printed as terminal QR, emitted only to dashboard, or both?
- Should group messages be enabled by default?
- Should outgoing messages from the user also emit events?
- Should privacy log mode default to true?

---

## 18. Initial decision proposal

- QR: both terminal print and event emit.
- Groups: receive but mark `isGroup=true`.
- User outgoing messages: emit later, not v0.1.
- Privacy logs: default true.

---

**Status:** planned, pending implementation.
