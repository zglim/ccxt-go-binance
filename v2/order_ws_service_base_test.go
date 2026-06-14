package binance

import (
	"fmt"

	"github.com/adshao/go-binance/v2/common/websocket"
	"github.com/adshao/go-binance/v2/common/websocket/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

// baseOrderWsServiceTestSuite collects the scaffolding shared by every spot
// order WebSocket service test: the dummy credentials, the gomock controller
// and mock client, and the credential error-branch assertions for Do/SyncDo.
//
// Concrete suites embed it, set up their own request/response fixtures and
// service under test in SetupTest, and hand the shared credential checks a
// reset func plus do/syncDo closures so the common assertions can drive each
// specific service while every file keeps its own request params and response
// expectations visible.
type baseOrderWsServiceTestSuite struct {
	suite.Suite

	apiKey     string
	secretKey  string
	signedKey  string
	timeOffset int64
	requestID  string

	ctrl   *gomock.Controller
	client *mock.MockClient
}

// setup initializes the dummy credentials and the gomock client shared by all
// order WS service suites. Concrete suites call it first from their SetupTest,
// passing the request ID they want to exercise, before wiring up their own
// request fixtures and service under test.
func (s *baseOrderWsServiceTestSuite) setup(requestID string) {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = 0
	s.requestID = requestID

	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)
}

// TearDownTest finishes the gomock controller. It is promoted to every concrete
// suite that embeds baseOrderWsServiceTestSuite.
func (s *baseOrderWsServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

// resetFunc rebuilds the service under test with the given credentials so a
// single credential error branch can be exercised in isolation. Each concrete
// suite supplies one that targets its own service field.
type resetFunc func(apiKey, secretKey, signKeyType string, timeOffset int64)

// runDoCredentialChecks exercises the credential error branches shared by every
// Do implementation: empty request ID, API key, secret key and sign key type.
// The do closure invokes the concrete service's Do with that suite's request
// fixture; none of these branches should ever reach the transport, so Write is
// expected zero times.
func (s *baseOrderWsServiceTestSuite) runDoCredentialChecks(reset resetFunc, do func(requestID string) error) {
	s.Run("EmptyRequestID", func() {
		reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)
		s.ErrorIs(do(""), websocket.ErrorRequestIDNotSet)
	})
	s.Run("EmptyApiKey", func() {
		reset("", s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)
		s.ErrorIs(do(s.requestID), websocket.ErrorApiKeyIsNotSet)
	})
	s.Run("EmptySecretKey", func() {
		reset(s.apiKey, "", s.signedKey, s.timeOffset)
		s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)
		s.ErrorIs(do(s.requestID), websocket.ErrorSecretKeyIsNotSet)
	})
	s.Run("EmptySignKeyType", func() {
		reset(s.apiKey, s.secretKey, "", s.timeOffset)
		s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)
		s.Error(do(s.requestID))
	})
}

// runSyncDoCredentialChecks exercises the credential error branches shared by
// every SyncDo implementation. The syncDo closure invokes the concrete
// service's SyncDo and is responsible for asserting that the typed response is
// nil, keeping each service's response type visible in its own test file.
func (s *baseOrderWsServiceTestSuite) runSyncDoCredentialChecks(reset resetFunc, syncDo func(requestID string) error) {
	s.Run("EmptyRequestID", func() {
		reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().
			WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		s.ErrorIs(syncDo(""), websocket.ErrorRequestIDNotSet)
	})
	s.Run("EmptyApiKey", func() {
		reset("", s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().
			WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		s.ErrorIs(syncDo(s.requestID), websocket.ErrorApiKeyIsNotSet)
	})
	s.Run("EmptySecretKey", func() {
		reset(s.apiKey, "", s.signedKey, s.timeOffset)
		s.client.EXPECT().
			WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		s.ErrorIs(syncDo(s.requestID), websocket.ErrorSecretKeyIsNotSet)
	})
	s.Run("EmptySignKeyType", func() {
		reset(s.apiKey, s.secretKey, "", s.timeOffset)
		s.client.EXPECT().
			WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		s.Error(syncDo(s.requestID))
	})
}
