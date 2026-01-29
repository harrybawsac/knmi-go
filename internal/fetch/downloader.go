// Package fetch provides HTTP download and zip extraction functionality.
package fetch

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Download fetches data from the given URL.
func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return data, nil
}

// ExtractZip extracts the first .txt file from a zip archive.
func ExtractZip(data []byte) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("opening zip: %w", err)
	}

	for _, f := range reader.File {
		if strings.HasSuffix(f.Name, ".txt") {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("opening file %s in zip: %w", f.Name, err)
			}
			defer rc.Close()

			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("reading file %s from zip: %w", f.Name, err)
			}

			return content, nil
		}
	}

	return nil, fmt.Errorf("no .txt file found in zip archive")
}
