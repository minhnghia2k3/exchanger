package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IUser interface {
	FindBy(ctx context.Context, field string, value any) (*User, error)
	Insert(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

type User struct {
	ID        int64      `json:"id"`
	RoleID    int64      `json:"role_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  Password   `json:"password"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login"`
	Role
}

type Role struct {
	ID          int64   `json:"id"`
	RoleName    string  `json:"role_name"`
	Level       int64   `json:"level"`
	Description *string `json:"description"`
}

type UserModel struct {
	db *sql.DB
}

type Password struct {
	Plain string
	Hash  []byte
}

func (m *UserModel) FindBy(ctx context.Context, field string, value any) (*User, error) {
	var user User

	allowField := map[string]bool{
		"id":       true,
		"role_id":  true,
		"username": true,
		"email":    true,
	}

	if !allowField[field] {
		return nil, errors.New("field not allowed")
	}

	query := fmt.Sprintf(`SELECT users.id, role_id, username, email, created_at, last_login
FROM users INNER JOIN roles ON role_id = roles.id
WHERE users.%s=$1`, field)

	ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, value).Scan(
		&user.ID,
		&user.RoleID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
func (m *UserModel) Insert(ctx context.Context, user *User) error {
	query := `INSERT INTO users(role_id, username, email, password)
VALUES($1, $2, $3, $4) RETURNING id, created_at, last_login`

	ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
	defer cancel()

	args := []any{user.RoleID, user.Username, user.Email, user.Password.Hash}

	err := m.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return ErrConflict
		default:
			return err
		}
	}

	return nil
}
func (m *UserModel) Update(ctx context.Context, user *User) error {
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
		query := `UPDATE users SET email = $1, username = $2, password = $3 WHERE id = $4`

		ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
		defer cancel()

		_, err := tx.ExecContext(ctx, query, user.Password.Hash, user.ID)

		if err != nil {
			var pqErr *pq.Error
			switch {
			case errors.As(err, &pqErr):
				return ErrConflict
			case errors.Is(err, sql.ErrNoRows):
				return ErrNotFound
			default:
				return err
			}
		}

		return nil
	})
}
func (m *UserModel) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
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

func (p *Password) Set(plain string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plain), 10)
	if err != nil {
		return err
	}

	p.Plain = plain
	p.Hash = hashedPassword

	return nil
}

func (p *Password) Verify() error {
	return bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(p.Plain))
}
