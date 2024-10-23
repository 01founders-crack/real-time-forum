package database

import (
	"database/sql"
	"fmt"
	"rtforum/server/models"

	"github.com/gofrs/uuid"
)

// takes in comment details and sends the information to corresponding column in 'Comments' table
func AddComment(comment models.Comment) error {
	statement, err := MyDB.Prepare("INSERT INTO Comments VALUES (?,?,?,?)")
	if err != nil {
		fmt.Println("1")
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Println("2")
		return err
	}
	_, err = statement.Exec(id, comment.PostId, comment.UserId, comment.Content)
	return err
}

func FindCommentsByPostId(postId uuid.UUID) ([]models.Comment, error) {
	var finalComments []models.Comment
	var comments []models.Comment
	rows, err := MyDB.Query("SELECT * FROM Comments WHERE PostId=?", postId)
	if err == sql.ErrNoRows {
		fmt.Println("3")
		return comments, nil
	} else if err != nil {
		fmt.Println("4")
		return comments, err
	}

	defer rows.Close()
	for rows.Next() {
		var comment models.Comment
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.UserId, &comment.Content)
		if err != nil {
			fmt.Println("5")
			return comments, err
		}

		comment.Username, err = FindNicknameById(comment.UserId)
		comments = append(comments, comment)
	}

	for i := len(comments) - 1; i >= 0; i-- {
		finalComments = append(finalComments, comments[i])
	}

	return finalComments, nil
}
