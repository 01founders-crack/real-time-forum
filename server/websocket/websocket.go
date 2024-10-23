package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"rtforum/server/database"
	"rtforum/server/models"
	"rtforum/server/sessions"
	"sync"

	"github.com/gorilla/websocket"
)

// UserConnection associates a connection with a username
type UserConnection struct {
	conn     *websocket.Conn
	username string
}

// ConnectionManager to manage all active WebSocket connections
type ConnectionManager struct {
	connections map[string]*UserConnection // map username to connection
	broadcast   chan []byte                // keep broadcast channel for compatibility
	lock        sync.RWMutex
}

// NewConnectionManager creates a new instance of ConnectionManager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*UserConnection),
		broadcast:   make(chan []byte),
	}
}

// Start listens for broadcast messages (kept for backward compatibility)
func (manager *ConnectionManager) Start() {
	for {
		message := <-manager.broadcast
		manager.broadcastMessage(message)
	}
}

// broadcastMessage sends a message to all connected clients
func (manager *ConnectionManager) broadcastMessage(message []byte) {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	for _, userConn := range manager.connections {
		err := userConn.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Write error for user %s: %v", userConn.username, err)
			userConn.conn.Close()
			// Don't delete here to avoid concurrent map write
		}
	}
}

// AddClient adds a new WebSocket connection to the manager
func (manager *ConnectionManager) AddClient(conn *websocket.Conn, username string) {
	manager.lock.Lock()
	manager.connections[username] = &UserConnection{
		conn:     conn,
		username: username,
	}
	manager.lock.Unlock()
	log.Printf("Client added for user %s. Total clients: %d", username, len(manager.connections))
}

// RemoveClient removes a WebSocket connection from the manager
func (manager *ConnectionManager) RemoveClient(username string) {
	manager.lock.Lock()
	delete(manager.connections, username)
	manager.lock.Unlock()
	log.Printf("Client removed for user %s. Total clients: %d", username, len(manager.connections))
}

// SendToUser sends a message to a specific user
func (manager *ConnectionManager) SendToUser(target string, message []byte) error {
	manager.lock.RLock()
	conn, exists := manager.connections[target]
	manager.lock.RUnlock()

	if !exists {
		return nil // User not connected, message will be retrieved when they connect
	}

	return conn.conn.WriteMessage(websocket.TextMessage, message)
}

var connectManager = NewConnectionManager()

func init() {
	go connectManager.Start()
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Validate session
	username, valid := sessions.ValidateSession(r)
	if !valid {
		log.Println("Session validation failed - Session is invalid")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrade to WebSocket
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// Add client to the connection manager
	connectManager.AddClient(conn, username)
	defer connectManager.RemoveClient(username)

	log.Printf("User %s connected via WebSocket", username)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		var newMessage struct {
			Type     string `json:"type"`
			Username string `json:"username"`
			Message  string `json:"message"`
			Target   string `json:"target"`
		}

		if err := json.Unmarshal(message, &newMessage); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle different message types
		switch newMessage.Type {
		case "message":
			// Store message in database
			senderId, err := database.FindIdByNickname(newMessage.Username)
			if err != nil {
				log.Printf("Error finding sender ID: %v", err)
				continue
			}

			targetId, err := database.FindIdByNickname(newMessage.Target)
			if err != nil {
				log.Printf("Error finding target ID: %v", err)
				continue
			}

			// Add message to database
			err = database.AddMessages(models.Message{
				SenderId:   senderId,
				ReceiverId: targetId,
				Content:    newMessage.Message,
			})
			if err != nil {
				log.Printf("Error storing message: %v", err)
				continue
			}

			err = database.AddNotification(senderId, targetId)
			if err != nil {
				log.Printf("Error adding notification: %v", err)
				continue
			}

			// Send message only to the target user
			if err := connectManager.SendToUser(newMessage.Target, message); err != nil {
				log.Printf("Error sending to target user: %v", err)
			}

			// Send a copy back to the sender for their chat window
			if err := connectManager.SendToUser(newMessage.Username, message); err != nil {
				log.Printf("Error sending to sender: %v", err)
			}

		case "status":
			// Broadcast status updates to all users
			connectManager.broadcast <- message
		}
	}
}
