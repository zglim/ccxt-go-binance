package binance

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2/common/websocket"
	"github.com/adshao/go-binance/v2/common/websocket/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

// spotWsApiServiceTestSuite tests the common spotWsApiService base skeleton
// that is shared by all spot websocket services.
func (s *spotWsApiServiceTestSuite) SetupTest() {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = int64(0)
	s.requestID = "base-test-request-id"

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)

	s.base = &spotWsApiService{
		c:          s.client,
		ApiKey:     s.apiKey,
		SecretKey:  s.secretKey,
		KeyType:    s.signedKey,
		TimeOffset: s.timeOffset,
	}
}

func (s *spotWsApiServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

type spotWsApiServiceTestSuite struct {
	suite.Suite
	apiKey    string
	secretKey string
	signedKey string
	timeOffset int64

	ctrl   *gomock.Controller
	client *mock.MockClient

	requestID string
	base      *spotWsApiService
}

func TestSpotWsApiBaseService(t *testing.T) {
	suite.Run(t, new(spotWsApiServiceTestSuite))
}

// --- sendRequest tests ---

func (s *spotWsApiServiceTestSuite) TestSendRequest_Success() {
	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(1)

	err := s.base.sendRequest(s.requestID, websocket.OrderListPlaceSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.NoError(err)
}

func (s *spotWsApiServiceTestSuite) TestSendRequest_EmptyRequestID() {
	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.base.sendRequest("", websocket.OrderListPlaceSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *spotWsApiServiceTestSuite) TestSendRequest_EmptyApiKey() {
	s.base.ApiKey = ""
	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.base.sendRequest(s.requestID, websocket.OrderListPlaceSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
}

func (s *spotWsApiServiceTestSuite) TestSendRequest_EmptySecretKey() {
	s.base.SecretKey = ""
	s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)

	err := s.base.sendRequest(s.requestID, websocket.OrderListPlaceSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
}

func (s *spotWsApiServiceTestSuite) TestSendRequest_WriteError() {
	writeErr := fmt.Errorf("write failed")
	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(writeErr).Times(1)

	err := s.base.sendRequest(s.requestID, websocket.OrderListPlaceSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.ErrorIs(err, writeErr)
}

// --- sendSyncRequest tests ---

func (s *spotWsApiServiceTestSuite) TestSendSyncRequest_Success() {
	expectedResponse := map[string]any{"id": s.requestID, "status": 200, "result": map[string]any{}}
	rawResponse, err := json.Marshal(expectedResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponse, nil).Times(1)

	resp, err := s.base.sendSyncRequest(s.requestID, websocket.OrderListCancelSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.Require().NoError(err)
	s.Equal(rawResponse, resp)
}

func (s *spotWsApiServiceTestSuite) TestSendSyncRequest_EmptyRequestID() {
	s.client.EXPECT().WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(0)

	resp, err := s.base.sendSyncRequest("", websocket.OrderListCancelSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.Nil(resp)
	s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
}

func (s *spotWsApiServiceTestSuite) TestSendSyncRequest_EmptyApiKey() {
	s.base.ApiKey = ""
	s.client.EXPECT().WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(0)

	resp, err := s.base.sendSyncRequest(s.requestID, websocket.OrderListCancelSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.Nil(resp)
	s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
}

func (s *spotWsApiServiceTestSuite) TestSendSyncRequest_WriteSyncError() {
	writeSyncErr := fmt.Errorf("write sync failed")
	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(nil, writeSyncErr).Times(1)

	resp, err := s.base.sendSyncRequest(s.requestID, websocket.OrderListCancelSpotWsApiMethod, params{"symbol": "BTCUSDT"})
	s.Nil(resp)
	s.ErrorIs(err, writeSyncErr)
}

// --- Channel proxy method tests ---

func (s *spotWsApiServiceTestSuite) TestGetReadChannel() {
	readC := make(chan []byte)
	s.client.EXPECT().GetReadChannel().Return(readC).Times(1)

	ch := s.base.GetReadChannel()
	s.NotNil(ch)
}

func (s *spotWsApiServiceTestSuite) TestGetReadErrorChannel() {
	errC := make(chan error)
	s.client.EXPECT().GetReadErrorChannel().Return(errC).Times(1)

	ch := s.base.GetReadErrorChannel()
	s.NotNil(ch)
}

func (s *spotWsApiServiceTestSuite) TestGetReconnectCount() {
	s.client.EXPECT().GetReconnectCount().Return(int64(3)).Times(1)

	count := s.base.GetReconnectCount()
	s.Equal(int64(3), count)
}

func (s *spotWsApiServiceTestSuite) TestReceiveAllDataBeforeStop() {
	s.client.EXPECT().Wait(5 * time.Second).Times(1)

	s.base.ReceiveAllDataBeforeStop(5 * time.Second)
}

// --- Embedded promotion test: verify that a concrete service inherits base methods ---

func (s *spotWsApiServiceTestSuite) TestConcreteServiceInheritsBaseMethods() {
	// Create a concrete service using the embedded base
	svc := &OrderListCancelWsService{
		spotWsApiService: spotWsApiService{
			c:         s.client,
			ApiKey:    s.apiKey,
			SecretKey: s.secretKey,
			KeyType:   s.signedKey,
		},
	}

	// Verify channel proxy methods are accessible through the concrete type
	s.client.EXPECT().GetReadChannel().Return(make(chan []byte)).Times(1)
	s.NotNil(svc.GetReadChannel())

	s.client.EXPECT().GetReadErrorChannel().Return(make(chan error)).Times(1)
	s.NotNil(svc.GetReadErrorChannel())

	s.client.EXPECT().GetReconnectCount().Return(int64(7)).Times(1)
	s.Equal(int64(7), svc.GetReconnectCount())

	s.client.EXPECT().Wait(2 * time.Second).Times(1)
	svc.ReceiveAllDataBeforeStop(2 * time.Second)
}
