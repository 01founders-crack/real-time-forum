package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"rtforum/server/database"
	"rtforum/server/sessions"
	"rtforum/server/utils"
	"strings"
	"testing"

	"github.com/gofrs/uuid" // Import the UUID package
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:") // In-memory test DB
	if err != nil {
		return nil, err
	}

	// Create the Users table
	query := `
    CREATE TABLE Users (
        ID TEXT PRIMARY KEY,
        Nickname TEXT,
        Age INTEGER,
        Gender TEXT,
        FirstName TEXT,
        LastName TEXT,
        Email TEXT,
        Password TEXT
    );
    CREATE TABLE Sessions (
        ID TEXT PRIMARY KEY,
        UserID TEXT,
        CreatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	// Set the global database instance to this test database
	database.MyDB = db

	return db, nil
}

// Insert a mock user for login tests (new function)
func setupMockUser() error {
	userID, err := uuid.NewV4() // Generate a UUID
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword("validPassword")
	if err != nil {
		return err
	}

	// Insert the mock user into the database with the generated UUID
	_, err = database.MyDB.Exec(`INSERT INTO Users (ID, Nickname, Age, Gender, FirstName, LastName, Email, Password)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID.String(), "testUser", 25, "male", "John", "Doe", "test@example.com", hashedPassword)
	return err
}

