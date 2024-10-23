package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

func AddSession(userId uuid.UUID) error {
	statement, err := MyDB.Prepare("INSERT INTO Sessions VALUES (?,?,?)")
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = statement.Exec(id, userId, time.Now())
	return err
}

func DeleteSession(userId uuid.UUID) error {
	statement, err := MyDB.Prepare("DELETE FROM Sessions WHERE UserId=?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(userId)
	return err
}

func FindSessionByUserId(userId uuid.UUID) (bool, error) {
	var user uuid.UUID
	err := MyDB.QueryRow("SELECT UserId FROM Sessions WHERE UserId=?", userId).Scan(&user)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	fmt.Println(userId)
	return true, nil
}
