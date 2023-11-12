package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

// NewAuthorizedRequest creates a new HTTP request with the Authorization header set.
func NewAuthorizedRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("https://chat.openai.com/%s", endpoint), body)
	if err != nil {
		return nil, err
	}

	authToken := viper.GetString("OPENAI_BEARER_TOKEN")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Add("authority", "chat.openai.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US")
	req.Header.Add("content-type", "application/json")

	return req, nil
}

func Client() *http.Client {
	return http.DefaultClient
}

func checkDotEnv() bool {
	_, err := os.Stat(".env")
	return !os.IsNotExist(err)
}

func CheckBearerToken() (string, error) {
	viper.AutomaticEnv()
	if checkDotEnv() {
		viper.SetConfigFile(".env") // Specify the .env file
		err := viper.ReadInConfig() // Read the .env file
		if err != nil {
			log.Fatal("Error while reading config file", err)
		}
	}

	authToken := viper.GetString("OPENAI_BEARER_TOKEN")

	if authToken == "" {
		return "", fmt.Errorf("OPENAI_BEARER_TOKEN is not set")
	}

	return authToken, nil
}

func NewFileUploadRequest(url string, file *os.File) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", getFileContentType(file))
	req.Header.Set("origin", "https://chat.openai.com")
	req.Header.Set("authority", "files.oaiusercontent.com")
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US")

	return req, nil

}

func getFileContentType(file *os.File) string {
	// Get file extension
	ext := filepath.Ext(file.Name())
	// Set content type based on extension
	var contentType string
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg":
		contentType = "image/jpg"
	case ".jpeg":
		contentType = "image/jpeg"
	case ".pdf":
		contentType = "application/pdf"
	case ".txt":
		contentType = "text/plain"
	case ".csv":
		contentType = "text/csv"
	case ".json":
		contentType = "application/json"
	case ".xml":
		contentType = "application/xml"
	case ".html":
		contentType = "text/html"
	case ".md":
		contentType = "text/markdown"
	default:
		contentType = "application/octet-stream"
	}

	return contentType
}
