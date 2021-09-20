package helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func CreateFromFile(path string) io.Writer {
	absPath, _ := filepath.Abs(path)
	file, _ := os.Open(absPath)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", absPath)
	return part
}

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
