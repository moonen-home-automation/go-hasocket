// Package websocket /*
package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log/slog"
)

func ConnectionFromURI(ctx context.Context, uri, authToken string) (*websocket.Conn, error) {
	// Init websocket connection
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, uri, nil)
	if err != nil {
		slog.Error("Failed to connect to websocket, Check URI\n", "uri", uri)
		return nil, err
	}

	// Read auth_required
	_, err = ReadMessage(conn, ctx)
	if err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Send auth message
	err = sendAuthMessage(conn, ctx, authToken)
	if err != nil {
		slog.Error("Unknown error creating websocket client\n")
		return nil, err
	}

	// Verify auth response
	err = verifyAuthResponse(conn, ctx)
	if err != nil {
		slog.Error("Auth token is invalid. Please double check it or create a new token in your Home Assistant profile\n")
		return nil, err
	}

	return conn, nil
}

type authMessage struct {
	MsgType     string `json:"type"`
	AccessToken string `json:"access_token"`
}

func sendAuthMessage(conn *websocket.Conn, ctx context.Context, token string) error {
	err := conn.WriteJSON(authMessage{MsgType: "auth", AccessToken: token})
	if err != nil {
		return err
	}
	return nil
}

var errInvalidToken = errors.New("invalid auth token")

type authResponse struct {
	MsgType string `json:"type"`
	Message string `json:"message"`
}

func verifyAuthResponse(conn *websocket.Conn, ctx context.Context) error {
	msg, err := ReadMessage(conn, ctx)
	if err != nil {
		return err
	}

	var authResp authResponse
	_ = json.Unmarshal(msg, &authResp)
	if authResp.MsgType != "auth_ok" {
		return errInvalidToken
	}

	return nil
}
