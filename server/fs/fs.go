package fs

import (
	"bytes"
	"cloud-server/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var log = logger.New(" FS ", logger.Cyan)

type FileMeta struct {
	UploadName  string    `json:"uploadName"`
	FileUUID    string    `json:"fileUUID"`
	Size        uint64    `json:"size"`
	UploadedAt  time.Time `json:"uploadedAt"`
	ContentType string    `json:"contentType"`
}

type File struct {
	Meta FileMeta
	Data []byte
}

const META_SEPARATOR = "\n---FILEDATA---\n"

func FindFile(fileID string) (File, error) {
	// Validate file ID to prevent path traversal
	if !isValidFileID(fileID) {
		return File{}, os.ErrNotExist
	}

	data, err := os.ReadFile("storage/" + fileID)
	if err != nil {
		log.Error("Failed to read file: ", fileID, err)
		return File{}, err
	}
	file := File{}

	sepBytes := []byte(META_SEPARATOR)
	sepIndex := bytes.Index(data, sepBytes)
	if sepIndex == -1 {
		file.Data = data
		return file, nil
	}

	metaBytes := data[:sepIndex]
	file.Data = data[sepIndex+len(sepBytes):]

	if err := json.Unmarshal(metaBytes, &file.Meta); err != nil {
		log.Error("Failed to unmarshal metadata for file: ", fileID, err)
		return file, err
	}

	return file, nil
}

func (file File) SaveFile(db *sql.DB, ownerID string) (string, error) {
	id := uuid.New().String()
	file.Meta.FileUUID = id

	metaJSON, err := json.Marshal(file.Meta)
	if err != nil {
		log.Error("Failed to marshal file metadata: ", err)
		return "", err
	}

	var comb bytes.Buffer
	comb.Write(metaJSON)
	comb.WriteString(META_SEPARATOR)
	comb.Write(file.Data)

	temp := "storage/temp/" + id
	final := "storage/" + id

	err = os.WriteFile(temp, comb.Bytes(), 0600)
	if err != nil {
		log.Print("Failed to write temp file: ", temp, err)
		return "", err
	}

	if err := os.Rename(temp, final); err != nil {
		os.Remove(temp)
		log.Print("Failed to rename temp file: ", temp, final, err)
		return "", err
	}

	// Store file ownership in database if database connection and owner ID are provided
	if db != nil && ownerID != "" {
		_, err = db.Exec("INSERT INTO files (id, owner_id) VALUES (?, ?)", id, ownerID)
		if err != nil {
			// If database insertion fails, clean up the file
			log.Error("Failed to insert file: ", id, err)
			os.Remove("storage/" + id)
			return "", err
		}
	}

	return id, nil
}

func DeleteFile(fileID string, db *sql.DB) error {
	// Validate file ID to prevent path traversal
	if !isValidFileID(fileID) {
		return fmt.Errorf("invalid file ID")
	}

	// Delete from filesystem first
	err := os.Remove("storage/" + fileID)
	if err != nil && !os.IsNotExist(err) {
		log.Error("Failed to delete file: ", fileID, err)
		return err
	}

	// Delete from database if database connection is provided
	if db != nil {
		_, err = db.Exec("DELETE FROM files WHERE id = ?", fileID)
		if err != nil {
			log.Error("Failed to delete file: ", fileID, err)
			return err
		}
	}

	return nil
}

func isValidFileID(fileID string) bool {
	// Validate that fileID is a valid UUID and doesn't contain path traversal characters
	if len(fileID) != 36 {
		return false
	}

	// Check for path traversal characters
	if strings.Contains(fileID, "/") || strings.Contains(fileID, "..") || strings.Contains(fileID, "\\") {
		return false
	}

	// Basic UUID format validation (simple check, can be enhanced)
	parts := strings.Split(fileID, "-")
	return len(parts) == 5

	return true
}
