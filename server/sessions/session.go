package sessions

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rtforum/server/database"

	"github.com/gorilla/sessions"
)

var sessionKey = []byte("your-very-secret-key-that-is-at-least-32-bytes-long")

var store = sessions.NewCookieStore(sessionKey)

func CreateSession(w http.ResponseWriter, r *http.Request, username string) error {
	session, _ := store.Get(r, "session-name")
	session.Values["username"] = username // Set username in session

	log.Printf("Creating session for user: %s", username)

	// Save session
	err := session.Save(r, w)
	if err != nil {
		log.Printf("Failed to save session for user: %s, error: %v", username, err)
		return err
	}

	log.Printf("Creating database session for user: %s", username)
	userId, err := database.FindIdByLogin(username)
	// Save session
	err = database.AddSession(userId)
	if err != nil {
		log.Printf("Failed to add session to database for user: %s, error: %v", username, err)
		return err
	}

	return nil
}

func ValidateSession(r *http.Request) (string, bool) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Failed to retrieve session: %v", err)
		return "", false
	}
	username, ok := session.Values["username"].(string)
	if !ok {
		log.Println("Session does not contain username")
		return "", false
	}

	log.Printf("Session validated successfully for user: %s", username)
	return username, true
}

func DestroySession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "session-name")
	if err != nil {
		fmt.Println("Failed to retrieve session:", err)
		return err
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		return errors.New("username not found")
	}

	// Find user ID by nickname, and handle case where user doesn't exist
	id, err := database.FindIdByNickname(username)
	if err == sql.ErrNoRows {
		// Gracefully handle case where user is not found
		fmt.Println("Error finding user by nickname: user does not exist")
		return errors.New("user does not exist")
	} else if err != nil {
		return err
	}

	err = database.DeleteSession(id)
	if err != nil {
		return err
	}

	// Clear all session values
	for k := range session.Values {
		delete(session.Values, k)
	}

	// Save the session to invalidate it
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("Failed to save session:", err)
		return err
	}

	return nil
}
