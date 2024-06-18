package services

import (
	"context"
	"encoding/json"
	ws "github.com/moonen-home-automation/go-hasocket/internal/websocket"
	"log/slog"
)

type ServiceCaller struct {
	ws  *ws.Writer
	ctx context.Context
}

func NewServiceCaller(ws *ws.Writer, ctx context.Context) ServiceCaller {
	return ServiceCaller{ws: ws, ctx: ctx}
}

func (s *ServiceCaller) Call(req ServiceRequest) (ServiceResult, error) {
	err := s.ws.WriteMessage(req, s.ctx)
	if err != nil {
		return ServiceResult{}, err
	}

	if req.ReturnResponse {
		return s.awaitResponse(req.Id)
	}
	return ServiceResult{}, nil
}

func (s *ServiceCaller) awaitResponse(id int64) (ServiceResult, error) {
	ctx, ctxCancel := context.WithCancel(s.ctx)
	defer ctxCancel()

	for {
		bytes, err := ws.ReadMessage(s.ws.Conn, ctx)
		if err != nil {
			slog.Error("Error reading from websocket:", err)
			return ServiceResult{}, err
		}

		base := ws.BaseMessage{
			// default to true for messages that don't include "success" at all
			Success: true,
		}
		err = json.Unmarshal(bytes, &base)
		if err != nil {
			slog.Error("Error parsing websocket message:", err)
			return ServiceResult{}, err
		}
		if !base.Success {
			slog.Warn("Received unsuccessful response", "response", string(bytes))
			return ServiceResult{}, err
		}

		if base.Type != "result" {
			continue
		}
		serviceResults := ServiceResult{}
		err = json.Unmarshal(bytes, &serviceResults)
		if err != nil {
			return ServiceResult{}, err
		}

		if serviceResults.Id != id {
			continue
		}

		return serviceResults, nil
	}
}
