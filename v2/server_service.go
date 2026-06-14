package binance

import (
	"context"
	"net/http"

	"github.com/adshao/go-binance/v2/common"
)

const (
	serverPingEndpoint = "/api/v3/ping"
	serverTimeEndpoint = "/api/v3/time"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: serverPingEndpoint,
	}
	_, err = s.c.callAPI(ctx, r, opts...)
	return err
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
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	return common.ParseServerTime(data)
}

// SetServerTimeService set server time
type SetServerTimeService struct {
	c *Client
}

// Do send request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx)
	if err != nil {
		return 0, err
	}
	timeOffset = common.ComputeTimeOffset(currentTimestamp(), serverTime)
	s.c.TimeOffset = timeOffset
	return timeOffset, nil
}
