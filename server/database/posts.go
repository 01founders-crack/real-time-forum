package database

import (
	"database/sql"
	"fmt"
	"rtforum/server/models"
	"strconv"

	"github.com/gofrs/uuid"
)

// takes in post details and sends the information to corresponding column in 'Posts' table
func AddPost(post models.Post) error {
	statement, err := MyDB.Prepare("INSERT INTO Posts VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println("7")
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Println("8")
		return err
	}
	_, err = statement.Exec(id, post.UserId, post.Title, post.Category, post.Content)
	return err
}

func FindAllPosts() ([]models.Post, error) {
	var posts []models.Post
	var finalPosts []models.Post
	rows, err := MyDB.Query("SELECT * FROM Posts")
	if err == sql.ErrNoRows {
		fmt.Println("9")
		return posts, nil
	} else if err != nil {
		fmt.Println("10")
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Category, &post.Content)
		if err != nil {
			fmt.Println("11")
			return posts, err
		}
		post.Username, err = FindNicknameById(post.UserId)
		if err != nil {
			return posts, err
		}
		post.Comments, err = FindCommentsByPostId(post.Id)
		if err != nil {
			fmt.Println("12")
			return posts, err
		}
		for i, v := range post.Id {
			if i == 4 || i == 6 || i == 8 || i == 10 {
				post.StringId += "-"
			}
			if len(strconv.FormatInt(int64(v), 16)) != 2 {
				post.StringId += "0"
			}
			post.StringId += strconv.FormatInt(int64(v), 16)
		}
		posts = append(posts, post)
	}

	for i := len(posts) - 1; i >= 0; i-- {
		finalPosts = append(finalPosts, posts[i])
	}
	return finalPosts, nil
}

func FindPostsByCategory(category string) ([]models.Post, error) {
	var posts []models.Post
	rows, err := MyDB.Query("SELECT * FROM Posts WHERE Category=?", category)
	if err == sql.ErrNoRows {
		fmt.Println("13")
		return posts, nil
	} else if err != nil {
		fmt.Println("14")
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.Id, &post.UserId, &post.Title, &post.Category, &post.Content)
		if err != nil {
			fmt.Println("15")
			return posts, err
		}
		post.Username, err = FindNicknameById(post.UserId)
		if err != nil {
			return posts, err
		}
		post.Comments, err = FindCommentsByPostId(post.Id)
		if err != nil {
			fmt.Println("16")
			return posts, err
		}
		for _, v := range post.Id {
			post.StringId += string(v)
			fmt.Println(v)
			fmt.Println(string(v))
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func FindPostsByTitle(title string) (uuid.UUID, error) {
	var postId uuid.UUID
	err := MyDB.QueryRow("SELECT Id FROM Posts WHERE Title=?", title).Scan(&postId)
	return postId, err
}
