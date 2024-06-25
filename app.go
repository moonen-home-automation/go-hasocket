package go_hasocket

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	ws "github.com/moonen-home-automation/go-hasocket/internal/websocket"
	events "github.com/moonen-home-automation/go-hasocket/pkg/events"
	services "github.com/moonen-home-automation/go-hasocket/pkg/services"
)

var ErrInvalidArgs = errors.New("invalid arguments provided")

var appInstance *App

type App struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	conn      *websocket.Conn

	wsWriter *ws.Writer

	serviceCaller *services.ServiceCaller

	eventListeners map[string][]*events.EventListener
}

func NewApp(uri, token string) (*App, error) {
	ctx := context.Background()
	conn, err := ws.ConnectionFromURI(ctx, uri, token)
	if conn == nil {
		return nil, err
	}

	wsWriter := &ws.Writer{Conn: conn}
	go wsWriter.KeepAlive(ctx)

	serviceCaller := services.NewServiceCaller(wsWriter, ctx)

	appInstance = &App{
		conn:          conn,
		wsWriter:      wsWriter,
		ctx:           ctx,
		serviceCaller: &serviceCaller,
	}

	return appInstance, nil
}

func GetApp() (*App, error) {
	if appInstance == nil {
		return appInstance, errors.New("app not defined")
	}
	return appInstance, nil
}

func (a *App) CallService(sr services.ServiceRequest) (services.ServiceResult, error) {
	return a.serviceCaller.Call(sr)
}

func (a *App) RegisterListener(eventType string) (events.EventListener, error) {
	listener := events.NewEventListener(a.wsWriter, a.ctx, eventType)

	err := listener.Register()
	return listener, err
}
