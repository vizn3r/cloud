package fs

import (
	"cloud-server/db"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreateShare(data *db.DB, fileID string) (string, error) {
	shareID := uuid.New().String()
	_, err := data.Connection.Exec(db.Q_SHARE_INSERT, shareID, fileID, time.Now().Add(24*time.Hour))
	return shareID, err
}

func FindFileByShare(data *db.DB, shareID string) (File, error) {
	var fileID string
	var downloads int
	var expiresAt *time.Time

	err := data.Connection.QueryRow(db.Q_SHARE_FIND_BY_ID, shareID).Scan(&fileID, &downloads, &expiresAt)
	if err != nil {
		return File{}, err
	}

	if expiresAt != nil && time.Now().After(*expiresAt) {
		data.Connection.Exec(db.Q_SHARE_DELETE, shareID)
		return File{}, fmt.Errorf("share expired")
	}
	return FindFile(fileID)
}
