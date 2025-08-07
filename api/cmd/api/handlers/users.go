package handlers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type UserInput struct {
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

const USER_FIELDS = "id, name, email, created_at, updated_at"

func UsersGetTx(tx *sql.Tx, id uuid.UUID) (*User, error) {
	user := User{}
	s := fmt.Sprintf(`SELECT %s FROM users WHERE id=$1`, USER_FIELDS)
	err := tx.QueryRow(s, id).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func UsersCreateTx(tx *sql.Tx, input *UserInput) (*User, error) {
	user := &User{}
	s := fmt.Sprintf(`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING %s`, USER_FIELDS)
	err := tx.QueryRow(s, input.Name, input.Email).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	return user, err
}

func UsersUpdateTx(tx *sql.Tx, id uuid.UUID, input *UserInput) (*User, error) {
	user := &User{}
	s := fmt.Sprintf(`UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING %s`, USER_FIELDS)
	err := tx.QueryRow(s, input.Name, input.Email, id).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func UsersDeleteTx(tx *sql.Tx, id uuid.UUID) (*User, error) {
	user := &User{}
	s := fmt.Sprintf(`DELETE FROM users WHERE id=$1 RETURNING %s`, USER_FIELDS)
	err := tx.QueryRow(s, id).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func UsersGetAllTx(tx *sql.Tx) ([]*User, error) {
	s := fmt.Sprintf(`SELECT %s FROM users ORDER BY created_at DESC`, USER_FIELDS)
	rows, err := tx.Query(s)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, err
}
