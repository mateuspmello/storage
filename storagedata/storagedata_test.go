package storagedata_test

import (
	"americanas/storagedata"
	"americanas/test"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

type fixture struct {
	sd *storagedata.StorageData
}

func setup() *fixture {

	sd := storagedata.New()
	return &fixture{
		sd: sd,
	}
}

func TestStorageFile(t *testing.T) {
	testCase := "TestStorageFile"

	f := setup()

	e := createFile("earth.png")

	req := map[string]interface{}{
		"path": "ht/monthly",
		"file": e,
		"name": "earth.png",
		"type": "png",
	}

	status, fileID, metadata, err := f.sd.StorageFile(req)

	md := metadata[string(fileID)].(map[string]interface{})
	modificationTime := md["modificationTime"]

	metadataExpected := map[string]interface{}{
		string(fileID): map[string]interface{}{
			"modificationTime": modificationTime,
			"name":             "earth.png",
			"path":             "ht/monthly/earth.png",
			"size":             int64(312866),
			"type":             "png",
		},
	}

	f.sd.DeleteByID(string(fileID))
	test.AssertEqual(t, testCase, metadata, metadataExpected)
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertNil(t, testCase, err)

}

func TestAllFiles(t *testing.T) {
	testCase := "TestAllFiles"
	f := setup()

	e := createFile("earth.png")

	req := map[string]interface{}{
		"path": "ht/monthly",
		"file": e,
		"name": "earth.png",
		"type": "png",
	}

	_, fileIDEarth, metadata, _ := f.sd.StorageFile(req)

	status, body, err := f.sd.AllFiles()
	var actual map[string]interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		fmt.Println(err)
	}

	md := metadata[string(fileIDEarth)].(map[string]interface{})
	metadataExpected := map[string]interface{}{
		string(fileIDEarth): map[string]interface{}{
			"modificationTime": md["modificationTime"],
			"name":             "earth.png",
			"path":             "ht/monthly/earth.png",
			"size":             float64(312866),
			"type":             "png",
		},
	}

	f.sd.DeleteByID(string(fileIDEarth))
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual, metadataExpected)
	test.AssertNil(t, testCase, err)

}

func TestUnderDir(t *testing.T) {
	testCase := "TestUnderDir"

	f := setup()

	e := createFile("earth.png")
	req := map[string]interface{}{
		"path": "space/planets",
		"file": e,
		"name": "earth.png",
		"type": "png",
	}
	_, fileIDEarth, metadataEarth, _ := f.sd.StorageFile(req)

	p := createFile("perseverance.png")
	req = map[string]interface{}{
		"path": "space/robots",
		"file": p,
		"name": "perseverance.png",
		"type": "png",
	}
	_, fileIDPerseverance, _, _ := f.sd.StorageFile(req)

	dir := "space/planets"
	status, body, err := f.sd.UnderDir(dir)
	var actual map[string]interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		fmt.Println(err)
	}

	md := metadataEarth[string(fileIDEarth)].(map[string]interface{})
	metadataExpected := map[string]interface{}{
		string(fileIDEarth): map[string]interface{}{
			"modificationTime": md["modificationTime"],
			"name":             "earth.png",
			"path":             "space/planets/earth.png",
			"size":             float64(312866),
			"type":             "png",
		},
	}

	f.sd.DeleteByID(string(fileIDEarth))
	f.sd.DeleteByID(string(fileIDPerseverance))
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual, metadataExpected)
	test.AssertNil(t, testCase, err)

}

func TestByID(t *testing.T) {
	testCase := "TestByID"

	f := setup()

	e := createFile("earth.png")
	req := map[string]interface{}{
		"path": "space/planets",
		"file": e,
		"name": "earth.png",
		"type": "png",
	}
	_, fileIDEarth, metadataEarth, _ := f.sd.StorageFile(req)

	status, body, err := f.sd.ByID(string(fileIDEarth))
	var actual map[string]interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		fmt.Println(err)
	}

	md := metadataEarth[string(fileIDEarth)].(map[string]interface{})
	metadataExpected := map[string]interface{}{
		"modificationTime": md["modificationTime"],
		"name":             "earth.png",
		"path":             "space/planets/earth.png",
		"size":             float64(312866),
		"type":             "png",
	}

	f.sd.DeleteByID(string(fileIDEarth))
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual, metadataExpected)
	test.AssertNil(t, testCase, err)

}

func TestMoveFile(t *testing.T) {
	testCase := "TestMoveFile"

	f := setup()

	e := createFile("earth.png")
	req := map[string]interface{}{
		"path": "space/planets",
		"file": e,
		"name": "earth.png",
		"type": "png",
	}
	_, fileIDEarth, _, _ := f.sd.StorageFile(req)

	id := string(fileIDEarth)
	toDir := "newproject/go"

	status, err := f.sd.MoveFile(id, toDir)

	m, errM := f.sd.GetMetadataJSON()
	if errM != nil {
		fmt.Println(err)
	}
	actualMetadata := m[id].(map[string]interface{})

	actualPath := actualMetadata["path"].(string)

	f.sd.DeleteByID(string(fileIDEarth))
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actualPath, `newproject/go/earth.png`)
	test.AssertNil(t, testCase, err)

}

func TestDeleteFile(t *testing.T) {
	testCase := "TestDeleteFile"

	f := setup()

	e := createFile("earth.png")
	req := map[string]interface{}{
		"path": "space/planets",
		"file": e,
		"name": "earth",
		"type": "png",
	}
	_, fileIDEarth, _, _ := f.sd.StorageFile(req)

	status, err := f.sd.DeleteByID(string(fileIDEarth))

	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertNil(t, testCase, err)

}

func TestOverwriteFile(t *testing.T) {
	testCase := "TestOverwriteFile"

	f := setup()

	e := createFile("earth.png")

	req := map[string]interface{}{
		"path": "ht/monthly",
		"file": e,
		"name": "earth.png",
		"type": "png",
	}

	_, fileID, meta, err := f.sd.StorageFile(req)
	fmt.Println(meta)
	if err != nil {
		fmt.Println(err)
	}

	m := createFile("mars.png")

	newMarsFile := map[string]interface{}{
		"path": "ht/monthly",
		"file": m,
		"name": "mars.png",
		"type": "png",
	}

	expected := map[string]interface{}{
		"name": "mars.png",
		"path": "ht/monthly/mars.png",
		"type": "png",
		"size": float64(338135),
	}

	status, body, err := f.sd.OverwriteFile(string(fileID), newMarsFile)
	var actual map[string]interface{}
	if err := json.Unmarshal(body, &actual); err != nil {
		fmt.Println(err)
	}

	f.sd.DeleteByID(string(fileID))
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual["name"], expected["name"])
	test.AssertEqual(t, testCase, actual["path"], expected["path"])
	test.AssertEqual(t, testCase, actual["type"], expected["type"])
	test.AssertEqual(t, testCase, actual["size"], expected["size"])
	test.AssertNil(t, testCase, err)

}

func createFile(fileName string) bytes.Buffer {
	path, err := filepath.Abs("../test_files/" + fileName)

	if err != nil {
		fmt.Println(err)
	}
	file, _ := os.Open(path)
	defer file.Close()

	var b bytes.Buffer
	io.Copy(&b, file)

	return b
}
