package main

import (
	"fmt"
	"net/http"
	"rtforum/server/auth"
	"rtforum/server/database"
	"rtforum/server/dummyData"
	"rtforum/server/handlers" // Adjust this import based on your project structure
	"rtforum/server/websocket"
)

func main() {
	defer database.MyDB.Close()
	database.Init()
	dummyData.AddDummyData()
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./static/images"))))

	// Serve the SPA (index.html)
	mux.HandleFunc("/", handlers.HandleSPA)
	mux.HandleFunc("/posts", handlers.PostHandler)
	mux.HandleFunc("/addPost", handlers.CreatePost)
	mux.HandleFunc("/addComment", handlers.CreateComment)
	mux.HandleFunc("/users", handlers.UserListHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/signin", auth.LoginUser)
	mux.HandleFunc("/logout", auth.LogoutUser)
	mux.HandleFunc("/signup", auth.Register)
	mux.HandleFunc("/addMessage", handlers.AddMessage)
	mux.HandleFunc("/ws", websocket.WebSocketHandler)

	// Start the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server is running on: http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err)
	}
}
