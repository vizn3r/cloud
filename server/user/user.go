package user

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"cloud-server/db"
	"cloud-server/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var log = logger.New("USER", logger.Purple)

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

	salt := generateToken()

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password for email:", email, err)
		return "", fmt.Errorf("failed to hash password")
	}

	userID := uuid.New().String()
	_, err = data.Connection.Exec(db.Q_USER_CREATE, userID, email, string(hash), salt)
	if err != nil {
		log.Error("Failed to create user:", email, err)
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	return userID, nil
}

func AuthenticateUser(data *db.DB, email, password string) (string, error) {
	var userID, storedHash, salt string
	var createdAt time.Time

	err := data.Connection.QueryRow(db.Q_USER_FIND_BY_EMAIL, email).Scan(&userID, &email, &storedHash, &salt, &createdAt)
	if err != nil {
		log.Error("Failed to find user by email:", email, err)
		return "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password+salt))
	if err != nil {
		log.Print("Password mismatch for user", email)
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
		log.Error("Failed to create session for user:", userID, err)
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
		log.Error("Invalid session token:", token, err)
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
		log.Error("Failed to find user by ID:", userID, err)
		return User{}, fmt.Errorf("user not found")
	}

	return User{
		ID:        userID,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}

func GetUserFiles(data *db.DB, userID string) ([]string, error) {
	rows, err := data.Connection.Query(db.Q_FILE_FIND_BY_OWNER, userID)
	if err != nil {
		log.Error("Failed to query user files for:", userID, err)
		return nil, fmt.Errorf("failed to get user files")
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var fileID string
		var uploadedAt time.Time
		var updatedAt sql.NullTime // Use nullable type for updated_at which can be NULL
		if err := rows.Scan(&fileID, &uploadedAt, &updatedAt); err != nil {
			log.Error("Failed to scan file row:", err)
			continue
		}
		files = append(files, fileID)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating over files:", err)
		return nil, fmt.Errorf("error reading files")
	}

	return files, nil
}

func generateToken() string {
	// Loop IN CASE rand.Read() errors out (no chance, but redundancy is cool)
	for {
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err == nil {
			return hex.EncodeToString(bytes)
		}
	}
}
