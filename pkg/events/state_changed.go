package events

import (
	"context"
	ws "github.com/moonen-home-automation/go-hasocket/internal/websocket"
	"time"
)

type StateChangedEventData struct {
	EntityId string            `json:"entity_id"`
	OldState StateChangedState `json:"old_state"`
	NewState StateChangedState `json:"new_state"`
}

type StateChangedState struct {
	EntityId     string         `json:"entity_id"`
	State        string         `json:"state"`
	Attributes   map[string]any `json:"attributes"`
	LastChanged  time.Time      `json:"last_changed"`
	LastReported time.Time      `json:"last_reported"`
	LastUpdated  time.Time      `json:"last_updated"`
}

type EventListener struct {
	ws  *ws.Writer
	ctx context.Context

	eventType      string
	subscriptionId int64
	closed         bool
}
