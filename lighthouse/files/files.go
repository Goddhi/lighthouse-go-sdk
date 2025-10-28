package files

import (
	"context"
	"net/url"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/cfg"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/httpx"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

type Service struct {
	h   *httpx.Client
	cfg cfg.Config
}

func New(h *httpx.Client, c cfg.Config) *Service {
	return &Service{h: h, cfg: c}
}

func (s *Service) List(ctx context.Context, lastKey *string) (*schema.FileList, error) {
	u := s.cfg.Hosts.API + "/api/user/files_uploaded"
	if lastKey != nil {
		u += "?lastKey=" + url.QueryEscape(*lastKey)
	}

	var dataResult struct {
		FileList   []schema.FileEntry `json:"fileList"`
		TotalFiles int                `json:"totalFiles"`
		LastKey    *string            `json:"lastKey,omitempty"`
	}

	_, err := s.h.WriteJSON(ctx, "GET", u, nil, &dataResult)
	if err != nil {
		return nil, err
	}

	return &schema.FileList{
		Data:       dataResult.FileList,
		TotalFiles: &dataResult.TotalFiles,
		LastKey:    dataResult.LastKey,
	}, nil
}

func (s *Service) Info(ctx context.Context, cid string) (*schema.FileInfo, error) {
	u := s.cfg.Hosts.API + "/api/lighthouse/file_info?cid=" + url.QueryEscape(cid)
	var out schema.FileInfo
	_, err := s.h.WriteJSON(ctx, "GET", u, nil, &out)
	return &out, err
}

func (s *Service) Pin(ctx context.Context, cid, name string) error {
	u := s.cfg.Hosts.API + "/api/lighthouse/pin"
	body := map[string]string{"cid": cid, "fileName": name}
	_, err := s.h.WriteJSON(ctx, "POST", u, body, nil)
	return err
}

func (s *Service) Delete(ctx context.Context, id string) error {
	u := s.cfg.Hosts.API + "/api/user/delete_file?id=" + url.QueryEscape(id)

	var response struct {
		Message string `json:"message"`
	}

	_, err := s.h.WriteJSON(ctx, "DELETE", u, nil, &response)
	if err != nil {
		return err
	}

	return nil
}
