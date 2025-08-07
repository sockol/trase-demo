package handlers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Post struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Content   string     `json:"content" db:"content"`
	UserId    uuid.UUID  `json:"user_id" db:"user_id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type PostInput struct {
	Title   string    `json:"title" db:"title"`
	Content string    `json:"content" db:"content"`
	UserId  uuid.UUID `json:"user_id" db:"user_id"`
}

const POST_FIELDS = "id, title, content, user_id, created_at, updated_at"

func PostsGetTx(tx *sql.Tx, id uuid.UUID) (*Post, error) {
	post := Post{}
	s := fmt.Sprintf(`SELECT %s FROM posts WHERE id=$1`, POST_FIELDS)
	err := tx.QueryRow(s, id).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &post, err
}

func PostsCreateTx(tx *sql.Tx, input *PostInput) (*Post, error) {
	post := &Post{}
	s := fmt.Sprintf(`INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING %s`, POST_FIELDS)
	err := tx.QueryRow(s, input.Title, input.Content, input.UserId).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.UpdatedAt)
	return post, err
}

func PostsUpdateTx(tx *sql.Tx, id uuid.UUID, input *PostInput) (*Post, error) {
	post := &Post{}
	s := fmt.Sprintf(`UPDATE posts SET title=$1, content=$2, user_id=$3 WHERE id = $4 RETURNING %s`, POST_FIELDS)
	err := tx.QueryRow(s, input.Title, input.Content, input.UserId, id).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return post, err
}

func PostsDeleteTx(tx *sql.Tx, id uuid.UUID) (*Post, error) {
	post := &Post{}
	s := fmt.Sprintf(`DELETE FROM posts WHERE id=$1 RETURNING %s`, POST_FIELDS)
	err := tx.QueryRow(s, id).Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return post, err
}

func PostsGetAllTx(tx *sql.Tx) ([]*Post, error) {
	s := fmt.Sprintf(`SELECT %s FROM posts ORDER BY created_at DESC`, POST_FIELDS)
	rows, err := tx.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		post := Post{}
		err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.UserId, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, err
}
