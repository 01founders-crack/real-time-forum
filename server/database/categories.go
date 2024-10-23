package database

import (
	"rtforum/server/models"

	"github.com/gofrs/uuid"
)

// takes in category details and sends the information to corresponding column in 'Categories' table
func AddCategory(category models.Category) error {
	statement, err := MyDB.Prepare("INSERT INTO Categories VALUES (?,?,?)")
	if err != nil {
		return err
	}
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}
	_, err = statement.Exec(id, category.Name, category.Desc)
	return err
}
