package binance

import (
	"fmt"

	"github.com/adshao/go-binance/v2/common/websocket"
	"github.com/adshao/go-binance/v2/common/websocket/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

// wsTestScaffold provides reusable scaffolding for websocket service tests.
// Embed it in concrete test suites to share credential fields, mock lifecycle,
// and standard error-branch assertions for Do / SyncDo.
//
// Usage:
//
//	type myServiceWsTestSuite struct {
//	    wsTestScaffold
//	    // service-specific fields …
//	}
//
//	func (s *myServiceWsTestSuite) SetupTest() {
//	    s.initScaffold("some-request-id")
//	    // initialise service-specific fields and build the request …
//	}
//
//	func (s *myServiceWsTestSuite) TearDownTest() {
//	    s.finishScaffold()
//	}
type wsTestScaffold struct {
	suite.Suite
	apiKey     string
	secretKey  string
	signedKey  string
	timeOffset int64
	requestID  string
	ctrl       *gomock.Controller
	client     *mock.MockClient
}

// initScaffold sets default credentials and creates the mock controller.
// Call this from the concrete suite's SetupTest.
func (s *wsTestScaffold) initScaffold(requestID string) {
	s.apiKey = "dummyApiKey"
	s.secretKey = "dummySecretKey"
	s.signedKey = "HMAC"
	s.timeOffset = 0
	s.requestID = requestID
	s.ctrl = gomock.NewController(s.T())
	s.client = mock.NewMockClient(s.ctrl)
}

// finishScaffold tears down the mock controller.
// Call this from the concrete suite's TearDownTest.
func (s *wsTestScaffold) finishScaffold() {
	s.ctrl.Finish()
}

// wsResetFunc resets the service under test with the given credentials.
type wsResetFunc func(apiKey, secretKey, signKeyType string, timeOffset int64)

// wsDoFunc calls Do on the service under test with the given requestID.
type wsDoFunc func(requestID string) error

// wsSyncDoFunc calls SyncDo on the service under test with the given requestID.
// The response is returned as interface{}; it will be nil on error paths.
type wsSyncDoFunc func(requestID string) (interface{}, error)

// assertDoErrorBranches covers the 4 standard Do error branches:
// empty requestID, empty API key, empty secret key, empty sign key type.
func (s *wsTestScaffold) assertDoErrorBranches(reset wsResetFunc, do wsDoFunc) {
	s.Run("EmptyRequestID", func() {
		reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil).Times(0)
		err := do("")
		s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
	})
	s.Run("EmptyApiKey", func() {
		reset("", s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)
		err := do(s.requestID)
		s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
	})
	s.Run("EmptySecretKey", func() {
		reset(s.apiKey, "", s.signedKey, s.timeOffset)
		s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)
		err := do(s.requestID)
		s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
	})
	s.Run("EmptySignKeyType", func() {
		reset(s.apiKey, s.secretKey, "", s.timeOffset)
		s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).Times(0)
		err := do(s.requestID)
		s.Error(err)
	})
}

// assertSyncDoErrorBranches covers the 4 standard SyncDo error branches:
// empty requestID, empty API key, empty secret key, empty sign key type.
func (s *wsTestScaffold) assertSyncDoErrorBranches(reset wsResetFunc, syncDo wsSyncDoFunc) {
	s.Run("EmptyRequestID", func() {
		reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().
			WriteSync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		resp, err := syncDo("")
		s.Nil(resp)
		s.ErrorIs(err, websocket.ErrorRequestIDNotSet)
	})
	s.Run("EmptyApiKey", func() {
		reset("", s.secretKey, s.signedKey, s.timeOffset)
		s.client.EXPECT().
			WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		resp, err := syncDo(s.requestID)
		s.Nil(resp)
		s.ErrorIs(err, websocket.ErrorApiKeyIsNotSet)
	})
	s.Run("EmptySecretKey", func() {
		reset(s.apiKey, "", s.signedKey, s.timeOffset)
		s.client.EXPECT().
			WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		resp, err := syncDo(s.requestID)
		s.Nil(resp)
		s.ErrorIs(err, websocket.ErrorSecretKeyIsNotSet)
	})
	s.Run("EmptySignKeyType", func() {
		reset(s.apiKey, s.secretKey, "", s.timeOffset)
		s.client.EXPECT().
			WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("write sync: error")).Times(0)
		resp, err := syncDo(s.requestID)
		s.Nil(resp)
		s.Error(err)
	})
}
