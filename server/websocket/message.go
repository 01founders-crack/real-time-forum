package websocket

type Message struct {
	From    string
	To      string
	Content []byte
}

// Parse a raw message into a structured Message object
func parseMessage(rawMessage []byte) Message {
	return Message{
		From:    "testUser",
		To:      "targetUser",
		Content: rawMessage,
	}
}
