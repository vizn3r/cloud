package fs

import (
	"cloud-server/db"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateShare(data *db.DB, fileID string, duration time.Duration) (string, error) {
	shareID := uuid.New().String()
	_, err := data.Connection.Exec(db.Q_SHARE_INSERT, shareID, fileID, time.Now().Add(duration))
	if err != nil {
		log.Printf("Failed to create share for file %s: %v", fileID, err)
	}
	return shareID, err
}

func FindFileByShare(data *db.DB, shareID string) (File, error) {
	// Validate share ID format
	if len(shareID) != 36 || strings.Contains(shareID, "/") || strings.Contains(shareID, "..") {
		return File{}, fmt.Errorf("invalid share ID")
	}

	var fileID string
	var downloads int
	var expiresAt *time.Time

	err := data.Connection.QueryRow(db.Q_SHARE_FIND_BY_ID, shareID).Scan(&fileID, &downloads, &expiresAt)
	if err != nil {
		log.Printf("Failed to find share %s: %v", shareID, err)
		return File{}, err
	}

	if expiresAt != nil && time.Now().After(*expiresAt) {
		data.Connection.Exec(db.Q_SHARE_DELETE, shareID)
		return File{}, fmt.Errorf("share expired")
	}
	return FindFile(fileID)
}
