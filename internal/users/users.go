package users

import (
	"database/sql"
	"errors"
	"forum/internal/database"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

type Session struct {
	UserID    string
	ExpiresAt time.Time
}

const SessionDuration = 2 * time.Hour // Example duration

var sessionStore = make(map[string]Session)
var storeLock = sync.RWMutex{} // To make map access thread-safe

var db *sql.DB

func getUserByEmail(email string) (*User, error) {
	db := database.GetDB()
	user := &User{}
	err := db.QueryRow("SELECT id, username, password, email, created_at FROM users WHERE email = ?", email).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getUserByEmailOrUsername(email, username string) (*User, error) {
	db := database.GetDB()
	user := &User{}
	err := db.QueryRow("SELECT id, username, password, email, created_at FROM users WHERE email = ? OR username = ?", email, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByEmailOrUsername(emailOrUsername string) (*User, error) {
	db := database.GetDB()
	user := &User{}
	err := db.QueryRow("SELECT id, username, password, email, created_at FROM users WHERE email = ? OR username = ?", emailOrUsername, emailOrUsername).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(username, password, email string) (*User, error) {
	db := database.GetDB()
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	result, err := db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, hashedPassword, email)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &User{ID: int(id), Username: username, Password: hashedPassword, Email: email, CreatedAt: time.Now()}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(username, password, email string) error {
	// Check if username or email already exists
	db := database.GetDB()
	existingUser, err := getUserByEmailOrUsername(email, username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username or email already taken")
	}

	// Hash the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	// Save user to the database
	_, err = db.Exec("INSERT INTO users (username, password, email) VALUES (?, ?, ?)", username, hashedPassword, email)
	return err
}

func Login(email, password string) (*User, error) {
	user, err := getUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("email not registered")
	}

	// Verify password
	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("incorrect password")
	}

	// Create a session (implementation depends on your session management strategy)

	return user, nil
}
