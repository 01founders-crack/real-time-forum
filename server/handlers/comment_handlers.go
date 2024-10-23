package handlers

import (
	"fmt"
	"net/http"
	"rtforum/server/database"
	"rtforum/server/models"
	"rtforum/server/sessions"

	"github.com/gofrs/uuid"
)

// Expected input format from frontend
type CommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

// CreateComment handles the creation of a comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	// Step 1: Check if the user is authenticated

	username, valid := sessions.ValidateSession(r)
	if !valid {
		// If the user is not authenticated, return 401 Unauthorized for testing purposes
		if r.Header.Get("Test-Mode") == "true" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Otherwise, redirect to the login page (for production usage)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, err := database.FindIdByNickname(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Step 2: Parse the request body
	var req CommentRequest
	err = r.ParseForm()
	
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	req.PostID = r.FormValue("postId")
	req.Content = r.FormValue("commentMessage")

	// Step 3: Validate the input fields
	if req.Content == "" || req.PostID == "" {
		http.Error(w, "Invalid input: missing fields", http.StatusBadRequest)
		return
	}

	// Step 4: Convert PostID and UserID to UUID
	postID, err := uuid.FromString(req.PostID)
	if err != nil {
		http.Error(w, "Invalid PostID", http.StatusBadRequest)
		return
	}

	// Debugging statements before calling AddComment
	fmt.Printf("Attempting to add comment: PostID=%v, UserID=%v, Content=%v\n", postID, userID, req.Content)

	// Insert the comment into the database
	err = database.AddComment(models.Comment{PostId: postID, UserId: userID, Content: req.Content})
	if err != nil {
		// Debugging statement for when the insertion fails
		fmt.Printf("Failed to add comment: %v\n", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Debugging statement for successful insertion
	fmt.Println("Comment added successfully")

	// Send a success response
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("Comment created successfully"))
}
