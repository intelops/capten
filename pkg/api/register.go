package api

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"capten/pkg/config"
)

type fileInfo struct {
	fileKey string
	path    string
}

func RegisterAgentInfo(customerId, agentHost string) {
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Errorf("failed to load the config %v", err)
		return
	}

	filesToUpload := []fileInfo{
		{
			fileKey: "ca_crt",
			path:    "./cert/ca.crt",
		},
		{
			fileKey: "client_crt",
			path:    "./cert/client.crt",
		},
		{
			fileKey: "client_key",
			path:    "./cert/client.key",
		},
	}

	headers := map[string]string{
		"customer_id": customerId,
		"endpoint":    agentHost,
	}

	url := fmt.Sprintf("%s://%s/register/agent",
		cfg.GetString("saas.scheme"),
		cfg.GetString("saas.endpoint"))

	if err := UploadFiles(filesToUpload, url, headers); err != nil {
		logrus.Error("failed to register agent", err)
	}
}

func UploadFiles(files []fileInfo, url string, headers map[string]string) error {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	for _, f := range files {
		part, err := w.CreateFormFile(f.fileKey, filepath.Base(f.path))
		if err != nil {
			return err
		}

		fileContents, err := os.ReadFile(f.path)
		if err != nil {
			return err
		}

		_, err = part.Write(fileContents)
		if err != nil {
			return err
		}
	}

	err := w.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", w.FormDataContentType())
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register status code %v", res.StatusCode)
	}

	fmt.Println(string(response))
	return nil
}
