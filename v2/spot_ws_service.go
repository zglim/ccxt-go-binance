package binance

import (
	"encoding/json"
	"time"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// spotWsService is the reusable skeleton shared by the spot websocket API
// services (order list place/oto/otoco, order list cancel, SOR order
// place/test). It owns the connection, the auth credentials and the common
// request/response plumbing so that a concrete service only has to declare its
// own method constant, request builder and response type.
//
// Concrete services embed spotWsService by value, which promotes both the
// credential fields (ApiKey, SecretKey, KeyType, TimeOffset) and the public
// channel accessors below, keeping each service's exported API unchanged.
type spotWsService struct {
	c          websocket.Client
	ApiKey     string
	SecretKey  string
	KeyType    string
	TimeOffset int64
}

// newSpotWsService opens a websocket connection and wires up a client shared by
// every spot websocket service constructor.
func newSpotWsService(apiKey, secretKey string) (spotWsService, error) {
	conn, err := websocket.NewConnection(WsApiInitReadWriteConn, WebsocketKeepalive, WebsocketTimeoutReadWriteConnection)
	if err != nil {
		return spotWsService{}, err
	}

	client, err := websocket.NewClient(conn)
	if err != nil {
		return spotWsService{}, err
	}

	return spotWsService{
		c:         client,
		ApiKey:    apiKey,
		SecretKey: secretKey,
		KeyType:   common.KeyTypeHmac,
	}, nil
}

// buildRequest signs and serializes a websocket API request for the given
// method and params. Both do and syncDo funnel through here, so request
// construction (and its error branch) lives in a single place.
func (s *spotWsService) buildRequest(requestID string, method websocket.WsApiMethodType, params params) ([]byte, error) {
	return websocket.CreateRequest(
		websocket.NewRequestData(
			requestID,
			s.ApiKey,
			s.SecretKey,
			s.TimeOffset,
			s.KeyType,
		),
		method,
		params,
	)
}

// do sends a request asynchronously; the response is delivered on the read
// channel returned by GetReadChannel.
func (s *spotWsService) do(requestID string, method websocket.WsApiMethodType, params params) error {
	rawData, err := s.buildRequest(requestID, method, params)
	if err != nil {
		return err
	}

	return s.c.Write(requestID, rawData)
}

// syncDo sends a request and blocks until the raw response (or an error) is
// returned. Decoding into a concrete response type is delegated to
// spotWsSyncDo so that each service keeps its own type-safe return value.
func (s *spotWsService) syncDo(requestID string, method websocket.WsApiMethodType, params params) ([]byte, error) {
	rawData, err := s.buildRequest(requestID, method, params)
	if err != nil {
		return nil, err
	}

	return s.c.WriteSync(requestID, rawData, websocket.WriteSyncWsTimeout)
}

// ReceiveAllDataBeforeStop waits until all responses will be received from websocket until timeout expired
func (s *spotWsService) ReceiveAllDataBeforeStop(timeout time.Duration) {
	s.c.Wait(timeout)
}

// GetReadChannel returns channel with API response data (including API errors)
func (s *spotWsService) GetReadChannel() <-chan []byte {
	return s.c.GetReadChannel()
}

// GetReadErrorChannel returns channel with errors which are occurred while reading websocket connection
func (s *spotWsService) GetReadErrorChannel() <-chan error {
	return s.c.GetReadErrorChannel()
}

// GetReconnectCount returns count of reconnect attempts by client
func (s *spotWsService) GetReconnectCount() int64 {
	return s.c.GetReconnectCount()
}

// spotWsSyncDo sends a request through s and decodes the raw response into a new
// value of T. Keeping the response type as an explicit type parameter lets each
// service's SyncDo stay a single, type-safe line while the shared skeleton owns
// both the send and the deserialization (and their error branches).
func spotWsSyncDo[T any](s *spotWsService, requestID string, method websocket.WsApiMethodType, params params) (*T, error) {
	rawData, err := s.syncDo(requestID, method, params)
	if err != nil {
		return nil, err
	}

	result := new(T)
	if err := json.Unmarshal(rawData, result); err != nil {
		return nil, err
	}

	return result, nil
}
