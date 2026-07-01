package context

import (
	"strings"

	"github.com/Brook-sys/auxitalk/pkg/types"
)

type Builder struct {
	MaxMessages int
}

func NewBuilder(maxMessages int) *Builder {
	if maxMessages <= 0 {
		maxMessages = 20
	}
	return &Builder{MaxMessages: maxMessages}
}

func (b *Builder) Build(session types.Session) string {
	var sb strings.Builder

	sb.WriteString("Session: ")
	sb.WriteString(session.ID)
	sb.WriteString(" (")
	sb.WriteString(session.Channel)
	sb.WriteString(")\n")

	messages := session.Messages
	if len(messages) > b.MaxMessages {
		messages = messages[len(messages)-b.MaxMessages:]
	}

	for _, msg := range messages {
		if msg.Direction == types.MessageDirectionInbound {
			sb.WriteString("User: ")
		} else {
			sb.WriteString("Assistant: ")
		}
		sb.WriteString(msg.Text)
		sb.WriteString("\n")
	}

	return sb.String()
}
