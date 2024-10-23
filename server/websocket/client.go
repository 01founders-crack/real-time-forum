package websocket

// import (
// 	"github.com/gorilla/websocket"
// )

// // Client struct representing a single Websocket connection
// type Client struct {
// 	conn     *websocket.Conn
// 	username string
// 	manager  *ConnectionManager
// 	send     chan []byte
// }

// // NewClient create a new WebSocket client
// func NewClient(conn *websocket.Conn, username string, manager *ConnectionManager) *Client {
// 	return &Client{
// 		conn:     conn,
// 		username: username,
// 		manager:  manager,
// 		send:     make(chan []byte),
// 	}
// }

// // ReadMessage listens for incoming WebSocket messages
// func (c *Client) ReadMessages() {
// 	defer func() {
// 		c.manager.unregister <- c
// 		c.conn.Close()
// 	}()

// 	for {
// 		_, message, err := c.conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		//Assuming this is a private message format: {to: "username", content: "message"}
// 		parsedMessage := parseMessage(message)
// 		c.manager.SendMessage(parsedMessage)
// 	}
// }

// // WriteMessage sends messages back to the client
// func (c *Client) WriteMessage() {
// 	defer c.conn.Close()

// 	for message := range c.send {
// 		err := c.conn.WriteMessage(websocket.TextMessage, message)
// 		if err != nil {
// 			break
// 		}
// 	}
// }
