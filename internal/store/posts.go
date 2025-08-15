package store

import (
	"context"
	"database/sql"
	"errors"
)

type Post struct {
	ID          int64  `json:"id"`
	Content     string `json:"content"`
	Title       string `json:"title"`
	AuthorEmail string `json:"author_email"`
	CreatedAt   string `json:"created_at"`
	PhotoURL    string `json:"photo_url"`
}
type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, author_email)
	VALUES ($1, $2, $3) RETURNING id, created_at, content, author_email
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query,
		post.Content, post.Title, post.AuthorEmail).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Content,
		&post.AuthorEmail,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, postId int64) (*Post, error) {
	query := `
	SELECT id, author_email, title, content, created_at 
	FROM posts 
	WHERE ID =  $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.AuthorEmail,
		&post.Title,
		&post.Content,
		&post.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}

	}
	return &post, nil
}

func (s *PostsStore) Delete(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE ID =  $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostsStore) Edit(ctx context.Context, post *Post) error {
	query := `UPDATE posts SET
content = $1, image = $2
WHERE ID = $3
RETURNING id;`

	err := s.db.QueryRowContext(ctx,
		query, post.Content, post.PhotoURL, post.ID).Scan(&post.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostsStore) GetList(ctx context.Context) ([]*Post, error) {
	query := `
	SELECT id, author_email, title, content, created_at, image
	FROM posts;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		err = rows.Scan(
			&post.ID,
			&post.AuthorEmail,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.PhotoURL)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
