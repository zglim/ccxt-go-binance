package binance

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type orderListPlaceDeprecatedServiceWsTestSuite struct {
	wsTestScaffold

	symbol            string
	side              SideType
	price             string
	quantity          string
	stopPrice         string
	listClientOrderID string

	orderListPlace        *OrderListPlaceWsService
	orderListPlaceRequest *OrderListPlaceWsRequest
}

func TestOrderListPlaceDeprecatedServiceWsPlace(t *testing.T) {
	suite.Run(t, new(orderListPlaceDeprecatedServiceWsTestSuite))
}

func (s *orderListPlaceDeprecatedServiceWsTestSuite) SetupTest() {
	s.initScaffold("e2a85d9f-07a5-4f94-8d5f-789dc3deb098")

	s.symbol = "BTCUSDT"
	s.side = SideTypeSell
	s.price = "23420.00000000"
	s.quantity = "0.00650000"
	s.stopPrice = "23410.00000000"
	s.listClientOrderID = "testOCOList"

	s.orderListPlace = &OrderListPlaceWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.orderListPlaceRequest = NewOrderListPlaceWsRequest().
		Symbol(s.symbol).
		Side(s.side).
		Price(s.price).
		Quantity(s.quantity).
		StopPrice(s.stopPrice).
		ListClientOrderID(s.listClientOrderID).
		NewOrderRespType(NewOrderRespTypeRESULT)
}

func (s *orderListPlaceDeprecatedServiceWsTestSuite) TearDownTest() {
	s.finishScaffold()
}

func (s *orderListPlaceDeprecatedServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.orderListPlace = &OrderListPlaceWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}

// --- Do tests ---

func (s *orderListPlaceDeprecatedServiceWsTestSuite) TestOrderListPlace() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.orderListPlace.Do(s.requestID, s.orderListPlaceRequest)
	s.NoError(err)
}

func (s *orderListPlaceDeprecatedServiceWsTestSuite) TestOrderListPlace_ErrorBranches() {
	s.assertDoErrorBranches(s.reset, func(reqID string) error {
		return s.orderListPlace.Do(reqID, s.orderListPlaceRequest)
	})
}

// --- SyncDo tests ---

func (s *orderListPlaceDeprecatedServiceWsTestSuite) TestOrderListPlaceSync() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	orderListPlaceResponse := CreateOrderListWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: CreateOrderListResult{
			OrderListId:       1274512,
			ContingencyType:   "OCO",
			ListStatusType:    "EXEC_STARTED",
			ListOrderStatus:   "EXECUTING",
			ListClientOrderId: s.listClientOrderID,
			TransactionTime:   1660801713793,
			Symbol:            s.symbol,
			Orders: []struct {
				Symbol        string `json:"symbol"`
				OrderId       int64  `json:"orderId"`
				ClientOrderId string `json:"clientOrderId"`
			}{
				{Symbol: s.symbol, OrderId: 12569138901, ClientOrderId: "BqtFCj5odMoWtSqGk2X9tU"},
				{Symbol: s.symbol, OrderId: 12569138902, ClientOrderId: "jLnZpj5enfMXTuhKB1d0us"},
			},
		},
	}

	rawResponseData, err := json.Marshal(orderListPlaceResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	req := s.orderListPlaceRequest
	response, err := s.orderListPlace.SyncDo(s.requestID, req)
	s.Require().NoError(err)
	s.Equal(*req.listClientOrderID, response.Result.ListClientOrderId)
	s.Equal(req.symbol, response.Result.Symbol)
	s.Equal("OCO", response.Result.ContingencyType)
}

func (s *orderListPlaceDeprecatedServiceWsTestSuite) TestOrderListPlaceSync_ErrorBranches() {
	s.assertSyncDoErrorBranches(s.reset, func(reqID string) (interface{}, error) {
		return s.orderListPlace.SyncDo(reqID, s.orderListPlaceRequest)
	})
}
