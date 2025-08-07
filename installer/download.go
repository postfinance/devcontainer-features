package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

type download struct{}

// Downloads the given url and returns the whole content as string.
func (d *download) AsString(url string) (string, error) {
	bytes, err := d.AsBytes(url)
	return string(bytes), err
}

// Downloads the given url and returns the whole content as byte slice.
func (d *download) AsBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file '%s'. Status code: %d", url, resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return bodyBytes, nil
}

// Downloads a file from the given url.
func (d *download) ToFile(url, filename string, progressName string) error {
	// Make sure to create the download folder
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check if the response status code is successful (200 OK)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file '%s'. Status code: %d", url, resp.StatusCode)
	}
	// Create the progress bar
	if progressName == "" {
		progressName = "Downloading"
	}
	progressBar := progressbar.DefaultBytes(
		resp.ContentLength,
		progressName,
	)
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// Copy the bytes
	_, err = io.Copy(io.MultiWriter(file, progressBar), resp.Body)
	if err != nil {
		return err
	}

	return nil
}
