package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"time"
)

// Authorization interface defines methods for user authentication and session management.
type Authorization interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	AddSessionToken(email, token string, expiresAt time.Time) error
	GetSessionToken(token string) (models.User, error)
	DeleteSessionToken(token string) error
}

// AuthStorage is a struct that implements the Authorization interface.
type AuthStorage struct {
	db *sql.DB
}

// NewAuthSqlite returns a new AuthStorage instance.
func NewAuthSqlite(db *sql.DB) *AuthStorage {
	return &AuthStorage{db: db}
}

// CreateUser creates a new user in the database.
func (r *AuthStorage) CreateUser(user *models.User) error {
	query := fmt.Sprintf("INSERT INTO user (username, email, password) values ($1, $2, $3)")
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail retrieves a user from the database by email.
func (s *AuthStorage) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT id, email, username, password FROM user WHERE email=$1;`
	row := s.db.QueryRow(query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by email: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user from the database by username.
func (s *AuthStorage) GetUserByUsername(username string) (models.User, error) {
	query := `SELECT id, email, username, password FROM user WHERE username=$1;`
	row := s.db.QueryRow(query, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by username: %w", err)
	}
	return user, nil
}

// AddSessionToken adds a session token to a user in the database.
func (s *AuthStorage) AddSessionToken(email, token string, expiresAt time.Time) error {
	query := `UPDATE user SET token = $1, expiresAt = $2 WHERE email = $3;`
	_, err := s.db.Exec(query, token, expiresAt, email)
	if err != nil {
		return fmt.Errorf("storage: save session token: %w", err)
	}
	return nil
}

// GetSessionToken retrieves a user from the database by session token.
func (s *AuthStorage) GetSessionToken(token string) (models.User, error) {
	query := `SELECT id, email, username, password, token, expiresAt FROM user WHERE token=$1;`

	row := s.db.QueryRow(query, token)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Token, &user.ExpiresAt)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by session token: %w", err)
	}
	return user, nil
}

// DeleteSessionToken removes a session token from a user in the database.
func (s *AuthStorage) DeleteSessionToken(token string) error {
	query := `UPDATE user SET token = NULL, expiresAt = NULL WHERE token = $1;`
	_, err := s.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("storage: delete session token: %w", err)
	}
	return nil
}
