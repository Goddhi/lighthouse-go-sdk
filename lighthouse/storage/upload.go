package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/httpx"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/cfg"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

type Service struct {
	h   *httpx.Client
	cfg cfg.Config
}

func New(h *httpx.Client, c cfg.Config) *Service { return &Service{h: h, cfg: c} }

// Upload

func (s *Service) UploadFile(ctx context.Context, path string, opts ...schema.UploadOption) (*schema.UploadResult, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()

	stat, _ := f.Stat()
	return s.UploadReader(ctx, filepath.Base(path), stat.Size(), f, opts...)
}

func (s *Service) UploadReader(ctx context.Context, name string, size int64, r io.Reader, opts ...schema.UploadOption) (*schema.UploadResult, error) {
	o := schema.DefaultUploadOptions()
	for _, opt := range opts { opt(o) }
	r = wrapWithProgress(r, size, nil) 

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer mw.Close()

		// file field
		fw, _ := mw.CreateFormFile("file", name)
		io.Copy(fw, r)
	}()
	url := s.cfg.Hosts.Upload + "/api/v0/add?cid-version=1"

	req, err := http.NewRequestWithContext(ctx, "POST", url, pr)
	if err != nil { return nil, err }
	req.Header.Set("Content-Type", mw.FormDataContentType())

	res, err := s.h.Inject(req)
	if err != nil { return nil, err }
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("http %d: %s", res.StatusCode, string(b))
	}

	var out schema.UploadResult
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
