package deals

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

func (s *Service) Status(ctx context.Context, cid string) ([]schema.DealStatus, error) {
    u := s.cfg.Hosts.API + "/api/lighthouse/deal_status?cid=" + url.QueryEscape(cid)

    var deals []schema.DealStatus
    _, err := s.h.WriteJSON(ctx, "GET", u, nil, &deals)
    if err != nil {
        return nil, err
    }
    return deals, nil  
}
