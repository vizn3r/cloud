package fs

import (
	"cloud-server/db"
	"time"

	"github.com/google/uuid"
)

type Share struct {
	UUID      string
	FileUUID  string
	CreatedAt time.Time
	ExpiresAt time.Time
	Downloads uint64

	db *db.DB
}

func NewShare(id string, db *db.DB) *Share {
	return &Share{
		UUID: id,
		db:   db,
	}
}

func (sh *Share) FindShare() (File, error) {
	err := sh.db.Connection.QueryRow(db.Q_SHARE_FIND_BY_ID, sh.UUID).Scan(sh.FileUUID, sh.Downloads)
	if err != nil {
		return File{}, err
	}

	return FindFile(sh.FileUUID)
}

func (sh *Share) SaveShare() (Share, error) {
	id := uuid.New().String()
	sh.UUID = id
	err := sh.db.Connection.QueryRow(db.Q_SHARE_INSERT, sh.UUID).Scan(sh.UUID, sh.FileUUID)
	if err != nil {
		return Share{}, err
	}

	return *sh, err
}
