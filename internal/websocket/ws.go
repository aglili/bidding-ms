package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotificationMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type ConnectionManager struct {
	connections map[uuid.UUID]*websocket.Conn
	sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uuid.UUID]*websocket.Conn),
	}
}

func (cm *ConnectionManager) Register(userID uuid.UUID, conn *websocket.Conn) {
	cm.Lock()
	defer cm.Unlock()
	cm.connections[userID] = conn
	log.Printf("User %s connected via WebSocket", userID)
}

func (cm *ConnectionManager) UnRegister(userID uuid.UUID) {
	cm.Lock()
	defer cm.Unlock()
	if conn, exists := cm.connections[userID]; exists {
		conn.Close()
		delete(cm.connections, userID)
		log.Printf("User %s disconnected from WebSocket", userID)
	}
}

func (cm *ConnectionManager) SendToUser(userID uuid.UUID, message NotificationMessage) error {
	cm.RLock()
	conn, exists := cm.connections[userID]
	cm.RUnlock()

	if !exists {
		return nil
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)

}
