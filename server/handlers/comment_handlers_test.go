package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"rtforum/server/database"
	"rtforum/server/sessions"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
)

// SetupTestDB initializes the test database and creates necessary tables
func SetupTestDB() {
	var err error
	database.MyDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Create the Comments table
	createCommentsTable := `
	CREATE TABLE Comments (
		id TEXT PRIMARY KEY,
		post_id TEXT,
		user_id TEXT,
		content TEXT
	);
	`

	_, err = database.MyDB.Exec(createCommentsTable)
	if err != nil {
		panic(err)
	}
}

// Helper function to create a session for testing
func createSession(req *http.Request, username string) {
	rr := httptest.NewRecorder()
	sessions.CreateSession(rr, req, username)

	// Add the session cookie to the request
	cookie := rr.Result().Cookies()[0]
	req.AddCookie(cookie)
}

// TestCreateComment_Success tests successful comment creation
func TestCreateComment_Success(t *testing.T) {
	// Initialize the test database
	SetupTestDB()

	// Create a valid request body
	postID := uuid.Must(uuid.NewV4()).String()
	userID := uuid.Must(uuid.NewV4()).String()
	body := strings.NewReader(`{"post_id": "` + postID + `", "user_id": "` + userID + `", "content": "This is a test comment"}`)

	req, err := http.NewRequest("POST", "/create-comment", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a valid session
	createSession(req, "testUser")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateComment)

	// Call the CreateComment handler
	handler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if rr.Body.String() != "Comment created successfully" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Comment created successfully")
	}
}

// TestCreateComment_MissingFields tests the case when some fields are missing
func TestCreateComment_MissingFields(t *testing.T) {
	// Create a request with missing content field
	body := strings.NewReader(`{"post_id": "123", "user_id": "456"}`)

	req, err := http.NewRequest("POST", "/create-comment", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a valid session
	createSession(req, "testUser")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateComment)

	// Call the CreateComment handler
	handler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if strings.TrimSpace(rr.Body.String()) != "Invalid input: missing fields" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Invalid input: missing fields")
	}
}

// TestCreateComment_InvalidUUID tests the case when an invalid UUID is provided
func TestCreateComment_InvalidUUID(t *testing.T) {
	// Create a request with an invalid PostID
	body := strings.NewReader(`{"post_id": "invalid-uuid", "user_id": "invalid-uuid", "content": "This is a test comment"}`)

	req, err := http.NewRequest("POST", "/create-comment", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a valid session
	createSession(req, "testUser")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateComment)

	// Call the CreateComment handler
	handler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	if !strings.Contains(rr.Body.String(), "Invalid PostID") {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Invalid PostID")
	}
}

// TestCreateComment_Unauthorized tests the case when the user is not authenticated
func TestCreateComment_Unauthorized(t *testing.T) {
	// Create a valid request body
	body := strings.NewReader(`{"post_id": "123", "user_id": "456", "content": "This is a test comment"}`)

	req, err := http.NewRequest("POST", "/create-comment", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Test-Mode", "true") // Add this line for testing

	// Do not create a session, so the user is unauthorized

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateComment)

	// Call the CreateComment handler
	handler.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	if strings.TrimSpace(rr.Body.String()) != "Unauthorized" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Unauthorized")
	}
}

// TestCreateComment_Unauthenticated checks that unauthenticated users are redirected
func TestCreateComment_Unauthenticated(t *testing.T) {
	req, err := http.NewRequest("POST", "/create-comment", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(CreateComment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("expected status code %v, got %v", http.StatusSeeOther, status)
	}

	if rr.Header().Get("Location") != "/login" {
		t.Errorf("expected to be redirected to /login, got %v", rr.Header().Get("Location"))
	}
}

// TestCreateComment_Authenticated checks that authenticated users can create comments
func TestCreateComment_Authenticated(t *testing.T) {
	// Create a valid request body
	postID := uuid.Must(uuid.NewV4()).String()
	userID := uuid.Must(uuid.NewV4()).String()
	body := strings.NewReader(`{"post_id": "` + postID + `", "user_id": "` + userID + `", "content": "This is a test comment"}`)

	req, err := http.NewRequest("POST", "/create-comment", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Simulate a valid session
	sessions.CreateSession(rr, req, "testUser")

	handler := http.HandlerFunc(CreateComment)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("expected status code %v, got %v", http.StatusCreated, status)
	}

	if rr.Body.String() != "Comment created successfully" {
		t.Errorf("expected 'Comment created successfully', got %v", rr.Body.String())
	}
}
