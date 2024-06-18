package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
)

type BaseMessage struct {
	Type    string `json:"type"`
	Id      int64  `json:"id"`
	Success bool   `json:"success"`
}

type ChanMsg struct {
	Id      int64
	Type    string
	Success bool
	Raw     []byte
}

func ReadMessage(conn *websocket.Conn, ctx context.Context) ([]byte, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

func ListenWebsocket(conn *websocket.Conn, ctx context.Context, c chan ChanMsg) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Function done")
			return
		default:
			bytes, err := ReadMessage(conn, ctx)
			if err != nil {
				slog.Error("Error reading from websocket:", err)
				return
			}

			base := BaseMessage{
				// default to true for messages that don't include "success" at all
				Success: true,
			}
			err = json.Unmarshal(bytes, &base)
			if err != nil {
				slog.Error("Error parsing websocket message:", err)
				return
			}
			if !base.Success {
				slog.Warn("Received unsuccessful response", "response", string(bytes))
				return
			}
			chanMsg := ChanMsg{
				Type:    base.Type,
				Id:      base.Id,
				Success: base.Success,
				Raw:     bytes,
			}

			select {
			case c <- chanMsg:
				fmt.Println("Message send to channel")
			case <-ctx.Done():
				fmt.Println("Function done while sending message")
				return
			}
		}
	}
}
