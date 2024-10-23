package sessions

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test creating a session
func TestCreateSession(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the CreateSession function
	err = CreateSession(rr, req, "testUser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the Set-Cookie header is set in the response
	cookie := rr.Result().Cookies()
	if len(cookie) == 0 {
		t.Errorf("Expected a session cookie, but none was set")
	}
}

// Test validating a session
func TestValidateSession(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a session
	err = CreateSession(rr, req, "testUser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Extract the session cookie and add it to a new request
	cookie := rr.Result().Cookies()[0]
	req.AddCookie(cookie)

	// Validate the session
	username, valid := ValidateSession(req)
	if !valid {
		t.Errorf("Expected session to be valid, but it was not")
	}
	if username != "testUser" {
		t.Errorf("Expected username to be 'testUser', got '%s'", username)
	}
}

// Test destroying a session
func TestDestroySession(t *testing.T) {
	// Step 1: Create a new request and response recorder
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Step 2: Create a session
	err = CreateSession(rr, req, "testUser")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Extract the session cookie from the response
	cookie := rr.Result().Cookies()[0]

	// Step 3: Create a new request to simulate a new client request
	newReq, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	newReq.AddCookie(cookie)

	// Step 4: Destroy the session with the new request
	rr = httptest.NewRecorder()
	err = DestroySession(rr, newReq)
	if err != nil {
		t.Fatalf("Expected no error when destroying session, got %v", err)
	}

	// Step 5: Create another new request to validate the session
	finalReq, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Try to validate the session using the new request (which should now be invalid)
	username, valid := ValidateSession(finalReq)
	if valid {
		t.Errorf("Expected session to be invalid, but it was valid")
	}
	if username != "" {
		t.Errorf("Expected no username, but got '%s'", username)
	}
}
