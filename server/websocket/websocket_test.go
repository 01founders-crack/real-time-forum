package websocket

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"rtforum/server/database"
	"rtforum/server/sessions"
	"rtforum/server/utils"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

// Helper function to set up the test database
func setupTestDB() (*sql.DB, error) {
	// Create an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Create the Users table with necessary columns, including Age
	usersTableQuery := `
    CREATE TABLE Users (
        ID TEXT PRIMARY KEY,
        Nickname TEXT,
        Age INTEGER,
        Gender TEXT,
        FirstName TEXT,
        LastName TEXT,
        Email TEXT,
        Password TEXT
    );`
	_, err = db.Exec(usersTableQuery)
	if err != nil {
		return nil, err
	}

	// Create the Sessions table
	sessionsTableQuery := `
    CREATE TABLE Sessions (
        ID TEXT PRIMARY KEY,
        UserID TEXT,
        CreatedAt TIMESTAMP
    );`
	_, err = db.Exec(sessionsTableQuery)
	if err != nil {
		return nil, err
	}

	// Set the global database instance to this test database
	database.MyDB = db

	return db, nil
}

// Helper function to insert a mock user into the database
func setupMockUser() error {
	// Generate a valid UUID for the mock user
	userID := uuid.Must(uuid.NewV4()).String()

	// Hash the password using bcrypt
	hashedPassword, err := utils.HashPassword("validPassword")
	if err != nil {
		return err
	}

	_, err = database.MyDB.Exec(`INSERT INTO Users (ID, Nickname, Age, Gender, FirstName, LastName, Email, Password)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, "testUser", 25, "male", "John", "Doe", "test@example.com", hashedPassword)
	return err
}

// TestWebSocketHandler ensures WebSocket connections work for authenticated users
func TestWebSocketHandler_Success(t *testing.T) {
	// Step 1: Set up the test database and mock user
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	err = setupMockUser()
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	// Step 2: Set up the test HTTP server
	server := httptest.NewServer(http.HandlerFunc(WebSocketHandler))
	defer server.Close()

	// Step 3: Create a request and response recorder
	req, _ := http.NewRequest("GET", "/ws", nil)
	rr := httptest.NewRecorder()

	// Step 4: Create a session and set the username
	err = sessions.CreateSession(rr, req, "testUser") // Simulate a session for testUser
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Step 5: Copy cookies from recorder to the request (to simulate session)
	for _, cookie := range rr.Result().Cookies() {
		req.AddCookie(cookie)
	}

	// Step 6: Connect to the WebSocket server using the same cookies
	wsUrl := "ws" + server.URL[len("http"):] + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, req.Header)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Step 7: Send a message
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Step 8: Receive the echoed message
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	expectedMessage := "Hello, WebSocket!"
	if string(message) != expectedMessage {
		t.Errorf("Expected %s, got %s", expectedMessage, message)
	}
}

// Test unauthorized WebSocket connection (no session)
func TestWebSocketHandler_Unauthorized(t *testing.T) {
	// Create a test HTTP server with the WebSocketHandler
	server := httptest.NewServer(http.HandlerFunc(WebSocketHandler))
	defer server.Close()

	// Simulate WebSocket client connection (without session)
	_, resp, err := websocket.DefaultDialer.Dial("ws"+server.URL[len("http"):]+" /ws", nil)
	if err == nil {
		t.Fatalf("Expected WebSocket connection to fail, but it succeeded")
	}
	if resp != nil && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status code 401 Unauthorized, but got %v", resp.StatusCode)
	}
}

// TestWebSocketBroadcast ensures authenticated users can broadcast messages via WebSocket
func TestWebSocketBroadcast(t *testing.T) {
	// Step 1: Set up the test database and mock user
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	err = setupMockUser()
	if err != nil {
		t.Fatalf("Failed to insert mock user: %v", err)
	}

	// Step 2: Set up the test HTTP server
	server := httptest.NewServer(http.HandlerFunc(WebSocketHandler))
	defer server.Close()

	// Step 3: Simulate a valid session
	req, _ := http.NewRequest("GET", "/ws", nil)
	rr := httptest.NewRecorder()

	// Create session and ensure it's saved
	sessions.CreateSession(rr, req, "testUser")

	// Extract the session cookie from the response
	sessionCookie := rr.Result().Cookies()[0]

	// Step 4: Connect to the WebSocket server with the session cookie
	wsUrl := "ws" + server.URL[len("http"):] + "/ws"
	header := http.Header{}
	header.Add("Cookie", sessionCookie.String()) // Attach session cookie

	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, header)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Step 5: Send a message
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Step 6: Receive the echoed message
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	expectedMessage := "Hello, WebSocket!"
	if string(message) != expectedMessage {
		t.Errorf("Expected %s, got %s", expectedMessage, message)
	}
}
