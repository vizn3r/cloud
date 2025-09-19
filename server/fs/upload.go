package fs

import (
	"time"

	"cloud-server/db"

	"github.com/google/uuid"
)

type UploadSession struct {
	ID     string        `json:"ID,omitempty"`
	UserID string        `json:"userID,omitempty"`
	File   *File         `json:"file,omitempty"`
	Chunks []UploadChunk `json:"chunks,omitempty"`

	ChunkSize uint `json:"chunkSize,omitempty"`
	ChunkNum  uint `json:"chunkNum,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

/*
*
* Example POST /init request:
*
* {
*   "chunkSize": 1024,
*   "chunkNum": 10,
*   "file": {
*     "meta": {
*       "uploadName": "test.txt",
*       "contentType": "text/plain"
*     }
*   }
* }
*
 */

type UploadChunk struct {
	ID        string `json:"ID,omitempty"`
	SessionID string `json:"sessionID,omitempty"`

	Size uint   `json:"size"`
	Num  uint   `json:"num"`
	Data []byte `json:"data,omitempty"`
}

func CreateUploadSession(data *db.DB, userID string, duration time.Duration) (*UploadSession, error) {
	session := &UploadSession{
		ID:        uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}

	file := &File{
		Meta: FileMeta{
			ID: uuid.New().String(),
		},
	}
	session.File = file

	if _, err := data.Connection.Exec(db.Q_UPLOAD_INSERT, session.ID, userID, file.Meta.ID, session.ExpiresAt); err != nil {
		return nil, err
	}

	return session, nil
}
