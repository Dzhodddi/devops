package store

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersStore struct {
	db *sql.DB
}

func (store *UsersStore) GetUserByID(id int64) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, created_at FROM users WHERE id = $1"
	err := store.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (store *UsersStore) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, email, created_at FROM users WHERE email = $1"
	err := store.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (store *UsersStore) CreateUser(username, email string) (*User, error) {
	user := &User{}
	query := "INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, username, email, created_at"
	err := store.db.QueryRow(query, username, email).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ViolatePK
		default:
			return nil, err
		}
	}
	return user, nil
}
