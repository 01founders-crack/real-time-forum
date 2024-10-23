// // connection_manager.go
package websocket

// import (
// 	"log"
// 	"sync"
// )

// // ConnectionManager struct
// type ConnectionManager struct {
// 	clients    map[string]*Client
// 	register   chan *Client
// 	unregister chan *Client
// 	broadcast  chan Message
// 	mu         sync.Mutex
// }

// // NewConnectionManager creates a new instance of ConnectionManager
// func NewConnectionManager() *ConnectionManager {
// 	return &ConnectionManager{
// 		clients:    make(map[string]*Client),
// 		register:   make(chan *Client),
// 		unregister: make(chan *Client),
// 		broadcast:  make(chan Message),
// 	}
// }

// // Start the connection manager, handling register/unregister
// func (cm *ConnectionManager) Start() {
// 	for {
// 		select {
// 		case client := <-cm.register:
// 			cm.mu.Lock()
// 			cm.clients[client.username] = client
// 			cm.mu.Unlock()
// 			log.Printf("User %s connected", client.username)

// 		case client := <-cm.unregister:
// 			cm.mu.Lock()
// 			delete(cm.clients, client.username)
// 			cm.mu.Unlock()
// 			log.Printf("User %s disconnected", client.username)

// 		case message := <-cm.broadcast:
// 			cm.mu.Lock()
// 			if client, ok := cm.clients[message.To]; ok {
// 				client.send <- message.Content
// 			}
// 			cm.mu.Unlock()
// 		}
// 	}
// }

// // SendMessage to a specific user
// func (cm *ConnectionManager) SendMessage(message Message) {
// 	cm.broadcast <- message
// }
