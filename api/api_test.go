package api_test

import (
	"americanas/api"
	"americanas/test"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

type fixture struct {
	router  *httprouter.Router
	api     *api.Api
	storage *StorageFake
}

type StorageFake struct {
	status int
	err    error
	body   []byte
}

func (s *StorageFake) StorageFile(body map[string]interface{}) (int, []byte, map[string]interface{}, error) {
	return s.status, s.body, nil, s.err
}

func (s *StorageFake) AllFiles() (int, []byte, error) {
	return s.status, s.body, s.err
}

func (s *StorageFake) UnderDir(dir string) (int, []byte, error) {
	return s.status, s.body, s.err
}

func (s *StorageFake) ByID(id string) (int, []byte, error) {
	return s.status, s.body, s.err
}

func (s *StorageFake) MoveFile(id, toDir string) (int, error) {
	return s.status, s.err
}

func (s *StorageFake) DeleteByID(id string) (int, error) {
	return s.status, s.err
}

func (s *StorageFake) OverwriteFile(id string, body map[string]interface{}) (int, []byte, error) {
	return s.status, s.body, s.err
}

func setup(t *testing.T) *fixture {
	s := &StorageFake{}
	router := httprouter.New()
	api := api.New(s)
	api.RegisterRouters(router)
	return &fixture{
		api:     api,
		storage: s,
		router:  router,
	}
}

func TestPOSTSendFile(t *testing.T) {
	testCase := "test-post-send-file-with-sucess"
	url := "/sendfile"
	absPath, _ := filepath.Abs("test_files/mars.png")
	file, _ := os.Open(absPath)
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", absPath)
	_, _ = io.Copy(part, file)
	writer.WriteField("path", "ht/monthly")
	writer.Close()

	fixture := setup(t)
	fixture.storage.status = http.StatusOK

	status, returnBody, header := fixture.requestMultiPart(url, "POST", body, *writer)
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, returnBody, `{"status":"success"}`)
	test.AssertEqual(t, testCase, header.Get("Location"), "/sendfile")

}

func TestGETGetFile(t *testing.T) {
	testCase := "test-get-get-file-with-sucess"
	url := "storagedata/test/mars.png"
	fixture := setup(t)
	fixture.storage.status = http.StatusOK

	reader := bytes.NewReader(getFileTest("mars.png"))
	buf := bufio.NewReader(reader)
	line, _ := buf.ReadBytes('\n')
	expected := strings.TrimSuffix(string(line), "\n")

	status, actual, _ := fixture.request(url, "GET", nil)
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual, expected)

}

func TestGETAllFiles(t *testing.T) {
	testCase := "test-get-all-files-with-sucess"
	url := "/allfiles"
	fixture := setup(t)
	fixture.storage.status = http.StatusOK
	fixture.storage.body = fakeBodyResult()

	status, body, header := fixture.request(url, "GET", nil)

	var actual map[string]interface{}
	json.Unmarshal([]byte(body), &actual)

	var expected map[string]interface{}
	json.Unmarshal(fakeBodyResult(), &expected)

	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual, expected)
	test.AssertEqual(t, testCase, header.Get("Location"), "/allfiles")

}

func TestGETUnderDir(t *testing.T) {
	testCase := "test-get-under-dir-list-with-sucess"
	url := "/underdir?data=%s"
	fixture := setup(t)
	fixture.storage.status = http.StatusOK
	fixture.storage.body = underDirFakeResult()

	dir := "ht/monthly"
	url = fmt.Sprintf(url, dir)
	status, body, header := fixture.request(url, "GET", nil)

	var actual map[string]interface{}
	json.Unmarshal([]byte(body), &actual)

	var expected map[string]interface{}
	json.Unmarshal(underDirFakeResult(), &expected)

	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual, expected)
	test.AssertEqual(t, testCase, header.Get("Location"), url)

}

func TestGETByID(t *testing.T) {
	testCase := "test-get-by-id-with-sucess"
	url := "/byid?data=%s"
	fixture := setup(t)
	fixture.storage.status = http.StatusOK
	fixture.storage.body = byIDFakeResult()

	id := "aab053840116dacaf13a062d909e5761"
	url = fmt.Sprintf(url, id)
	status, body, header := fixture.request(url, "GET", nil)

	var actual map[string]interface{}
	json.Unmarshal([]byte(body), &actual)

	var expected map[string]interface{}
	json.Unmarshal(byIDFakeResult(), &expected)

	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, actual, expected)
	test.AssertEqual(t, testCase, header.Get("Location"), url)

}

