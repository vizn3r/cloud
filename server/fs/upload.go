package fs

import (
	"math/bits"
	"os"
	"time"

	"cloud-server/db"

	"github.com/google/uuid"
)

type UploadSession struct {
	ID       string    `json:"ID,omitempty"`
	UserID   string    `json:"userID,omitempty"`
	FileMeta *FileMeta `json:"fileMeta,omitempty"`
	file     *File

	ChunkSize uint     `json:"chunkSize,omitempty"`
	ChunkNum  uint64   `json:"chunkNum,omitempty"`
	ChunkMap  []uint64 `json:"chunkMap,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

type UploadChunk struct {
	ID        uint64 `json:"ID,omitempty"`
	SessionID string `json:"sessionID,omitempty"`

	Size uint64 `json:"size"`
	Data []byte `json:"data,omitempty"`
}

func CreateUploadSession(data *db.DB, userID string, duration time.Duration) (*UploadSession, error) {
	session := &UploadSession{
		ID:        uuid.New().String(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}

	file := &FileMeta{
		ID: uuid.New().String(),
	}
	session.FileMeta = file

	if _, err := data.Connection.Exec(db.Q_UPLOAD_INSERT, session.ID, userID, file.ID, session.ExpiresAt); err != nil {
		return nil, err
	}

	return session, nil
}

func FindUploadSession(data *db.DB, sessionID string) (*UploadSession, error) {
	session := &UploadSession{}
	err := data.Connection.QueryRow(db.Q_UPLOAD_FIND_BY_ID, sessionID).Scan(&session.UserID, &session.FileMeta.ID, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

//func remainingChunks(chunkMap []uint64) uint {
//	count := uint(0)
//	for _, chunk := range chunkMap {
//		count += uint(bits.OnesCount64(chunk))
//	}
//	return count
//}

func updateChunkMap(chunkMap []uint64, chunkID uint64) []uint64 {
	arrIndex := chunkID / 64
	bitIndex := chunkID % 64
	chunkMap[arrIndex] &^= 1 << bitIndex
	return chunkMap
}

func getNextChunk(chunkMap []uint64) (uint64, bool) {
	for i, chunk := range chunkMap {
		if chunk != 0 {
			bitPos := bits.TrailingZeros64(chunk)
			return uint64(uint(i)*64 + uint(bitPos)), true
		}
	}
	return 0, false
}

func (s *UploadSession) SaveChunk(data *db.DB, chunk *UploadChunk) (nextChunkID uint64, err error) {
	temp := tempDest + s.FileMeta.ID
	final := finalDest + s.FileMeta.ID

	if s.file == nil {
		s.file = &File{}
		s.file.io, err = os.OpenFile(temp, os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			log.Error("Failed to create file: ", err)
			return 0, err
		}

		// Pre-allocate the file
		if err = s.file.io.Truncate(int64(s.FileMeta.Size)); err != nil {
			log.Error("Failed to truncate file: ", err)
			return 0, err
		}
	}

	// Write the chunk
	offset := uint64(s.ChunkSize) * chunk.ID
	_, err = s.file.io.WriteAt(chunk.Data, int64(offset))
	if err != nil {
		log.Error("Failed to write chunk: ", err)
		return 0, err
	}

	s.ChunkMap = updateChunkMap(s.ChunkMap, chunk.ID)
	// Update the chunks
	s.ChunkNum -= 1

	if _, err := data.Connection.Exec(db.Q_UPLOAD_UPDATE_CHUNKS, s.ID, s.ChunkNum, chunk.ID); err != nil {
		log.Error("Failed to save chunk: ", err)
		return 0, err
	}

	if s.ChunkNum == 0 {
		if err := os.Rename(temp, final); err != nil {
			log.Error("Failed to rename file: ", err)
			return 0, err
		}

		s.file.io.Close()
		s.file = nil
		return 0, nil
	}

	if nextChunkID, ok := getNextChunk(s.ChunkMap); ok {
		return nextChunkID, nil
	}

	log.Error("Failed to get next chunk")
	return 0, nil
}
