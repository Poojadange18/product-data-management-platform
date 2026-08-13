package repository

import (
	"database/sql"
	"errors"

	"product-data-management-platform/model"
)

var ErrDuplicate = errors.New("duplicate record")

type UserRepository struct{ DB *sql.DB }

func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{DB: db} }

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.QueryRow(`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`, user.Email, user.PasswordHash).Scan(&user.ID)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.DB.QueryRow(`SELECT id, email, password_hash FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return user, nil
}
