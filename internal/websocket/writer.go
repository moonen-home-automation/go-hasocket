package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"sync"
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
