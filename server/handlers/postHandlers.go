package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rtforum/server/database"
	"rtforum/server/models"
	"rtforum/server/sessions"
	"strings"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate the session to ensure the user is authenticated
	username, valid := sessions.ValidateSession(r)
	if !valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, err := database.FindIdByNickname(username)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Step 2: Parse the request body to get the post data
	var req struct {
		Title    string `json:"title"`
		Category string `json:"category"`
		Content  string `json:"content"`
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	req.Title = r.FormValue("title")
	req.Content = r.FormValue("content")
	req.Category = r.FormValue("category")

	// Step 3: Validate that the necessary fields are provided
	if strings.TrimSpace(req.Content) == "" || strings.TrimSpace(req.Title) == "" {
		http.Error(w, "Invalid input: missing fields", http.StatusBadRequest)
		return
	}

	// Step 4: Convert UserId to a UUID and validate it

	// Step 5: Add debugging statement before adding the post
	fmt.Printf("Attempting to add Post: UserID=%v, Title=%v, Content=%v, Category=%v\n", userId, req.Title, req.Content, req.Category)

	// Step 6: Insert the post into the database
	err = database.AddPost(models.Post{UserId: userId, Title: req.Title, Content: req.Content, Category: req.Category})
	if err != nil {
		// Debugging statement in case of failure
		fmt.Printf("Failed to add post: %v\n", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Debugging statement for successful insertion
	fmt.Println("Post added successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// Step 7: Send a success response
	// 	w.WriteHeader(http.StatusCreated)
	// 	w.Write([]byte("Post created successfully"))
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	posts, err := database.FindAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	// Marshal the recipes slice to JSON
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, err.Error()+"WHY2", http.StatusBadRequest)
		return
	}
}
