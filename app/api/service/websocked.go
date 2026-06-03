package service

import (
	"aim/app/api/model"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

type WebSocketStruct struct {
	hub *model.WebSockedHub
}

func NewWebSocket(hub *model.WebSockedHub) *WebSocketStruct {
	return &WebSocketStruct{
		hub: hub,
	}
}
func (w *WebSocketStruct) ConnectWebsocket(coon *websocket.Conn, userID int64, deviceID string) {
	client := &model.Client{
		UserID:      userID,
		DeviceID:    deviceID,
		Conn:        coon,
		Send:        make(chan []byte, 256),
		IsConnected: true,
	}
	w.Register(client)
	go w.WritePump(userID, deviceID)
	go w.ReadPump(userID, deviceID)
}
func (w *WebSocketStruct) PushToUser(userID int64, deviceID string, data []byte) (exist bool) {
	w.hub.Mu.RLock()
	clients, ok := w.hub.Client[userID]
	w.hub.Mu.RUnlock()
	if !ok {
		return false
	}
	w.hub.Mu.RLock()
	client, ok := clients[deviceID]
	w.hub.Mu.RUnlock()
	if !ok {
		return false
	}
	if !client.IsConnected {
		return false
	}
	select {
	case client.Send <- data:
		return true
	default:
		w.Unregister(userID, deviceID)
		return false
	}
}
func (w *WebSocketStruct) ReadPump(userID int64, deviceID string) {
	c, ok := w.hub.Client[userID][deviceID]
	if !ok {
		return
	}
	defer func() {
		w.Unregister(userID, deviceID)
	}()
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg struct {
			Type string `json:"type"`
		}
		if err := sonic.Unmarshal(msg, &wsMsg); err != nil {
			continue
		}

		switch wsMsg.Type {
		case "ping":
			c.Send <- []byte(`{"type":"pong"}`)
		case "logout":
			return // 退出循环，触发 defer 清理
		}
	}
}

func (w *WebSocketStruct) WritePump(userID int64, deviceID string) {
	c, ok := w.hub.Client[userID][deviceID]
	if !ok {
		return
	}
	defer func() {
		w.Unregister(userID, deviceID)
	}()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
func (w *WebSocketStruct) Register(c *model.Client) {
	w.hub.Mu.Lock()
	defer w.hub.Mu.Unlock()
	if w.hub.Client[c.UserID] == nil {
		w.hub.Client[c.UserID] = make(map[string]*model.Client)
	}
	w.hub.Client[c.UserID][c.DeviceID] = c

}
func (w *WebSocketStruct) Unregister(userID int64, deviceID string) {
	w.hub.Mu.RLock()
	_, ok := w.hub.Client[userID]
	w.hub.Mu.RUnlock()
	if !ok {
		return
	}
	w.hub.Mu.Lock()
	defer w.hub.Mu.Unlock()
	if deviceID == "" {
		for _, i := range w.hub.Client[userID] {
			i.IsConnected = false
			i.LastLogoutSecond = time.Now().Unix()
			close(i.Send)
			i.Conn.Close()
		}
	} else {
		i, ok := w.hub.Client[userID][deviceID]
		if !ok {
			return
		}
		if !i.IsConnected {
			return
		}
		i.IsConnected = false
		i.LastLogoutSecond = time.Now().Unix()
		close(i.Send)
		i.Conn.Close()
	}
}
func (w *WebSocketStruct) ClearClientOnTime() {
	go func() {
		for {
			now := time.Now()
			needClearTimeOfKey := now.Unix() - 7*24*60*60
			clearTime := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
			if now.After(clearTime) {
				clearTime = clearTime.AddDate(0, 0, 1)
			}
			select {
			case <-time.After(clearTime.Sub(now)):
			}
			type UserAndDeviceID struct {
				UserID   int64
				DeviceID string
			}
			needDeleteDeviceOfUserList := make([]UserAndDeviceID, 0)
			needDeleteUserIDList := make([]int64, 0)
			w.hub.Mu.RLock()
			for userID, deviceClient := range w.hub.Client {
				needDeleteUserID := true
				needDeleteDeviceList := make([]string, 0, len(w.hub.Client[userID]))
				for deviceID, client := range deviceClient {
					if !client.IsConnected && client.LastLogoutSecond < needClearTimeOfKey {
						needDeleteDeviceList = append(needDeleteDeviceList, deviceID)
					} else {
						needDeleteUserID = false
					}
				}
				for _, deviceID := range needDeleteDeviceList {
					needDeleteDeviceOfUserList = append(needDeleteDeviceOfUserList, UserAndDeviceID{UserID: userID, DeviceID: deviceID})
				}
				if needDeleteUserID {
					needDeleteUserIDList = append(needDeleteUserIDList, userID)
				}
			}
			w.hub.Mu.RUnlock()
			w.hub.Mu.Lock()
			for _, info := range needDeleteDeviceOfUserList {
				delete(w.hub.Client[info.UserID], info.DeviceID)
			}
			for _, info := range needDeleteUserIDList {
				delete(w.hub.Client, info)
			}
			w.hub.Mu.Unlock()
		}
	}()
}
