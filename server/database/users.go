package database

import (
	"database/sql"
	"errors"
	"fmt"
	"rtforum/server/models"
	"rtforum/server/utils"

	"github.com/gofrs/uuid"
)

// takes in registration details and sends the information to corresponding column in 'Users' table
func AddUser(user models.User) error {
	statement, err := MyDB.Prepare("INSERT INTO Users VALUES (?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = statement.Exec(id, user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, user.Password)
	return err
}

// checks if nickname already exists, result can be used to confirm log in details or in registration
func NicknameExists(nickname string) (bool, error) {
	found := ""
	err := MyDB.QueryRow("SELECT Nickname FROM Users WHERE Nickname=?", nickname).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// checks if email already exists, result can be used to confirm log in details or in registration
func EmailExists(email string) (bool, error) {
	found := ""
	err := MyDB.QueryRow("SELECT Email FROM Users WHERE Email=?", email).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func IsLoginValid(login, enteredPassword string) (bool, error) {
	if login == "" || enteredPassword == "" {
		return false, errors.New("missing credentials")
	}
	actualPass := ""
	exists, err := NicknameExists(login)
	if err != nil {
		return false, err
	}
	if exists {
		err = MyDB.QueryRow("SELECT Password FROM Users WHERE Nickname=?", login).Scan(&actualPass)
		if err != nil {
			return false, err
		}
	} else {
		exists, err = EmailExists(login)
		if err != nil {
			return false, err
		}
		if exists {
			err = MyDB.QueryRow("SELECT Password FROM Users WHERE Email=?", login).Scan(&actualPass)
			if err != nil {
				return false, err
			}
		} else {
			// Return false explicitly when the login does not match any user
			return false, nil
		}
	}

	// Validate the entered password
	if utils.ComparePasswords(actualPass, enteredPassword) {
		return true, nil
	}

	return false, nil
}

func FindNicknameById(userId uuid.UUID) (string, error) {
	var nickname string
	err := MyDB.QueryRow("SELECT Nickname FROM Users WHERE Id=?", userId).Scan(&nickname)
	return nickname, err
}

func FindIdByNickname(nickname string) (uuid.UUID, error) {
	var userId uuid.UUID
	err := MyDB.QueryRow("SELECT Id FROM Users WHERE Nickname=?", nickname).Scan(&userId)
	return userId, err
}

func FindIdByEmail(email string) (uuid.UUID, error) {
	var userId uuid.UUID
	err := MyDB.QueryRow("SELECT Id FROM Users WHERE Email=?", email).Scan(&userId)
	return userId, err
}

func FindIdByLogin(login string) (uuid.UUID, error) {
	var userId uuid.UUID
	userId, err := FindIdByNickname(login)
	if err == sql.ErrNoRows {
		err = nil
		userId, err = FindIdByEmail(login)
	}
	return userId, err
}

func GetAllUsers() ([]models.User, error) {
	var userList []models.User
	rows, err := MyDB.Query("SELECT Id, Nickname FROM Users ORDER BY Nickname")
	if err != nil {
		fmt.Println(err)
		return userList, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Id, &user.Nickname)
		if err != nil {
			return userList, err
		}
		userList = append(userList, user)
	}

	return userList, err
}
