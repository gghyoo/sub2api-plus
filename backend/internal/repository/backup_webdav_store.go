package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// WebDAVBackupStore implements service.BackupObjectStore using WebDAV protocol.
type WebDAVBackupStore struct {
	client   *http.Client
	baseURL  string
	username string
	password string
}

// NewWebDAVBackupStoreFactory returns a factory that creates WebDAV-backed stores.
func NewWebDAVBackupStoreFactory() func(ctx context.Context, cfg *service.BackupWebDAVConfig) (service.BackupObjectStore, error) {
	return func(_ context.Context, cfg *service.BackupWebDAVConfig) (service.BackupObjectStore, error) {
		if cfg.URL == "" {
			return nil, fmt.Errorf("webdav URL is required")
		}

		baseURL := strings.TrimRight(cfg.URL, "/")
		return &WebDAVBackupStore{
			client:   &http.Client{Timeout: 30 * time.Second},
			baseURL:  baseURL,
			username: cfg.Username,
			password: cfg.Password,
		}, nil
	}
}

func (w *WebDAVBackupStore) doReq(req *http.Request) (*http.Response, error) {
	if w.username != "" || w.password != "" {
		req.SetBasicAuth(w.username, w.password)
	}
	return w.client.Do(req)
}

func (w *WebDAVBackupStore) Upload(ctx context.Context, key string, body io.Reader, _ string) (int64, error) {
	// Ensure parent directories exist.
	if err := w.mkcolRecursive(ctx, key); err != nil {
		return 0, fmt.Errorf("webdav mkcol: %w", err)
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return 0, fmt.Errorf("read body: %w", err)
	}

	u, err := w.buildURL(key)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u, bytes.NewReader(data))
	if err != nil {
		return 0, fmt.Errorf("create put request: %w", err)
	}
	req.Header.Set("Content-Type", "application/gzip")

	resp, err := w.doReq(req)
	if err != nil {
		return 0, fmt.Errorf("webdav put: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("webdav put failed: %s (%.200s)", resp.Status, string(bodyBytes))
	}

	return int64(len(data)), nil
}

func (w *WebDAVBackupStore) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	u, err := w.buildURL(key)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create get request: %w", err)
	}

	resp, err := w.doReq(req)
	if err != nil {
		return nil, fmt.Errorf("webdav get: %w", err)
	}

	if resp.StatusCode >= 300 {
		_ = resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("webdav get failed: %s (%.200s)", resp.Status, string(bodyBytes))
	}

	return resp.Body, nil
}

func (w *WebDAVBackupStore) Delete(ctx context.Context, key string) error {
	u, err := w.buildURL(key)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return fmt.Errorf("create delete request: %w", err)
	}

	resp, err := w.doReq(req)
	if err != nil {
		return fmt.Errorf("webdav delete: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 && resp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webdav delete failed: %s (%.200s)", resp.Status, string(bodyBytes))
	}

	return nil
}

func (w *WebDAVBackupStore) PresignURL(_ context.Context, key string, _ time.Duration) (string, error) {
	// WebDAV does not support presigned URLs; return the direct URL.
	// The caller (browser) will need to provide Basic Auth credentials.
	u, err := w.buildURL(key)
	if err != nil {
		return "", err
	}
	return u, nil
}

func (w *WebDAVBackupStore) HeadBucket(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "PROPFIND", w.baseURL+"/", strings.NewReader(
		`<?xml version="1.0" encoding="utf-8"?>
<propfind xmlns="DAV:"><prop><resourcetype/></prop></propfind>`))
	if err != nil {
		return fmt.Errorf("create propfind request: %w", err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Depth", "0")

	resp, err := w.doReq(req)
	if err != nil {
		return fmt.Errorf("webdav propfind: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webdav connection test failed: %s (%.200s)", resp.Status, string(bodyBytes))
	}

	return nil
}

func (w *WebDAVBackupStore) buildURL(key string) (string, error) {
	u, err := url.Parse(w.baseURL)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}
	u.Path = path.Join(u.Path, key)
	return u.String(), nil
}

// mkcolRecursive creates intermediate directories for the given key.
func (w *WebDAVBackupStore) mkcolRecursive(ctx context.Context, key string) error {
	dir := path.Dir(key)
	if dir == "/" || dir == "." {
		return nil
	}

	parts := strings.Split(dir, "/")
	current := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		current = current + "/" + p
		u, err := w.buildURL(current)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, "MKCOL", u, nil)
		if err != nil {
			return fmt.Errorf("create mkcol request: %w", err)
		}

		resp, err := w.doReq(req)
		if err != nil {
			return fmt.Errorf("webdav mkcol %s: %w", current, err)
		}
		_ = resp.Body.Close()
		// 405 Method Not Allowed means the collection already exists — that's fine.
		if resp.StatusCode >= 300 && resp.StatusCode != http.StatusMethodNotAllowed {
			return fmt.Errorf("webdav mkcol %s failed: %s", current, resp.Status)
		}
	}

	return nil
}
