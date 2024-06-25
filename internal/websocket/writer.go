package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/moonen-home-automation/go-hasocket/internal"
	"sync"
	"time"
)

type Writer struct {
	Conn  *websocket.Conn
	mutex sync.Mutex
}

func (w *Writer) WriteMessage(msg interface{}, ctx context.Context) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	err := w.Conn.WriteJSON(msg)
	if err != nil {
		return err
	}

	return nil
}

type pingMsg struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
}

func (w *Writer) KeepAlive(ctx context.Context) {
	for {
		id := internal.GetId()
		err := w.WriteMessage(pingMsg{Id: id, Type: "ping"}, ctx)
		if err != nil {
			continue
		}
		time.Sleep(time.Second * 10)
	}
}
