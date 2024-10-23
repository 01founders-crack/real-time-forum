package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"rtforum/server/database"
	"rtforum/server/models"
	"rtforum/server/sessions"
	"rtforum/server/utils"
	"strconv"
)

const MaxNicknameLength = 30 // Define the maximum length for the nickname

// RegisterUser validates inputs, hashes the password, and adds a new user to the database
func RegisterUser(nickname string, age int, gender, firstName, lastName, email, password string) error {

	// Check for empty fields
	if nickname == "" {
		return errors.New("nickname cannot be empty")
	}
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if password == "" {
		return errors.New("password cannot be empty")
	}

	// Check for maximum field lengths
	if len(nickname) > MaxNicknameLength {
		return errors.New("nickname exceeds maximum length")
	}

	// Validate email format
	if !utils.ValidateEmail(email) {
		return errors.New("invalid email format")
	}

	// Validate password strength
	if !utils.ValidatePassword(password) {
		return errors.New("password does not meet strength requirements")
	}

	// Check if the nickname already exists
	exists, err := database.NicknameExists(nickname)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("nickname already exists")
	}

	// Check if the email already exists
	exists, err = database.EmailExists(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// Add the user to the database
	err = database.AddUser(models.User{Nickname: nickname, Age: age, Gender: gender, FirstName: firstName, LastName: lastName, Email: email, Password: hashedPassword})
	if err != nil {
		return err
	}

	return nil
}

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		fmt.Println("Error converting to int:", err) // Debug statement
		http.Error(w, "Failed to convert to int", http.StatusInternalServerError)
		return
	}
	err = RegisterUser(r.FormValue("nickname"), age, r.FormValue("gender"), r.FormValue("firstName"), r.FormValue("lastName"), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		fmt.Println("Error registering user:", err) // Debug statement
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// message := fmt.Sprintf("registration successful")
	// // Return success message
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(message))
}

// LoginUser is the function that handlers user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Yay")
	// Define the request struct for parsing the input
	var req struct {
		Identifier string `json:"identifier"` // Can be email or nickname
		Password   string `json:"password"`
	}

	// Decode the incoming request
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	http.Error(w, "Invalid input", http.StatusBadRequest)
	// 	return
	// }
	r.ParseForm()
	req.Identifier = r.FormValue("input-login-email")
	req.Password = r.FormValue("input-login-password")
	// Log login attempt
	log.Printf("Received login attempt: Identifier=%s", req.Identifier)

	// Validate login
	isValid, err := database.IsLoginValid(req.Identifier, req.Password)
	if err != nil {
		if error.Error(err) == "missing credentials" {
			log.Printf("Error missing credentials: %v", err) // Add more logging
			http.Error(w, "missing credentials", http.StatusBadRequest)
			return
		}
		log.Printf("Error validating login: %v", err) // Add more logging
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isValid {
		log.Printf("Invalid credentials for: %s", req.Identifier)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session after successful login
	err = sessions.CreateSession(w, r, req.Identifier)
	if err != nil {
		log.Printf("Failed to create session: %v", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	// message := fmt.Sprintf("Login successful, Identifier: %v", req.Identifier)
	// // Return success message
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(message))
}

// LogoutUser handles user logout by destroying the session
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Debugging message to track logout
	fmt.Println("Starting LogoutUser...")

	// Use the existing DestroySession function from sessions.go
	err := sessions.DestroySession(w, r)
	if err != nil {
		fmt.Println("Error destroying session:", err) // Debug statement
		http.Error(w, "Failed to destroy session", http.StatusInternalServerError)
		return
	}

	fmt.Println("Session destroyed successfully") // Debug statement

	// Redirect the user to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// message := fmt.Sprintf("Logout successful")
	// // Return success message
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(message))
}
