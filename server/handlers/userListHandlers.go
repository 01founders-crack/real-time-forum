package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rtforum/server/database"
	"rtforum/server/models"
	"rtforum/server/sessions"
	"strconv"
)

func UserListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
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
	userlist, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users := database.SortMessages(username, userlist)
	var userList []struct {
		UserId             string           `json:"userId"`
		Nickname           string           `json:"nickname"`
		Status             string           `json:"status"`
		CurrentUser        bool             `json:"currentUser"`
		Messages           []models.Message `json:"messages"`
		TotalNotifications int              `json:"totalNotifications"`
		Notifications      map[string]int   `json:"notifications"`
	}

	for _, v := range users {
		var user struct {
			UserId             string           `json:"userId"`
			Nickname           string           `json:"nickname"`
			Status             string           `json:"status"`
			CurrentUser        bool             `json:"currentUser"`
			Messages           []models.Message `json:"messages"`
			TotalNotifications int              `json:"totalNotifications"`
			Notifications      map[string]int   `json:"notifications"`
		}
		for i, v := range v.Id {
			if i == 4 || i == 6 || i == 8 || i == 10 {
				user.UserId += "-"
			}
			if len(strconv.FormatInt(int64(v), 16)) != 2 {
				user.UserId += "0"
			}
			user.UserId += strconv.FormatInt(int64(v), 16)
		}
		user.Nickname = v.Nickname
		if username == v.Nickname {
			user.CurrentUser = true
		} else {
			user.CurrentUser = false
		}
		active, err := database.FindSessionByUserId(v.Id)
		if err != nil {
			fmt.Println("fknjwdeojbnjkbkj", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if active {
			user.Status = "Online"
		} else {
			user.Status = "Offline"
		}
		user.TotalNotifications, err = database.TotalNumberOfMessages(v.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Notifications, err = database.NotificationList(v.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !user.CurrentUser {
			currentUserId, err := database.FindIdByNickname(username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			user.Messages, err = database.FindAllMessagesBetweenTwoUsers(currentUserId, v.Id)
			for i := range user.Messages {
				user.Messages[i].ReceiverName, err = database.FindNicknameById(user.Messages[i].ReceiverId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				user.Messages[i].SenderName, err = database.FindNicknameById(user.Messages[i].SenderId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				user.Messages[i].TimeString = user.Messages[i].CreatedAt.Format("2006-01-02 15:04")
			}
		}
		userList = append(userList, user)
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Marshal the recipes slice to JSON
	err = json.NewEncoder(w).Encode(userList)
	if err != nil {
		http.Error(w, err.Error()+"WHY2", http.StatusBadRequest)
		return
	}
}
