package handler

import (
	"net/http"
)

type Res = http.ResponseWriter

// func (r *Res) WriteFile(path string) error {
// 	file, err := os.ReadFile(path)
// 	if err != nil {
// 		return err
// 	}
// 	r.Write(file)
// 	return nil
// }