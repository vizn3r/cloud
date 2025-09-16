package user

import (
	"cloud-server/db"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

type Session struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
}

func CreateUser(data *db.DB, email, password string) (string, error) {
	// Validate email format
	if len(email) < 3 || len(email) > 255 {
		return "", fmt.Errorf("invalid email")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}

	userID := uuid.New().String()
	_, err = data.Connection.Exec(db.Q_USER_CREATE, userID, email, string(hash))
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	return userID, nil
}

func AuthenticateUser(data *db.DB, email, password string) (string, error) {
	var userID, storedHash string
	var createdAt time.Time

	err := data.Connection.QueryRow(db.Q_USER_FIND_BY_EMAIL, email).Scan(&userID, &email, &storedHash, &createdAt)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return userID, nil
}

func CreateSession(data *db.DB, userID string, duration time.Duration) (string, error) {
	sessionID := uuid.New().String()
	token := generateToken()
	expiresAt := time.Now().Add(duration)

	_, err := data.Connection.Exec(db.Q_SESSION_CREATE, sessionID, userID, token, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create session")
	}

	return token, nil
}

func ValidateSession(data *db.DB, token string) (string, error) {
	var userID string
	var expiresAt time.Time

	// Clean up expired sessions first
	data.Connection.Exec(db.Q_SESSION_DELETE_EXPIRED)

	err := data.Connection.QueryRow(db.Q_SESSION_FIND_BY_TOKEN, token).Scan(&userID, &expiresAt)
	if err != nil {
		return "", fmt.Errorf("invalid session")
	}

	if time.Now().After(expiresAt) {
		data.Connection.Exec(db.Q_SESSION_DELETE, token)
		return "", fmt.Errorf("session expired")
	}

	return userID, nil
}

func GetUserByID(data *db.DB, userID string) (User, error) {
	var email string
	var createdAt time.Time

	err := data.Connection.QueryRow(db.Q_USER_FIND_BY_ID, userID).Scan(&userID, &email, &createdAt)
	if err != nil {
		return User{}, fmt.Errorf("user not found")
	}

	return User{
		ID:        userID,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
