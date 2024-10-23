package database

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

func NumberOfMessagesFromOneUser(currentUser, targetUser uuid.UUID) (int, error) {
	var notificationNumber int
	err := MyDB.QueryRow("SELECT NumberOfUnread FROM Notifications WHERE CurrentUserId=? AND SenderId=?", currentUser, targetUser).Scan(&notificationNumber)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return notificationNumber, err
}

func TotalNumberOfMessages(currentUser uuid.UUID) (int, error) {
	var notificationNumber int
	rows, err := MyDB.Query("SELECT NumberOfUnread FROM Notifications WHERE CurrentUserId=?", currentUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	}
	defer rows.Close()
	for rows.Next() {
		var num int
		err = rows.Scan(&num)
		notificationNumber += num
	}
	return notificationNumber, err
}

func NotificationList(currentUserId uuid.UUID) (map[string]int, error) {
	allNotifications := make(map[string]int)
	rows, err := MyDB.Query("SELECT SenderId, NumberOfUnread FROM Notifications WHERE CurrentUserId=?", currentUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()
	for rows.Next() {
		var senderId uuid.UUID
		var num int
		err = rows.Scan(&senderId, &num)
		if err != nil {
			return nil, err
		}
		nickname, err := FindNicknameById(senderId)
		if err != nil {
			return nil, err
		}
		allNotifications[nickname] = num
	}

	return allNotifications, nil
}

func AddNotification(senderId, receiverId uuid.UUID) error {
	statement, err := MyDB.Prepare("INSERT INTO Notifications VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	previousNotifs, err := NumberOfMessagesFromOneUser(receiverId, senderId)
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = statement.Exec(id, receiverId, senderId, previousNotifs+1)
	return err
}

func DeleteNotification(currentUserId, targetId uuid.UUID) error {
	statement, err := MyDB.Prepare("DELETE FROM Notifications WHERE CurrentUserId=? AND SenderId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(currentUserId, targetId)
	return err
}
