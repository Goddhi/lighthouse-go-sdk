package ipns

import (
	"context"
	"net/url"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/cfg"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/internal/httpx"
	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

type Service struct {
	h *httpx.Client
	cfg cfg.Config
}

func New(h *httpx.Client, c cfg.Config) *Service {
	return &Service{h: h, cfg: c}
}

func (s *Service) GenerateKey(ctx context.Context, keyName string) (*schema.IPNSKeyResponse, error) {
	u := s.cfg.Hosts.API + "/api/ipns/generate_key?keyName=" + url.QueryEscape(keyName)

	var response schema.IPNSKeyResponse
	_, err := s.h.WriteJSON(ctx, "GET", u, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *Service) PublishRecord(ctx context.Context, cid, keyName string) (*schema.IPNSPublishResponse, error) {
	u := s.cfg.Hosts.Upload + "/api/ipns/publish_recored?cid=" + url.QueryEscape(cid) + "&keyName=" + url.QueryEscape(keyName)

	var response schema.IPNSPublishResponse
	_, err := s.h.WriteJSON(ctx, "GET", u, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) ListKeys(ctx context.Context) ([]schema.IPNSRecord, error) {
	u := s.cfg.Hosts.API + "/api/ipns/get_ipns_records"
	
	var records []schema.IPNSRecord
	_, err := s.h.WriteJSON(ctx, "GET", u, nil, &records)
	if err != nil {
		return nil, err
	}
	
	return records, nil
}

func (s *Service) RemoveKey(ctx context.Context, keyName string) (*schema.IPNSRemoveResponse, error) {
	u := s.cfg.Hosts.API + "/api/ipns/remove_key?keyName=" + url.QueryEscape(keyName)

	var response schema.IPNSRemoveResponse
	_, err := s.h.WriteJSON(ctx, "DELETE", u, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}