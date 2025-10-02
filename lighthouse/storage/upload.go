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
	"sync/atomic"


	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/httpx"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/cfg"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

type Service struct {
	h   *httpx.Client
	cfg cfg.Config
}

func New(h *httpx.Client, c cfg.Config) *Service { 
	return &Service{h: h, cfg: c} 
}

type progressReader struct {
	r        io.Reader
	uploaded *int64 // atomic counter
	total    int64
	onProg   schema.ProgressCallback
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	if n > 0 {
		uploaded := atomic.AddInt64(p.uploaded, int64(n))
		if p.onProg != nil {
			p.onProg(schema.Progress{
				Uploaded: uploaded,
				Total:    p.total,
			})
		}
	}
	return n, err
}

func (s *Service) UploadFile(ctx context.Context, path string, opts ...schema.UploadOption) (*schema.UploadResult, error) {
	f, err := os.Open(path)
	if err != nil {
		 return nil, err 
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return s.UploadReader(ctx, filepath.Base(path), stat.Size(), f, opts...)
}

func (s *Service) UploadReader(ctx context.Context, name string, size int64, r io.Reader, opts ...schema.UploadOption) (*schema.UploadResult, error) {
	o := schema.DefaultUploadOptions()
	for _, opt := range opts { 
		opt(o) 
	}

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	boundary := mw.Boundary()

	headerSize := int64(len(fmt.Sprintf("--%s\r\nContent-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\nContent-Type: application/octet-stream\r\n\r\n", boundary, name)))
	footerSize := int64(len(fmt.Sprintf("\r\n--%s--\r\n", boundary)))
	totalSize := headerSize + size + footerSize


	// Goroutine to write multipart data
	go func() {
		defer pw.Close()
		defer mw.Close()

		// file field
		fw, err := mw.CreateFormFile("file", name)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		buf := make([]byte, 32*1024)
		for {
			select {
			case <-ctx.Done():
				pw.CloseWithError(ctx.Err())
				return
			default:
			}

			n, err := r.Read(buf)
			if n > 0 {
				if _, writeErr := fw.Write(buf[:n]); writeErr != nil {
					pw.CloseWithError(writeErr)
					return
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	var uploaded int64
	progressPipeReader := io.Reader(pr)  // Type as io.Reader, not *io.PipeReader
	if o.OnProgress != nil {
		progressPipeReader = &progressReader{
			r:        pr,
			uploaded: &uploaded,
			total:    totalSize,
			onProg:   o.OnProgress,
		}
	}
	url := s.cfg.Hosts.Upload + "/api/v0/add?cid-version=1"
	req, err := http.NewRequestWithContext(ctx, "POST", url, progressPipeReader)
	if err != nil {
		 return nil, err 
		}
	req.Header.Set("Content-Type", mw.FormDataContentType())

	res, err := s.h.Inject(req)
	if err != nil { 
		return nil, err 
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("http %d: %s", res.StatusCode, string(b))
	}

	var result schema.UploadResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	if o.OnProgress != nil {
		o.OnProgress(schema.Progress{
			Uploaded: totalSize,
			Total:    totalSize,
		})
	}
	return &result, nil
}
