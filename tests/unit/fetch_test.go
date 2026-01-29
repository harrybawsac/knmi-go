package unit

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/harrybawsac/knmi-go/internal/fetch"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errContains string
	}{
		{
			name: "successful download",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test content"))
			},
			wantErr: false,
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:     true,
			errContains: "500",
		},
		{
			name: "not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr:     true,
			errContains: "404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			data, err := fetch.Download(server.URL)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tt.errContains)) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(data) == 0 {
				t.Error("expected non-empty data")
			}
		})
	}
}

func TestExtractZip(t *testing.T) {
	tests := []struct {
		name        string
		setupZip    func() []byte
		wantErr     bool
		errContains string
	}{
		{
			name: "valid zip with txt file",
			setupZip: func() []byte {
				return createTestZip(t, map[string]string{
					"etmgeg_260.txt": "# KNMI data\ntest content",
				})
			},
			wantErr: false,
		},
		{
			name: "empty zip",
			setupZip: func() []byte {
				return createTestZip(t, map[string]string{})
			},
			wantErr:     true,
			errContains: "no .txt file",
		},
		{
			name: "invalid zip data",
			setupZip: func() []byte {
				return []byte("not a zip file")
			},
			wantErr:     true,
			errContains: "zip",
		},
		{
			name: "multiple txt files picks first",
			setupZip: func() []byte {
				return createTestZip(t, map[string]string{
					"file1.txt": "content 1",
					"file2.txt": "content 2",
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipData := tt.setupZip()

			content, err := fetch.ExtractZip(zipData)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tt.errContains)) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(content) == 0 {
				t.Error("expected non-empty content")
			}
		})
	}
}

// createTestZip creates a zip archive with the given files.
func createTestZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	dir := t.TempDir()
	zipPath := filepath.Join(dir, "test.zip")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatalf("failed to create zip file: %v", err)
	}

	w := zip.NewWriter(f)
	for name, content := range files {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatalf("failed to create file in zip: %v", err)
		}
		if _, err := fw.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write to zip: %v", err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	data, err := os.ReadFile(zipPath)
	if err != nil {
		t.Fatalf("failed to read zip file: %v", err)
	}

	return data
}
