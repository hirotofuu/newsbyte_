package repository

import (
	"github.com/hirotofuu/newsbyte/internal/models"
)

type Databaserepo interface {
	GetUserByEmail(email string) (*models.User, error)
	InsertUser(user models.User) (int, error)
	AllUsers() ([]*models.User, error)
}
