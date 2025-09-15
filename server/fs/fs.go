package fs

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
)

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

func FindFile(path string) (File, error) {
	data, err := os.ReadFile("storage/" + path)
	if err != nil {
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
		return file, err
	}

	return file, nil
}

func (file File) SaveFile() (string, error) {
	id := uuid.New().String()
	file.Meta.FileUUID = id

	metaJSON, err := json.Marshal(file.Meta)
	if err != nil {
		return "", err
	}

	var comb bytes.Buffer
	comb.Write(metaJSON)
	comb.WriteString(META_SEPARATOR)
	comb.Write(file.Data)

	if err := os.MkdirAll("storage/temp", 0755); err != nil {
		return "", err
	}

	temp := "storage/temp/" + id
	final := "storage/" + id

	err = os.WriteFile(temp, comb.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	if err := os.Rename(temp, final); err != nil {
		os.Remove(temp)
		return "", err
	}

	return id, nil
}
