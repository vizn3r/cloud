package fs

import (
	"fmt"
	"math/rand/v2"
	"os"
)

func FindFile(path string) []byte {
	data, err := os.ReadFile("storage/" + path)
	if err != nil {
		return nil
	}
	return data
}

func SaveFile(data []byte) int {
	id := rand.IntN(999999999)
	os.WriteFile("storage/"+fmt.Sprintf("%d", id), data, 0644)
	return id
}
