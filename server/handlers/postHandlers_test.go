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
	_ "github.com/mattn/go-sqlite3" // SQLite driver for testing
)

// SetupTestDB creates the necessary tables for testing
func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:") // In-memory database for testing
	if err != nil {
		return nil, err
	}

	// Create the Posts table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS Posts (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		title TEXT,
		content TEXT,
		category TEXT
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// TestCreatePostFunctionality tests the overall creation of a post and database interaction
func TestCreatePostFunctionality(t *testing.T) {
	// Set up the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	// Replace the actual database connection with our test database
	database.MyDB = db

	// Create a valid request body
	userID := uuid.Must(uuid.NewV4()).String()
	body := strings.NewReader(`{"userId": "` + userID + `", "title": "Test Post", "category": "General", "content": "This is a test post"}`)

	req, err := http.NewRequest("POST", "/create-post", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Simulate a valid session
	rr := httptest.NewRecorder()
	sessions.CreateSession(rr, req, "testUser")

	// Call the CreatePost handler
	handler := http.HandlerFunc(CreatePost)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Verify that the post was inserted into the database
	var insertedTitle, insertedContent, insertedCategory string
	err = db.QueryRow(`SELECT title, content, category FROM Posts WHERE user_id = ?`, userID).Scan(&insertedTitle, &insertedContent, &insertedCategory)
	if err != nil {
		t.Fatalf("Failed to retrieve post from database: %v", err)
	}

	// Check that the inserted data matches the original post data
	if insertedTitle != "Test Post" {
		t.Errorf("expected title to be 'Test Post', got '%v'", insertedTitle)
	}
	if insertedContent != "This is a test post" {
		t.Errorf("expected content to be 'This is a test post', got '%v'", insertedContent)
	}
	if insertedCategory != "General" {
		t.Errorf("expected category to be 'General', got '%v'", insertedCategory)
	}
}

// TestCreatePost_Unauthorized tests that an unauthenticated user cannot create a post
func TestCreatePost_Unauthorized(t *testing.T) {
	// Create a valid request body
	body := strings.NewReader(`{"userId": "123", "title": "Test Post", "category": "General", "content": "This is a test post"}`)

	req, err := http.NewRequest("POST", "/create-post", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Do not create a session, so the user is unauthorized
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePost)

	// Call the CreatePost handler
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Check the response body
	if strings.TrimSpace(rr.Body.String()) != "Unauthorized" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Unauthorized")
	}
}

func TestCreatePost_Authenticated(t *testing.T) {
	// Set up the test database and ensure the Posts table exists
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	// Replace the actual database connection with our test database
	database.MyDB = db

	// Create a valid request body
	userID := uuid.Must(uuid.NewV4()).String()
	body := strings.NewReader(`{"userId": "` + userID + `", "title": "Test Post", "category": "General", "content": "This is a test post"}`)

	req, err := http.NewRequest("POST", "/create-post", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Simulate a valid session
	rr := httptest.NewRecorder()
	sessions.CreateSession(rr, req, "testUser")

	// Call the CreatePost handler
	handler := http.HandlerFunc(CreatePost)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	if rr.Body.String() != "Post created successfully" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Post created successfully")
	}
}

func TestCreatePost_InvalidInput(t *testing.T) {
	// Set up the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	// Replace the actual database connection with our test database
	database.MyDB = db

	// Create a request with missing title (invalid input)
	body := strings.NewReader(`{"userId": "123", "title": "", "category": "General", "content": "This is a test post"}`)
	req, err := http.NewRequest("POST", "/create-post", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Simulate a valid session (authentication should succeed)
	rr := httptest.NewRecorder()
	sessions.CreateSession(rr, req, "testUser")

	// Call the CreatePost handler
	handler := http.HandlerFunc(CreatePost)
	handler.ServeHTTP(rr, req)

	// Check for Bad Request due to missing title
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
