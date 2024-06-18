package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/moonen-home-automation/go-hasocket/internal"
	ws "github.com/moonen-home-automation/go-hasocket/internal/websocket"
	"log/slog"
)

type SubEvent struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

type UnSubEvent struct {
	Id           int64  `json:"id"`
	Type         string `json:"type"`
	Subscription int64  `json:"subscription"`
}

type EventData struct {
	Type         string
	RawEventJSON []byte
}

type BaseEventMsg struct {
	Event struct {
		EventType string `json:"event_type"`
	} `json:"event"`
}

func NewEventListener(ws *ws.Writer, ctx context.Context, eventType string) EventListener {
	return EventListener{ws: ws, ctx: ctx, eventType: eventType}
}

func (e *EventListener) Register() error {
	id := internal.GetId()
	e.subscriptionId = id
	sub := SubEvent{
		Id:        e.subscriptionId,
		Type:      "subscribe_events",
		EventType: e.eventType,
	}
	return e.ws.WriteMessage(sub, e.ctx)
}

func (e *EventListener) Close() error {
	id := internal.GetId()
	sub := UnSubEvent{
		Id:           id,
		Type:         "unsubscribe_events",
		Subscription: e.subscriptionId,
	}
	fmt.Println("Channel closed, unsubscribing")
	err := e.ws.WriteMessage(sub, e.ctx)
	if err != nil {
		return err
	}
	e.closed = true
	return nil
}

func (e *EventListener) Listen(dataChan chan EventData) {
	for {
		if e.closed {
			return
		}
		bytes, err := ws.ReadMessage(e.ws.Conn, e.ctx)
		if err != nil {
			slog.Error("Error reading from websocket:", err)
			continue
		}

		base := ws.BaseMessage{
			// default to true for messages that don't include "success" at all
			Success: true,
		}
		err = json.Unmarshal(bytes, &base)
		if err != nil {
			slog.Error("Error parsing websocket message:", err)
			continue
		}
		if !base.Success {
			slog.Warn("Received unsuccessful response", "response", string(bytes))
			continue
		}

		if base.Type != "event" {
			continue
		}
		baseEventMsg := BaseEventMsg{}
		err = json.Unmarshal(bytes, &baseEventMsg)
		if err != nil {
			continue
		}

		if baseEventMsg.Event.EventType != e.eventType {
			continue
		}
		eventData := EventData{
			Type:         baseEventMsg.Event.EventType,
			RawEventJSON: bytes,
		}
		dataChan <- eventData
	}
}
