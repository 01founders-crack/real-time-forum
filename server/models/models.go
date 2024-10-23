package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Nickname  string    `json:"nickname"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type Post struct {
	Id       uuid.UUID `json:"id"`
	StringId string    `json:"stringId"`
	UserId   uuid.UUID `json:"userId"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Content  string    `json:"content"`
	Comments []Comment `json:"comments"`
	Username string    `json:"username"` //only for front
}

type Comment struct {
	Id       uuid.UUID `json:"id"`
	PostId   uuid.UUID `json:"postId"`
	UserId   uuid.UUID `json:"userId"`
	Content  string    `json:"content"`
	Username string    `json:"username"` //only for front
}

type Message struct {
	Id           uuid.UUID `json:"id"`
	StringId     string    `json:"stringId"`
	SenderId     uuid.UUID `json:"senderId"`
	ReceiverId   uuid.UUID `json:"receiverId"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
	SenderName   string    `json:"senderName"`   //only for front
	ReceiverName string    `json:"receiverName"` //only for front
	TimeString   string    `json:"time"`
}

type Session struct {
	Id     uuid.UUID `json:"id"`
	UserId uuid.UUID `json:"userId"`
	Status string    `json:"status"`
}

type Category struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Desc string    `json:"desc"`
}
