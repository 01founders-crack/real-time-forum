package database

import (
	"fmt"
	"rtforum/server/models"
	"sort"
	"time"

	"github.com/gofrs/uuid"
)

type MessageTime struct {
	Id              uuid.UUID
	Nickname        string
	LastMessageSent time.Time
}

type By func(p1, p2 *MessageTime) bool

type MessageSorter struct {
	message []MessageTime
	by      func(p1, p2 *MessageTime) bool
}

func (by By) Sort(message []MessageTime) {
	ps := &MessageSorter{
		message: message,
		by:      by,
	}
	sort.Sort(ps)
}

func (s *MessageSorter) Len() int {
	return len(s.message)
}

func (s *MessageSorter) Swap(i, j int) {
	s.message[i], s.message[j] = s.message[j], s.message[i]
}

func (s *MessageSorter) Less(i, j int) bool {
	return s.by(&s.message[i], &s.message[j])
}

func SortMessages(currentuser string, userList []models.User) []MessageTime {
	var userMessages []MessageTime
	currentuserId, err := FindIdByNickname(currentuser)
	if err != nil {
		fmt.Println("error finding user:", currentuser)
	}
	for _, v := range userList {
		timeSent, err := FindTimeOfLastMessageBetweenTwoUsers(currentuserId, v.Id)
		if err != nil {
			fmt.Println("error finding messages:", currentuser)
		}
		userMessage := MessageTime{v.Id, v.Nickname, timeSent}
		userMessages = append(userMessages, userMessage)
	}
	By(func(p1, p2 *MessageTime) bool {
		return p1.LastMessageSent.After(p2.LastMessageSent)
	}).Sort(userMessages)
	return userMessages
}
