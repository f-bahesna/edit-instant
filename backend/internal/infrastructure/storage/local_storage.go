package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// LocalStorage implements FileStorage using the local filesystem.
type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	os.MkdirAll(basePath, 0755)
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}
}

func (s *LocalStorage) Store(ctx context.Context, filename string, data []byte) (string, error) {
	fullPath := filepath.Join(s.basePath, filename)
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("local storage write failed: %w", err)
	}
	url := fmt.Sprintf("%s/%s", s.baseURL, filename)
	return url, nil
}

func (s *LocalStorage) GetURL(ctx context.Context, filename string) (string, error) {
	return fmt.Sprintf("%s/%s", s.baseURL, filename), nil
}
