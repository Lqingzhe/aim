package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketHub struct {
	Client map[int64]map[string]*Client
	Mu     sync.RWMutex
}
type Client struct {
	UserID           int64
	DeviceID         string
	Conn             *websocket.Conn
	Send             chan []byte
	IsConnected      bool
	LastLogoutSecond int64
}