func TestRegisterUser(t *testing.T) {
	// Set up the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	// Setup test data
	nickname := "johnDoe"
	age := 25
	gender := "male"
	firstName := "John"
	lastName := "Doe"
	email := "john.doe@example.com"
	password := "validPassword123"

	// Run the RegisterUser function
	err = RegisterUser(nickname, age, gender, firstName, lastName, email, password)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRegisterUser_ValidRegistration(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	// Test with valid registration data
	err = RegisterUser("johnDoe", 25, "male", "John", "Doe", "john.doe@example.com", "validPassword123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRegisterUser_InvalidEmail(t *testing.T) {
	// Set up the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	// Test with an invalid email
	err = RegisterUser("johnDoe", 25, "male", "John", "Doe", "invalid-email", "validPassword123")
	if err == nil || err.Error() != "invalid email format" {
		t.Errorf("Expected 'invalid email format' error, got %v", err)
	}
}

func TestRegisterUser_WeakPassword(t *testing.T) {
	// Set up the database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	err = RegisterUser("johnDoe", 25, "male", "John", "Doe", "john.doe@example.com", "weak")
	if err == nil || err.Error() != "password does not meet strength requirements" {
		t.Errorf("Expected 'password does not meet strength requirements' error, got %v", err)
	}
}

func TestRegisterUser_EmptyNickName(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	err = RegisterUser("", 25, "male", "John", "Doe", "john.doe@example.com", "validPassword123")
	if err == nil || err.Error() != "nickname cannot be empty" {
		t.Errorf("Expected 'nickname cannot be empty' error, got %v", err)
	}
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	// Register First User
	err = RegisterUser("johnDoe", 25, "male", "John", "Doe", "john.doe@example.com", "validPassword123")
	if err != nil {
		t.Fatalf("Expected no error for the first user, got %v", err)
	}

	// Register Second User
	secondErr := RegisterUser("jamesDoe", 25, "male", "James", "Doe", "john.doe@example.com", "validPassword123")
	if secondErr == nil || secondErr.Error() != "email already exists" {
		t.Errorf("Expected 'email already exists', got %v", secondErr)
	}
}

func TestRegisterUser_ExceedingMaximumNicknameLength(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	longNickname := "aNicknameThatIsWayTooLongForTheSystemToHandleProperly"
	err = RegisterUser(longNickname, 25, "male", "John", "Doe", "john.doe@example.com", "validPassword123")
	if err == nil || err.Error() != "nickname exceeds maximum length" {
		t.Errorf("Expected 'nickname exceeds maximum length' error, got %v", err)
	}
}

func TestRegisterUser_PasswordMissingLetters(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up a test database: %v", err)
	}
	defer db.Close()

	err = RegisterUser("johnDoe", 25, "male", "John", "Doe", "john.doe@example.com", "12345678")
	if err == nil || err.Error() != "password does not meet strength requirements" {
		t.Errorf("Expected 'password does not meet strength requirements' error, got %v", err)
	}
}

// Test cases past this point have not passed

// TestLogoutUser tests the user logout functionality
// TestLogoutUser tests the user logout functionality
func TestLogoutUser(t *testing.T) {
	// Set up the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	// Insert mock user into the test database
	setupMockUser := func() error {
		// Hash the password
		hashedPassword, err := utils.HashPassword("validPassword")
		if err != nil {
			return err
		}

		random, _ := uuid.NewV4()
		// Insert the user into the database
		_, err = database.MyDB.Exec(`INSERT INTO Users (ID, Nickname, Age, Gender, FirstName, LastName, Email, Password)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			random, "testUser", 25, "male", "John", "Doe", "test@example.com", hashedPassword)
		return err
	}

	// Call setupMockUser before creating a session
	if err := setupMockUser(); err != nil {
		t.Fatalf("Failed to set up mock user: %v", err)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Mock the session for the test
	sessions.CreateSession(rr, req, "testUser")

	// Call the LogoutUser handler
	handler := http.HandlerFunc(LogoutUser)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	// Check if the session was destroyed
	_, valid := sessions.ValidateSession(req)
	if valid {
		t.Errorf("expected session to be destroyed, but it was still valid")
	}

	// Check the redirect location
	if rr.Header().Get("Location") != "/login" {
		t.Errorf("expected to be redirected to /login, but got %v", rr.Header().Get("Location"))
	}
}

// Test for successful login
func TestLoginUser_Success(t *testing.T) {
	// Set up the test database
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	// Insert mock user
	if err := setupMockUser(); err != nil {
		t.Fatalf("Failed to set up mock user: %v", err)
	}

	// Prepare valid login request body
	body := strings.NewReader(`{"identifier": "testUser", "password": "validPassword"}`)
	req, err := http.NewRequest("POST", "/login", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()

	// Call the LoginUser function
	handler := http.HandlerFunc(LoginUser)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	if rr.Body.String() != "Login successful" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Login successful")
	}
}

// Test for invalid password
func TestLoginUser_InvalidPassword(t *testing.T) {
	// Setup mock database with user
	setupTestDB()
	setupMockUser()

	// Prepare invalid login request body (wrong password)
	body := strings.NewReader(`{"identifier": "testUser", "password": "wrongPassword"}`)
	req, err := http.NewRequest("POST", "/login", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()

	// Call the LoginUser function
	handler := http.HandlerFunc(LoginUser)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Check the response body
	if strings.TrimSpace(rr.Body.String()) != "Invalid credentials" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Invalid credentials")
	}
}

// Test for missing credentials
func TestLoginUser_MissingCredentials(t *testing.T) {
	setupTestDB()
	setupMockUser()
	// Prepare invalid login request body (missing password)
	body := strings.NewReader(`{"identifier": "testUser"}`)
	req, err := http.NewRequest("POST", "/login", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()

	// Call the LoginUser function
	handler := http.HandlerFunc(LoginUser)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body
	if strings.TrimSpace(rr.Body.String()) != "missing credentials" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Missing credentials")
	}
}

// Test for non-existent user
func TestLoginUser_NonExistentUser(t *testing.T) {
	setupTestDB()
	setupMockUser()
	// Prepare invalid login request body (non-existent user)
	body := strings.NewReader(`{"identifier": "nonexistent", "password": "somePassword"}`)
	req, err := http.NewRequest("POST", "/login", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()

	// Call the LoginUser function
	handler := http.HandlerFunc(LoginUser)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	// Check the response body
	if strings.TrimSpace(rr.Body.String()) != "Invalid credentials" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "Invalid credentials")
	}
}
