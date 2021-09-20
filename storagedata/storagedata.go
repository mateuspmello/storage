package storagedata

import (
	"americanas/helper"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StorageData struct {
}

func (s *StorageData) StorageFile(body map[string]interface{}) (int, []byte, map[string]interface{}, error) {

	file, fullPath, typeFile, err := s.saveFileInDisk(body)
	if err != nil {
		return http.StatusBadRequest, nil, nil, err
	}
	defer file.Close()

	metadata, err := file.Stat()
	if err != nil {
		return http.StatusBadRequest, nil, nil, err
	}

	hash := md5.Sum([]byte(metadata.ModTime().String() + metadata.Name()))
	fileID := hex.EncodeToString(hash[:])
	modTime := metadata.ModTime()
	hh := modTime.Format("01/02/2006 15:04:05")
	metadataJSON := map[string]interface{}{
		fileID: map[string]interface{}{
			"name":             metadata.Name(),
			"path":             fullPath,
			"type":             typeFile,
			"size":             metadata.Size(),
			"modificationTime": hh,
		},
	}

	metadataFile, err := s.openMetadataJSON(os.O_CREATE | os.O_WRONLY)
	if err != nil {
		return http.StatusBadRequest, nil, nil, err
	}
	defer metadataFile.Close()

	fileMap := make(map[string]interface{})
	mtpath, _ := filepath.Abs("metadata.json")
	byteValue, _ := ioutil.ReadFile(mtpath)
	err = json.Unmarshal(byteValue, &fileMap)
	if err != nil {
		return http.StatusBadRequest, nil, nil, err
	}
	fileMap = helper.MergeMaps(fileMap, metadataJSON)
	fileOut, _ := json.MarshalIndent(fileMap, "", "	")
	_, err = io.Copy(metadataFile, bytes.NewReader(fileOut))
	if err != nil {
		return http.StatusBadRequest, nil, nil, err
	}

	return http.StatusOK, []byte(fileID), metadataJSON, nil
}

func (s *StorageData) AllFiles() (int, []byte, error) {

	mapFileMetadata, err := s.GetMetadataJSON()
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	ret, err := json.Marshal(mapFileMetadata)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, ret, nil
}

func (s *StorageData) UnderDir(dir string) (int, []byte, error) {

	mapFileMetadata, err := s.GetMetadataJSON()

	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	mapreturn := make(map[string]interface{})
	for k, v := range mapFileMetadata {
		item, ok := v.(map[string]interface{})
		if !ok {
			return http.StatusBadRequest, nil, err
		}
		d, ok := item["path"].(string)
		if !ok {
			return http.StatusBadRequest, nil, err
		}

		if strings.HasPrefix(d, dir) {
			m := make(map[string]interface{})
			m[k] = v
			mapreturn = helper.MergeMaps(mapreturn, m)
		}

	}

	ret, err := json.Marshal(mapreturn)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, ret, nil
}

func (s *StorageData) ByID(id string) (int, []byte, error) {

	mapFileMetadata, err := s.GetMetadataJSON()
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	mapWithId, ok := mapFileMetadata[id].(map[string]interface{})
	if !ok {
		return http.StatusBadRequest, nil, err
	}

	ret, err := json.Marshal(mapWithId)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, ret, nil
}

func (s *StorageData) MoveFile(id, toDir string) (int, error) {

	mapFileMetadata, err := s.GetMetadataJSON()
	if err != nil {
		return http.StatusBadRequest, err
	}

	mapWithId, ok := mapFileMetadata[id].(map[string]interface{})
	if !ok {
		return http.StatusBadRequest, err
	}

	path, ok := mapWithId["path"].(string)
	if !ok {
		return http.StatusBadRequest, err
	}

	nameFile, ok := mapWithId["name"].(string)
	if !ok {
		return http.StatusBadRequest, err
	}

	toDirComplete := ""
	if toDirComplete, err = filepath.Abs(toDir); err != nil {
		return http.StatusBadRequest, err
	}

	toDirAndFile := toDirComplete + "/" + nameFile

	fromPath := ""
	if fromPath, err = filepath.Abs(path); err != nil {
		return http.StatusBadRequest, err
	}

	os.MkdirAll(toDirComplete, os.ModePerm)
	err = os.Rename(fromPath, toDirAndFile)
	if err != nil {
		return http.StatusBadRequest, err
	}

	mapWithId["path"] = toDir + "/" + nameFile
	// aa := make(map[string]interface{})
	// aa[id] = mapWithId
	mapFileMetadata[id] = mapWithId
	err = s.WriteMetadataInDisk(mapFileMetadata)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func (s *StorageData) DeleteByID(id string) (int, error) {

	mapFileMetadata, err := s.GetMetadataJSON()
	if err != nil {
		return http.StatusBadRequest, err
	}

	mapWithId, ok := mapFileMetadata[id].(map[string]interface{})
	if !ok {
		return http.StatusBadRequest, err
	}

	path := ""
	if path, ok = mapWithId["path"].(string); !ok {
		return http.StatusBadRequest, err
	}

	pathComplete := ""
	if pathComplete, err = filepath.Abs(path); err != nil {
		return http.StatusBadRequest, err
	}

	err = os.Remove(pathComplete)
	if err != nil {
		return http.StatusBadRequest, err
	}

	delete(mapFileMetadata, id)
	err = s.WriteMetadataInDisk(mapFileMetadata)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func (s *StorageData) OverwriteFile(id string, body map[string]interface{}) (int, []byte, error) {

	_, err := s.DeleteByID(id)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	file, fullPath, typeFile, err := s.saveFileInDisk(body)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	defer file.Close()

	mapFileMetadata, err := s.GetMetadataJSON()
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	metadata, err := file.Stat()
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	dataToOverWrite := map[string]interface{}{
		"name":             metadata.Name(),
		"path":             fullPath,
		"type":             typeFile,
		"size":             metadata.Size(),
		"modificationTime": metadata.ModTime(),
	}

	mapFileMetadata[id] = dataToOverWrite

	err = s.WriteMetadataInDisk(mapFileMetadata)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	ret, err := json.Marshal(dataToOverWrite)
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, ret, nil

}

func (s *StorageData) openMetadataJSON(flags int) (*os.File, error) {
	mtpath, err := filepath.Abs("metadata.json")
	if err != nil {
		return nil, err
	}

	metadataFile, err := os.OpenFile(mtpath, flags, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return metadataFile, nil
}

func (s *StorageData) saveFileInDisk(body map[string]interface{}) (*os.File, string, string, error) {

	err := validateBody(body)
	if err != nil {
		fmt.Printf("[saveFileInDisk] Wrong body format. Error: %s", err)
		return nil, "", "", err
	}

	f := body["file"].(bytes.Buffer)
	path := body["path"].(string)
	name := body["name"].(string)
	typeFile := body["type"].(string)

	_ = os.MkdirAll(path, os.ModePerm)
	fullPath := fmt.Sprintf("%s/%s.%s", path, name, typeFile)
	_, fileExists := os.Stat(fullPath)

	count := 1
	for fileExists == nil {
		fullPath = fmt.Sprintf("%s/%s(%v).%s", path, name, count, typeFile)
		count++
		_, fileExists = os.Stat(fullPath)
	}

	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, "", "", err
	}
	// Copy the file to the destination path
	_, err = io.Copy(file, bytes.NewReader(f.Bytes()))
	if err != nil {
		return nil, "", "", err
	}

	return file, fullPath, typeFile, err
}

func (s *StorageData) WriteMetadataInDisk(newMetadata map[string]interface{}) error {

	fileMetadata, err := s.openMetadataJSON(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)
	if err != nil {
		return err
	}
	defer fileMetadata.Close()

	mapMetadataIdent, err := json.MarshalIndent(newMetadata, "", "	")
	if err != nil {
		return err
	}

	_, err = io.Copy(fileMetadata, bytes.NewReader(mapMetadataIdent))
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageData) GetMetadataJSON() (map[string]interface{}, error) {
	mtpath, err := filepath.Abs("metadata.json")
	if err != nil {
		return nil, err
	}

	fileMetadata, err := ioutil.ReadFile(mtpath)
	if err != nil {
		return nil, err
	}

	mapFileMetadata := make(map[string]interface{})
	err = json.Unmarshal(fileMetadata, &mapFileMetadata)
	if err != nil {
		return nil, err
	}
	return mapFileMetadata, nil
}

func validateBody(body map[string]interface{}) error {
	var fieldsMissing []string

	value, ok := body["path"].(string)
	if !ok || value == "" {
		fieldsMissing = append(fieldsMissing, "path")
	}

	_, aa := body["file"].(bytes.Buffer)
	if !aa {
		fieldsMissing = append(fieldsMissing, "file")
	}

	if len(fieldsMissing) > 0 {
		return errors.New("missing fields: " + strings.Join(fieldsMissing, ", "))
	}

	return nil
}

func New() *StorageData {

	sd := StorageData{}

	return &sd
}
