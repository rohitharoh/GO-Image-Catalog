package api

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"path/filepath"
	"strconv"

	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/tb/image-catalog/imageapp/services"

	"fmt"
	"github.com/tb/image-catalog/common-packages/system"
	"io"
	"os"
)

type Controller struct {
	system.Controller
}

func (controller *Controller) ListFilesInS3(c web.C, w http.ResponseWriter, r *http.Request, logger *log.Entry) ([]byte, error) {
	skip := r.URL.Query().Get("skip")
	skipValue, _ := strconv.Atoi(skip)

	decoder := json.NewDecoder(r.Body)
	var data map[string]string
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("error -->", err)
		return nil, err
	}
	filePath := data["filePath"]

	refId := c.Env["refId"]
		
	response, err := services.ListFilesInS3(refId.(string), filePath, skipValue)
	if err != nil {

		return nil, err
	}
	return response, nil
}

func (controller *Controller) UploadMultipleFiles(c web.C, w http.ResponseWriter, r *http.Request, logger *log.Entry) ([]byte, error) {

	_ = r.ParseMultipartForm(32 << 20)
	folderPath := r.URL.Query().Get("folderPath")
	refId := c.Env["refId"]

	m := r.MultipartForm
	files := m.File["fileUpload"]
	finalResponse := []map[string]interface{}{}
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		fileName := files[i].Filename
		ext := filepath.Ext(fileName)


		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
		file, _ := files[i].Open()
		defer file.Close()

		uniqueFileName := fileName

		filePath := system.UPLOAD_FOLDER + uniqueFileName

		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			fmt.Println("err opening file")
			fmt.Println("error -->", err)
			return nil, err
		}

		_, _ = io.Copy(f, file)

		response, err := services.UploadFileMedia(f, filePath, fileName, folderPath, refId.(string))

		if err != nil {
			fmt.Println("error --->", err)
			return nil, err
		}

		finalResponse = append(finalResponse, response)
	} else {
			return nil, system.NotAnImageFile
	}
	}
	finalResponse11, _ := json.Marshal(finalResponse)
	return finalResponse11, nil
}
