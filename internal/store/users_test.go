package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserModel_FindBy(t *testing.T) {
	// Initialize the mock database and model
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := UserModel{db}

	// Define table-driven test cases
	tests := []struct {
		name          string
		field         string
		value         any
		mockBehavior  func()
		expectedUser  *User
		expectedError error
	}{
		{
			name:  "Successful Find by ID with Role",
			field: "id",
			value: int64(1),
			mockBehavior: func() {
				// Mock adding a row to the roles table
				sqlmock.NewRows([]string{"id", "role_name", "level", "description"}).
					AddRow(2, "Admin", 1, "Admin role description")

				// Mock adding a row to the users table with the associated role_id
				userRows := sqlmock.NewRows([]string{"id", "role_id", "username", "email", "password", "created_at", "last_login", "activated", "role.id", "role_name", "level", "description"}).
					AddRow(1, 2, "testuser", "test@example.com", []byte("hashedpassword"), time.Now(), nil, true, 2, "Admin", 1, "Admin role description")

				// Expect the roles row to be queried
				mock.ExpectQuery("SELECT users.id, role_id, username, email, password, created_at, last_login, activated, roles.* FROM users").
					WithArgs(int64(1)).
					WillReturnRows(userRows)
			},
			expectedUser: &User{
				ID:        1,
				RoleID:    2,
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  password{hash: []byte("hashedpassword")},
				Activated: true,
				Role: &Role{
					ID:          2,
					RoleName:    "Admin",
					Level:       1,
					Description: ptr("Admin role description"),
				},
				CreatedAt: time.Now(),
				LastLogin: nil,
			},
			expectedError: nil,
		},
		{
			name:  "Field Not Allowed",
			field: "invalid_field",
			value: "test",
			mockBehavior: func() {
				// No query expectation since this should fail validation before query execution.
			},
			expectedUser:  nil,
			expectedError: errors.New("field not allowed"),
		},
		{
			name:  "User Not Found",
			field: "id",
			value: int64(1),
			mockBehavior: func() {
				mock.ExpectQuery("SELECT users.id, role_id, username, email, password, created_at, last_login, activated, roles.* FROM users INNER JOIN roles").
					WithArgs(int64(1)).
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:  nil,
			expectedError: ErrNotFound,
		},
	}

	// Execute each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock behavior for the current test case
			tc.mockBehavior()

			// Call the FindBy method
			user, err := model.FindBy(context.Background(), tc.field, tc.value)

			// Assert that the error matches the expected error
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				// Compare the returned user with the expected user
				assert.Equal(t, tc.expectedUser.ID, user.ID)
				assert.Equal(t, tc.expectedUser.Role.ID, user.Role.ID)
			}

			// Check that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Helper function to create pointer for string values
func ptrString(s string) *string {
	return &s
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
		RoleID:   1,
	}

	err = user.Password.Set("password123")
	assert.NoError(t, err)

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

func TestUserModel_Update(t *testing.T) {
	// Initialize the mock database and model
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := UserModel{db}

	// Define table-driven test cases
	tests := []struct {
		name          string
		user          *User
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Successful Update",
			user: &User{
				ID:       1,
				Email:    "newemail@example.com",
				Username: "newusername",
				Password: password{hash: []byte("hashedpassword")},
			},
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users SET email = \\$1, username = \\$2, password = \\$3 WHERE id = \\$4").
					WithArgs("newemail@example.com", "newusername", []byte("hashedpassword"), int64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Conflict Error (Email already exists)",
			user: &User{
				ID:       1,
				Email:    "existingemail@example.com",
				Username: "newusername",
				Password: password{hash: []byte("hashedpassword")},
			},
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users SET email = \\$1, username = \\$2, password = \\$3 WHERE id = \\$4").
					WithArgs("existingemail@example.com", "newusername", []byte("hashedpassword"), int64(1)).
					WillReturnError(&pq.Error{Code: "23505"}) // Conflict due to unique violation
				mock.ExpectRollback()
			},
			expectedError: ErrConflict,
		},
		{
			name: "User Not Found",
			user: &User{
				ID:       1,
				Email:    "nonexistent@example.com",
				Username: "nonexistent",
				Password: password{hash: []byte("hashedpassword")},
			},
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users SET email = \\$1, username = \\$2, password = \\$3 WHERE id = \\$4").
					WithArgs("nonexistent@example.com", "nonexistent", []byte("hashedpassword"), int64(1)).
					WillReturnError(ErrNotFound) // No rows affected
				mock.ExpectRollback()
			},
			expectedError: ErrNotFound,
		},
	}

	// Execute each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock behavior for the current test case
			tc.mockBehavior()

			// Call the update method
			err = model.Update(context.Background(), tc.user)

			// Assert that the error matches the expected error
			assert.ErrorIs(t, err, tc.expectedError)

			// Check that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserModel_Delete(t *testing.T) {
	// Initialize the mock database and model
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	model := UserModel{db}

	// Define table-driven test cases
	tests := []struct {
		name          string
		id            int64
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Successful Delete",
			id:   1,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
					WithArgs(int64(1)).
					WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "User Not Found",
			id:   1,
			mockBehavior: func() {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
					WithArgs(int64(1)).
					WillReturnError(ErrNotFound) // No rows affected
				mock.ExpectRollback()
			},
			expectedError: ErrNotFound,
		},
	}

	// Execute each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock behavior for the current test case
			tc.mockBehavior()

			// Call the delete method
			err = model.Delete(context.Background(), tc.id)

			// Assert that the error matches the expected error
			assert.ErrorIs(t, err, tc.expectedError)

			// Check that all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
