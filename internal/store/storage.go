package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeOut      = 5 * time.Second
	ErrNotFound       = errors.New("record not found")
	MissionedAssigned = errors.New("mission assigned")
	TargetAmountError = errors.New("target amount error")
	ViolatePK         = errors.New("violate pk error")
	MissionCompleted  = errors.New("missiion completed")
)

type Storage struct {
	Users interface {
		CreateUser(username, email string) (*User, error) 
		GetUserByID(id int64) (*User, error)
		GetUserByEmail(email string) (*User, error)
	}
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, postId int64) (*Post, error)
		Delete(ctx context.Context, postID int64) error
		Edit(ctx context.Context, post *Post) error
		GetList(ctx context.Context) ([]*Post, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users: &UsersStore{db: db},
		Posts: &PostsStore{db: db},
	}
}
