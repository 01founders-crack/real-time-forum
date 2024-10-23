package database

import (
	"database/sql"
	"rtforum/server/models"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

// takes in message details and sends the information to corresponding column in 'Messages' table
func AddMessages(message models.Message) error {
	statement, err := MyDB.Prepare("INSERT INTO Messages VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = statement.Exec(id, message.SenderId, message.ReceiverId, message.Content, time.Now())
	return err
}

func FindAllMessagesBetweenTwoUsers(currentUserId, targetUserId uuid.UUID) ([]models.Message, error) {
	var messages []models.Message
	rows, err := MyDB.Query("SELECT * FROM Messages WHERE (SenderId=? AND ReceiverId=?) OR (SenderId=? AND ReceiverId=?) ORDER BY CreatedAt", currentUserId, targetUserId, targetUserId, currentUserId)
	if err == sql.ErrNoRows {
		return messages, nil
	} else if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.Content, &message.CreatedAt)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func FindTimeOfLastMessageBetweenTwoUsers(currentUserId, targetUserId uuid.UUID) (time.Time, error) {
	var messageTime time.Time
	rows, err := MyDB.Query("SELECT CreatedAt FROM Messages WHERE (SenderId=? AND ReceiverId=?) OR (SenderId=? AND ReceiverId=?)", currentUserId, targetUserId, targetUserId, currentUserId)
	if err == sql.ErrNoRows {
		return messageTime, nil
	} else if err != nil {
		return messageTime, err
	}
	defer rows.Close()
	for rows.Next() {
		var currentMessageTime time.Time
		err = rows.Scan(&currentMessageTime)
		if err != nil {
			return messageTime, err
		}
		if currentMessageTime.After(messageTime) {
			messageTime = currentMessageTime
		}
	}
	return messageTime, nil
}

func FindAllMessagesBetweenTwoUsersBatched(currentUserId, targetUserId uuid.UUID, offset int) ([]models.Message, error) {
	var messages []models.Message
	rows, err := MyDB.Query("SELECT * FROM Messages WHERE (SenderId=? AND ReceiverId=?) OR (SenderId=? AND ReceiverId=?) LIMIT=10 OFFSET=? ORDER BY CreatedAt", currentUserId, targetUserId, targetUserId, currentUserId, offset)
	if err == sql.ErrNoRows {
		return messages, nil
	} else if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		var message models.Message
		err = rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.Content, &message.CreatedAt)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func FindAllMessages() ([]models.Message, error) {
	// var posts []models.Post
	// var finalPosts []models.Post
	var messages []models.Message
	rows, err := MyDB.Query("SELECT * FROM Messages")
	if err == sql.ErrNoRows {
		return messages, nil
	} else if err != nil {
		return messages, err
	}
	defer rows.Close()
	for rows.Next() {
		// 	var post models.Post
		var message models.Message
		err = rows.Scan(&message.Id, &message.SenderId, &message.ReceiverId, &message.Content, &message.CreatedAt)
		if err != nil {
			return messages, err
		}
		// 	post.Username, err = FindNicknameById(post.UserId)
		message.SenderName, err = FindNicknameById(message.SenderId)
		if err != nil {
			return messages, err
		}
		message.ReceiverName, err = FindNicknameById(message.ReceiverId)
		if err != nil {
			return messages, err
		}
		for i, v := range message.Id {
			if i == 4 || i == 6 || i == 8 || i == 10 {
				message.StringId += "-"
			}
			if len(strconv.FormatInt(int64(v), 16)) != 2 {
				message.StringId += "0"
			}
			message.StringId += strconv.FormatInt(int64(v), 16)
		}
		// 	posts = append(posts, post)
		messages = append(messages, message)
	}

	// for i := len(posts) - 1; i >= 0; i-- {
	// 	finalPosts = append(finalPosts, posts[i])
	// }
	return messages, nil
}
