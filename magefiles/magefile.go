//go:generate mage cacheLinodeOpenAPI

package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	linodeOpenAPIURL  = "https://raw.githubusercontent.com/linode/linode-api-openapi/refs/heads/main/openapi.json"
	linodeOpenAPIPath = ".cache/linode/openapi.json"
)

// CacheLinodeOpenAPI fetches the Linode OpenAPI spec with etag-based caching.
func CacheLinodeOpenAPI() error {
	etagPath := linodeOpenAPIPath + ".etag"

	if err := os.MkdirAll(filepath.Dir(linodeOpenAPIPath), 0o755); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, linodeOpenAPIURL, nil)
	if err != nil {
		return err
	}
	if etag, err := os.ReadFile(etagPath); err == nil && len(etag) > 0 {
		req.Header.Set("If-None-Match", strings.TrimSpace(string(etag)))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(in io.ReadCloser) {
		_ = in.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusNotModified {
		return nil
	}

	out, err := os.Create(linodeOpenAPIPath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	if etag := resp.Header.Get("ETag"); etag != "" {
		return os.WriteFile(etagPath, []byte(etag), 0o644)
	}
	return nil
}
