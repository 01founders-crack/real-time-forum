// websocket.js
// Create a new WebSocket connection
const socket = new WebSocket('ws://localhost:8080/ws'); // Update the URL if needed

// Event handler for when the connection is open
socket.onopen = () => {
    console.log('Connected to WebSocket server');
    // Optional: Send an initial message or handshake if needed
};

// Event handler for receiving messages from the server
socket.onmessage = (event) => {
    console.log('Message received:', event.data);
    // Update the UI when a new message is received
    updateUI(event.data);
};

// Event handler for when the connection is closed
socket.onclose = (event) => {
    console.log('Disconnected from WebSocket server', event);
    // Optionally, handle reconnection logic here
};

// Event handler for WebSocket errors
socket.onerror = (error) => {
    console.error('WebSocket error:', error);
};

// Function to send messages to the server
function sendMessage(message) {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(message);
        console.log('Message sent to server', message);
    } else {
        console.error('WebSocket is not open');
    }
}

//Sending Messages
// Function to handle message submission from a form
document.getElementById('messageForm').addEventListener('submit', (event) => {
    event.preventDefault();
    const message = document.getElementById('messageInput').value;
    console.log('Sending message:', message);
    sendMessage(message);
    document.getElementById('messageInput').value = ''; // Clear the input after sending
});

// Dynamically Updating the UI
// Function to update the UI with a new message
function updateUI(message) {
    const messageList = document.getElementById('messageList');
    const newMessage = document.createElement('li');
    newMessage.textContent = message;
    messageList.appendChild(newMessage); // Add the new message to the list
    console.log('Message added to UI:', message);
}

// <!-- Message List to Display Messages -->
// <ul id="messageList"></ul>

// <!-- Form to Send Messages -->
// <form id="messageForm">
//     <input type="text" id="messageInput" placeholder="Enter message" required />
//     <button type="submit">Send</button>
// </form>