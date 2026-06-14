package binance

import (
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// spotWsApiService provides common websocket service capabilities for spot websocket API services.
// It encapsulates connection initialization, request building, async/sync sending, and channel proxying.
// Concrete services embed this struct and only define their service-specific Do/SyncDo methods.
type spotWsApiService struct {
	c          websocket.Client
	ApiKey     string
	SecretKey  string
	KeyType    string
	TimeOffset int64
}

// newSpotWsApiService creates a spotWsApiService with a live websocket connection.
func newSpotWsApiService(apiKey, secretKey string) (*spotWsApiService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return nil, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &spotWsApiService{
		c:         client,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
	}, nil
}

// sendRequest builds a signed websocket request and sends it asynchronously.
func (s *spotWsApiService) sendRequest(requestID string, method websocket.WsApiMethodType, p params) error {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		method,
		p,
	)
	if err != nil {
		return err
	}

	return s.c.Write(requestID, rawData)
}

// sendSyncRequest builds a signed websocket request, sends it, and waits for the raw response bytes.
func (s *spotWsApiService) sendSyncRequest(requestID string, method websocket.WsApiMethodType, p params) ([]byte, error) {
	rawData, err := websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		method,
		p,
	)
	if err != nil {
		return nil, err
	}

	return s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *spotWsApiService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *spotWsApiService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *spotWsApiService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *spotWsApiService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}
