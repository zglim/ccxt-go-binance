package options

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/adshao/go-binance/v2/common"
)

const (
	pingEndpoint       = "/eapi/v1/ping"
	serverTimeEndpoint = "/eapi/v1/time"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: pingEndpoint,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	d := map[string]string{}
	return json.Unmarshal(data, &d)
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...RequestOption) (serverTime int64, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: serverTimeEndpoint,
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	return common.ParseServerTime(data)
}
