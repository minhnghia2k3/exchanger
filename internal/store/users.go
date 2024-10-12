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
	CreateAndInvite(ctx context.Context, token string, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

type User struct {
	ID        int64      `json:"id"`
	RoleID    int64      `json:"role_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  password   `json:"-"`
	Activated bool       `json:"activated"`
	Role      *Role      `json:"role,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login"`
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

type password struct {
	plain string
	hash  []byte
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

	query := fmt.Sprintf(`SELECT users.id, role_id, username, email, password, created_at, last_login, activated, roles.*
FROM users INNER JOIN roles ON role_id = roles.id
WHERE users.%s=$1`, field)

	ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
	defer cancel()

	user.Role = &Role{}

	err := m.db.QueryRowContext(ctx, query, value).Scan(
		&user.ID,
		&user.RoleID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.LastLogin,
		&user.Activated,
		&user.Role.ID,
		&user.Role.RoleName,
		&user.Role.Level,
		&user.Role.Description,
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

	args := []any{user.RoleID, user.Username, user.Email, user.Password.hash}

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

func (m *UserModel) CreateAndInvite(ctx context.Context, token string, user *User) error {
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
		err := m.Insert(ctx, user)
		if err != nil {
			return err
		}

		// TODO: send user activation token to email
		return nil
	})
}

func (m *UserModel) Update(ctx context.Context, user *User) error {
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
		query := `UPDATE users SET email = $1, username = $2, password = $3 WHERE id = $4`

		ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
		defer cancel()

		args := []any{user.Email, user.Username, user.Password.hash, user.ID}

		_, err := tx.ExecContext(ctx, query, args...)

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
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
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
	})
}

func (p *password) Set(plain string) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plain), 10)
	if err != nil {
		return err
	}

	p.plain = plain
	p.hash = hashedPassword

	return nil
}

func (p *password) Verify() error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(p.plain))
}
