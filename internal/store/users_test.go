package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserModel_FindBy(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := UserModel{db}

	// Valid user data for testing
	user := &User{
		ID:        1,
		RoleID:    2,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		LastLogin: nil,
	}

	// Define the test cases
	tests := []struct {
		name      string
		field     string
		value     any
		mockQuery func()
		wantUser  *User
		wantError error
	}{
		{
			name:  "Success - User Found by Username",
			field: "username",
			value: "testuser",
			mockQuery: func() {
				mock.ExpectQuery(`SELECT users.id, role_id, username, email, created_at, last_login FROM users`).
					WithArgs("testuser").
					WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "username", "email", "created_at", "last_login"}).
						AddRow(user.ID, user.RoleID, user.Username, user.Email, user.CreatedAt, user.LastLogin))
			},
			wantUser:  user,
			wantError: nil,
		},
		{
			name:  "Error - Field Not Allowed",
			field: "password",
			value: "somepassword",
			mockQuery: func() {
				// No mock required since the function will return an error before the query is executed
			},
			wantUser:  nil,
			wantError: errors.New("field not allowed"),
		},
		{
			name:  "Error - No User Found",
			field: "email",
			value: "nonexistent@example.com",
			mockQuery: func() {
				mock.ExpectQuery(`SELECT users.id, role_id, username, email, created_at, last_login FROM users`).
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantUser:  nil,
			wantError: ErrNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock expectations
			tc.mockQuery()

			// Call the function
			gotUser, err := model.FindBy(context.Background(), tc.field, tc.value)

			// Assertions
			if tc.wantError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantUser, gotUser)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := UserModel{db}

	// Define the user data
	user := &User{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123", // This should be encrypted in the Insert function
		RoleID:   1,
	}

	// Define the test cases
	tests := []struct {
		name      string
		mockQuery func()
		wantError error
	}{
		{
			name: "Success - User Created",
			mockQuery: func() {
				// Mock the "FindBy" function call (no existing user found)
				mock.ExpectQuery(`SELECT users.id, role_id, username, email, created_at, last_login FROM users`).
					WithArgs(user.Email).
					WillReturnError(sql.ErrNoRows)

				// Mock successful user creation
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Username, user.Email, sqlmock.AnyArg(), user.RoleID). // sqlmock.AnyArg() for the hashed password
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantError: nil,
		},
		{
			name: "Error - User Already Exists (Conflict)",
			mockQuery: func() {
				// Mock the "FindBy" function call (existing user found)
				mock.ExpectQuery(`SELECT users.id, role_id, username, email, created_at, last_login FROM users`).
					WithArgs(user.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "username", "email"}).
						AddRow(1, 1, "existinguser", "newuser@example.com"))

				// No insert query expected since user already exists
			},
			wantError: ErrConflict,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock expectations
			tc.mockQuery()

			// Call the Insert function
			err = model.Insert(context.Background(), user)

			// Assert the expected error
			if tc.wantError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Ensure all mock expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
