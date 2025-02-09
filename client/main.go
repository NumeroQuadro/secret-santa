package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	serverURL := "http://localhost:5555"

	filename := "goofyahhdocument.txt"
	err := uploadFile(serverURL, filename)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return
	}
	fmt.Println("File uploaded successfully.")

	downloadedFilename := "downloaded_" + filename
	err = downloadFile(serverURL, filename, downloadedFilename)
	if err != nil {
		fmt.Printf("Failed to download file: %v\n", err)
		return
	}
	fmt.Println("File downloaded successfully.")

	err = deleteFile(serverURL, filename)
	if err != nil {
		fmt.Printf("Failed to delete file: %v\n", err)
		return
	}
	fmt.Println("File deleted successfully.")
}

func uploadFile(serverURL, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	resp, err := http.Post(serverURL+"/put?filename="+filename, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	return nil
}

func downloadFile(serverURL, filename, destFilename string) error {
	resp, err := http.Get(serverURL + "/get?filename=" + filename)
	if err != nil {
		return fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	destFile, err := os.Create(destFilename)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to destination file: %w", err)
	}

	return nil
}

func deleteFile(serverURL, filename string) error {
	req, err := http.NewRequest(http.MethodDelete, serverURL+"/delete?filename="+filename, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %s", resp.Status)
	}

	return nil
}
