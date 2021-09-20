package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type Api struct {
	storageDocument Storage
}

type Storage interface {
	StorageFile(body map[string]interface{}) (int, []byte, map[string]interface{}, error)
	AllFiles() (int, []byte, error)
	UnderDir(dir string) (int, []byte, error)
	ByID(id string) (int, []byte, error)
	MoveFile(id, toDir string) (int, error)
	DeleteByID(id string) (int, error)
	OverwriteFile(id string, body map[string]interface{}) (int, []byte, error)
}

var (
	IOMaxBufferSize = int64(16000000)
)

func (api *Api) RegisterRouters(router *httprouter.Router) {
	router.POST("/sendfile", api.sendFile)
	router.GET("/allfiles", api.allFiles)
	router.GET("/underdir", api.underDir)
	router.GET("/byid", api.byID)
	router.POST("/movefile", api.moveFile)
	router.POST("/delete", api.delete)
	router.POST("/overwrite", api.overwrite)
	router.ServeFiles("/storagedata/*filepath", http.Dir("/home/mateus-mello/go/src/americanas/storagedata"))

}

func (api *Api) sendFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	record, err := api.readBodyMultiPart(w, r)
	if err != nil {
		fmt.Println("[sendFile] Error in read body:", err.Error())
		errMap := map[string]interface{}{"error": "Invalid body"}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}
	statusCode, _, _, err := api.storageDocument.StorageFile(record)
	if statusCode != http.StatusOK {
		fmt.Printf("[sendFile] Error in sendFile with statusCode: %v - error %v", statusCode, err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}

	w.Header().Set("Location", "/sendfile")
	successMsg := map[string]interface{}{"status": "success"}
	api.send(w, statusCode, successMsg)
}

func (api *Api) allFiles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	statusCode, body, err := api.storageDocument.AllFiles()
	if statusCode != http.StatusOK {
		fmt.Printf("[allFiles] Error in allFiles with statusCode: %v - error %v", statusCode, err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}
	var responseMap map[string]interface{}
	_ = json.Unmarshal(body, &responseMap)

	w.Header().Set("Location", "/allfiles")
	api.send(w, statusCode, responseMap)
}

func (api *Api) underDir(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := api.getKeyFromURL(*r.URL)
	statusCode, body, err := api.storageDocument.UnderDir(url)
	if statusCode != http.StatusOK {
		fmt.Printf("[underDir] Error in underDir with statusCode: %v - error %v", statusCode, err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}

	w.Header().Set("Location", "/underdir?data="+url)
	api.send(w, statusCode, body)
}

func (api *Api) byID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := api.getKeyFromURL(*r.URL)
	statusCode, body, err := api.storageDocument.ByID(url)
	if statusCode != http.StatusOK {
		fmt.Printf("[byID] Error in byID with statusCode: %v - error %v", statusCode, err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}

	var responseMap map[string]interface{}
	_ = json.Unmarshal(body, &responseMap)

	w.Header().Set("Location", "/byid?data="+url)
	api.send(w, statusCode, responseMap)
}

func (api *Api) moveFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := api.getKeyFromURL(*r.URL)
	body, err := api.readBody(w, r.Body)
	if err != nil {
		fmt.Printf("[moveFile] Error in readBody. error %v", err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}
	toDir, ok := body["directory"].(string)
	if !ok {
		fmt.Printf("[moveFile] Error in body. error %v", err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}
	statusCode, err := api.storageDocument.MoveFile(id, toDir)
	if statusCode != http.StatusOK {
		fmt.Printf("[moveFile] Error in moveFile with statusCode: %v - error %v", statusCode, err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}
	w.Header().Set("Location", "/movefile?data="+id)
	successMsg := map[string]interface{}{"status": "success"}
	api.send(w, statusCode, successMsg)
}

func (api *Api) delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := api.getKeyFromURL(*r.URL)

	statusCode, err := api.storageDocument.DeleteByID(id)
	if statusCode != http.StatusOK {
		fmt.Printf("[delete] Error in delete with statusCode: %v - error %v", statusCode, err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}

	w.Header().Set("Location", "/delete?data="+id)
	successMsg := map[string]interface{}{"status": "success"}
	api.send(w, statusCode, successMsg)
}

func (api *Api) overwrite(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := api.getKeyFromURL(*r.URL)
	body, err := api.readBody(w, r.Body)
	if err != nil {
		fmt.Printf("[overwrite] Error in readBody. error %v", err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}

	statusCode, ret, err := api.storageDocument.OverwriteFile(id, body)
	if statusCode != http.StatusOK {
		fmt.Printf("[overwrite] Error in overwrite with statusCode: %v - error %v", statusCode, err.Error())
		errMap := map[string]interface{}{"error": err.Error()}
		api.send(w, http.StatusBadRequest, errMap)
		return
	}

	w.Header().Set("Location", "/movefile?data="+id)
	api.send(w, statusCode, ret)
}

func (api *Api) getKeyFromURL(url url.URL) string {
	keys, ok := url.Query()["data"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("[getKeyFromURL] Url Param 'data' is missing")
		return ""
	}
	return keys[0]
}

func (api *Api) readBodyMultiPart(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	// 10 << 20 specifies a maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	var buf bytes.Buffer
	io.Copy(&buf, file)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	record := make(map[string]interface{})

	record["name"] = handler.Filename
	record["size"] = handler.Size
	record["type"] = handler.Header.Get("Content-Type")
	record["file"] = buf
	record["path"] = r.FormValue("path")
	return record, err
}

func (api *Api) readBody(w http.ResponseWriter, body io.ReadCloser) (map[string]interface{}, error) {
	body = http.MaxBytesReader(w, body, IOMaxBufferSize)
	var record map[string]interface{}
	err := json.NewDecoder(body).Decode(&record)
	return record, err
}

func (api *Api) send(w http.ResponseWriter, statusCode int, value interface{}) {
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, TRACE, GET, HEAD, POST, PUT")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
	w.WriteHeader(statusCode)

	if value != nil {
		_ = json.NewEncoder(w).Encode(value)
	}
}

func New(storageDocument Storage) *Api {
	api := Api{
		storageDocument: storageDocument,
	}
	return &api
}