func TestPOSTMoveFile(t *testing.T) {
	testCase := "test-post-send-file-with-sucess"
	url := "/movefile?data=%s"

	fixture := setup(t)
	fixture.storage.status = http.StatusOK
	fixture.storage.body = directory()

	id := "aab053840116dacaf13a062d909e5761"
	url = fmt.Sprintf(url, id)
	status, returnBody, header := fixture.request(url, "POST", bytes.NewBuffer([]byte(directory())))
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, returnBody, `{"status":"success"}`)
	test.AssertEqual(t, testCase, header.Get("Location"), url)

}

func TestPOSTDeleteFile(t *testing.T) {
	testCase := "test-post-delete-file-with-sucess"
	url := "/delete?data=%s"

	fixture := setup(t)
	fixture.storage.status = http.StatusOK

	id := "aab053840116dacaf13a062d909e5761"
	url = fmt.Sprintf(url, id)
	status, returnBody, header := fixture.request(url, "POST", bytes.NewBuffer([]byte(id)))
	test.AssertEqual(t, testCase, status, http.StatusOK)
	test.AssertEqual(t, testCase, returnBody, `{"status":"success"}`)
	test.AssertEqual(t, testCase, header.Get("Location"), url)

}

func (f *fixture) createRequest(url string, method string, body io.Reader) *http.Request {
	server := httptest.NewServer(f.router)
	url = fmt.Sprintf("%v/%v", server.URL, url)

	req, _ := http.NewRequest(method, url, body)

	return req
}

func (f *fixture) sendRequest(req *http.Request) (int, string, http.Header) {
	client := &http.Client{}
	resp, _ := client.Do(req)

	reader := bufio.NewReader(resp.Body)
	line, _ := reader.ReadBytes('\n')
	resp.Body.Close()

	return resp.StatusCode, strings.TrimSuffix(string(line), "\n"), resp.Header
}

func (f *fixture) requestMultiPart(url string, method string, body io.Reader, writer multipart.Writer) (int, string, http.Header) {

	req := f.createRequest(url, method, body)
	fmt.Println(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return f.sendRequest(req)

}

func (f *fixture) request(url string, method string, body io.Reader) (int, string, http.Header) {

	req := f.createRequest(url, method, body)
	return f.sendRequest(req)
}

func getFileTest(filename string) []byte {
	path, err := filepath.Abs("../test_files/" + filename)
	if err != nil {
		fmt.Println(err)
	}
	file, _ := os.Open(path)
	defer file.Close()

	var b bytes.Buffer
	io.Copy(&b, file)

	return b.Bytes()

}

func underDirFakeResult() []byte {
	return []byte(`{
		"aab053840116dacaf13a062d909e5761": {
			"modificationTime": "2021-09-11T23:07:39.696063108-03:00",
			"name": "golang.png",
			"path": "ht/monthly/golang.png",
			"size": 24357,
			"type": "png"
		},
	}`)
}

func byIDFakeResult() []byte {
	return []byte(`{
		"modificationTime": "2021-09-11T23:07:39.696063108-03:00",
		"name": "golang.png",
		"path": "ht/monthly/golang.png",
		"size": 24357,
		"type": "png"
	}`)
}

func directory() []byte {
	return []byte(`{		
		"directory": "ht/monthly"
	}`)
}

func fakeBodyResult() []byte {

	return []byte(`{
		"ht/mar.png": {
			"name": "mar.png",
			"path": "ht/",								
			"type": "png",
			"dateCreated": "",
			"dateModify": "2021-09-09T08:20:27Z",
			"size": "24357",
		},
		"ht/mar.png": {
			"name": "mar.png",
			"path": "",								
			"type": "",
			"dateCreated": "",
			"dateModify": "",
			"size": "",
		},
		"ht/mar.png": {
			"name": "mar.png",
			"path": "",								
			"type": "",
			"dateCreated": "",
			"dateModify": "",
			"size": "",
		},
		"ht/mar.png": {
			"name": "mar.png",
			"path": "",								
			"type": "",
			"dateCreated": "",
			"dateModify": "",
			"size": "",
		},
	}`)

}
