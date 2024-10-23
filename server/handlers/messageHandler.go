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

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		var data struct {
			TargetUser string `json:"nickname"`
		}
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(data)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data) // Echo back the received data
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 	return
	// }
	// var data struct {
	// 	CurrentUser uuid.UUID `json:"currentUser"`
	// 	TargetUser  uuid.UUID `json:"targetUser"`
	// }
	// err := json.NewDecoder(r.Body).Decode(&data)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// messages, err := database.FindAllMessagesBetweenTwoUsers(data.CurrentUser, data.TargetUser)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// for _, message := range messages {
	// 	message.ReceiverName, err = database.FindNicknameById(message.ReceiverId)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	message.SenderName, err = database.FindNicknameById(message.SenderId)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// }
	// // Set the content type to application/json
	// w.Header().Set("Content-Type", "application/json")

	// // Marshal the recipes slice to JSON
	// err = json.NewEncoder(w).Encode(messages)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}

func GetAllMessages(w http.ResponseWriter, r *http.Request) {

	messages, err := database.FindAllMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Marshal the recipes slice to JSON
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
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
	var req models.Message
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	req.SenderId = userId
	req.Content = r.FormValue("messageInput")
	req.ReceiverId, err = database.FindIdByNickname(r.FormValue("messageTargetUser"))
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, "Invalid input: missing fields", http.StatusBadRequest)
		return
	}

	fmt.Printf("Attempting to add Message: RecieverId=%v, SenderId=%v, Content=%v\n", req.ReceiverId, req.SenderId, req.Content)

	err = database.AddMessages(req)
	if err != nil {
		// Debugging statement in case of failure
		fmt.Printf("Failed to add message: %v\n", err)
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}
	fmt.Println("Post added successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
