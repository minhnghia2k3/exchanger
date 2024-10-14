package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type IUsers interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	CreateAndInvite(ctx context.Context, token string, user *User) error
	Delete(ctx context.Context, id int64) error
	Activate(ctx context.Context, token string) error
	Login(ctx context.Context, email, password string) (*User, error)
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

type UserStorage struct {
	db *sql.DB
}

type password struct {
	plain string
	hash  []byte
}

func (m *UserStorage) GetByID(ctx context.Context, id int64) (*User, error) {
	var user User

	query := `SELECT users.id, role_id, username, email, password, created_at, last_login, activated, roles.*
	FROM users INNER JOIN roles ON role_id = roles.id
	WHERE users.id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
	defer cancel()

	user.Role = &Role{}

	err := m.db.QueryRowContext(ctx, query, id).Scan(
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

func (m *UserStorage) CreateAndInvite(ctx context.Context, token string, user *User) error {
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
		// Create new user
		err := m.create(ctx, user)
		if err != nil {
			return err
		}

		// Store invitation token to database
		expiryDuration := time.Now().Add(3 * 24 * time.Hour)
		if err = m.createUserInvitation(ctx, user.ID, token, expiryDuration); err != nil {
			return err
		}

		return nil
	})
}

func (m *UserStorage) Delete(ctx context.Context, id int64) error {
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
		query := `DELETE FROM users WHERE id = $1`
		ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
		defer cancel()

		_, err := m.db.ExecContext(ctx, query, id)
		log.Println("error:", err)
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

func (m *UserStorage) Activate(ctx context.Context, token string) error {
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
		// 1. get user by token
		user, err := m.getUserByInvitation(ctx, token)
		if err != nil {
			return err
		}

		// 2. update activation status
		user.Activated = true
		if err = m.update(ctx, user); err != nil {
			return err
		}

		// 3. Delete user invitations
		if err = m.deleteUserInvitation(ctx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (m *UserStorage) Login(ctx context.Context, email, password string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
	defer cancel()

	var user User
	user.Role = &Role{}

	err := m.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password.hash)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUnauthorized
		default:
			return nil, err
		}
	}

	// Validate password
	if err = user.Password.Verify(password); err != nil {
		return nil, ErrUnauthorized
	}

	return &user, nil
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

func (p *password) Verify(plain string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plain))
}

func (m *UserStorage) create(ctx context.Context, user *User) error {
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

func (m *UserStorage) update(ctx context.Context, user *User) error {
	query := `UPDATE users SET activated = $1 WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), QueryContextTimeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, user.Activated, user.ID)

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

func (m *UserStorage) getUserByInvitation(ctx context.Context, token string) (*User, error) {
	var user User

	query := `SELECT id, activated FROM users 
    INNER JOIN user_invitations ON user_invitations.user_id = users.id
	WHERE user_invitations.token = $1 AND user_invitations.expire_at > $2`

	sum := sha256.Sum256([]byte(token))
	hashedToken := hex.EncodeToString(sum[:])

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, hashedToken, time.Now()).Scan(&user.ID, &user.Activated)

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

func (m *UserStorage) createUserInvitation(ctx context.Context, userID int64, token string, expiry time.Time) error {
	query := `INSERT INTO user_invitations(user_id, token, expire_at)
	VALUES($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, userID, token, expiry)

	if err != nil {
		return err
	}

	return nil

}

func (m *UserStorage) deleteUserInvitation(ctx context.Context, userID int64) error {
	return withTx(ctx, m.db, func(tx *sql.Tx) error {
		query := `DELETE FROM user_invitations WHERE user_id = $1`

		ctx, cancel := context.WithTimeout(ctx, QueryContextTimeout)
		defer cancel()

		_, err := tx.ExecContext(ctx, query, userID)
		if err != nil {
			return err
		}

		return nil
	})
}
