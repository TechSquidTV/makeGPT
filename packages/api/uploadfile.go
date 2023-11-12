package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/TechSquidTV/makeGPT/packages/utils"
	"github.com/charmbracelet/log"
)

// Valid types: "profile_picture", "gizmo"
type UploadFilePurpose string

// // This currently reports as unused, which may be because it is not a true type.
// // However, this appears working in the GPTTool type.
// const (
// 	profile_picture UploadFilePurpose = "profile_picture"
// 	gizmo           UploadFilePurpose = "gizmo"
// )

type FileProperties struct {
	Filename string `json:"file_name"`
	Filesize int64  `json:"file_size"`
}

type UploadFilePayload struct {
	Usecase UploadFilePurpose `json:"use_case"`
	FileProperties
}

type UploadFileResponse struct {
	// "success"
	Status    string `json:"status"`
	UploadURL string `json:"upload_url"`
	FileID    string `json:"file_id"`
}

func UploadFile(path string, purpose UploadFilePurpose) error {
	//

	file := FileProperties{
		Filename: filepath.Base(path),
		Filesize: utils.GetFileSize(path),
	}
	payload := UploadFilePayload{
		Usecase:        UploadFilePurpose(purpose),
		FileProperties: file,
	}

	log.Debugf("Uploading file %s (%d bytes) for %s", file.Filename, file.Filesize, purpose)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal upload payload: \n%v\n", err)
		return err
	}
	req, err := NewAuthorizedRequest("POST", "backend-api/files", strings.NewReader(string(payloadJSON)))
	if err != nil {
		log.Warnf("Request:  \n%v\n", req)
		log.Fatalf("Failed to create request for file: \n%v\n", err)
		return err
	}
	var response UploadFileResponse
	resp, err := Client().Do(req)
	if err != nil {
		log.Warnf("Response: \n%v\n", resp)
		log.Fatalf("Failed to get upload URL: \n%v\n", err)
		return err
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Warnf("Request Headers: \n%v\n", req.Header)
		log.Warnf("Request Body: \n%v\n", payload)
		log.Warnf("Response: \n%v\n", response)
		log.Fatal("Failed to decode upload response", err)
		return err
	}
	resp.Body.Close()

	log.Infof("Uploading file to:  %s", response.UploadURL)
	req = nil
	err = nil
	fileData, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error opening file for upload: \n%v\n", err)
		return err
	}

	req, err = NewFileUploadRequest(response.UploadURL, fileData)
	if err != nil {
		log.Fatalf("Error creating file upload request: \n%v\n", err)
		return err
	}

	fresp, err := Client().Do(req)
	if err != nil {
		log.Warnf("Upload response: \n%v\n", resp)
		log.Fatalf("Failed to upload file: \n%v\n", err)
		return err
	}
	defer fresp.Body.Close()
	if fresp.StatusCode != 201 {
		log.Warnf("Upload response: \n%v\n", fresp)
		log.Fatalf("Failed to upload file: \n%v\n", err)
		return err
	}
	log.Infof("Uploaded file %s to %s with ID %s", file.Filename, response.UploadURL, response.FileID)
	return nil
}
