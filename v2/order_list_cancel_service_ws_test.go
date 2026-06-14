package binance

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func (s *orderListCancelServiceWsTestSuite) SetupTest() {
	s.setup("c5899911-d3f4-47ae-8835-97da553d27d0")

	s.symbol = "BTCUSDT"
	s.orderListID = int64(1274512)
	s.listClientOrderID = "6023531d7edaad348f5aff"

	s.orderListCancel = &OrderListCancelWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.orderListCancelRequest = NewOrderListCancelWsRequest().
		Symbol(s.symbol).
		OrderListID(s.orderListID).
		ListClientOrderID(s.listClientOrderID)
}

type orderListCancelServiceWsTestSuite struct {
	baseOrderWsServiceTestSuite

	symbol            string
	orderListID       int64
	listClientOrderID string

	orderListCancel        *OrderListCancelWsService
	orderListCancelRequest *OrderListCancelWsRequest
}

func TestOrderListCancelServiceWsPlace(t *testing.T) {
	suite.Run(t, new(orderListCancelServiceWsTestSuite))
}

func (s *orderListCancelServiceWsTestSuite) TestOrderListCancel() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.orderListCancel.Do(s.requestID, s.orderListCancelRequest)
	s.NoError(err)
}

func (s *orderListCancelServiceWsTestSuite) TestOrderListCancel_CredentialErrors() {
	s.runDoCredentialChecks(s.reset, func(requestID string) error {
		return s.orderListCancel.Do(requestID, s.orderListCancelRequest)
	})
}

func (s *orderListCancelServiceWsTestSuite) TestOrderListCancelSync() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	orderListCancelResponse := CancelOrderListWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: CancelOrderListResult{
			OrderListId:       s.orderListID,
			ContingencyType:   "OCO",
			ListStatusType:    "ALL_DONE",
			ListOrderStatus:   "ALL_DONE",
			ListClientOrderId: s.listClientOrderID,
			TransactionTime:   1660801720215,
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

	rawResponseData, err := json.Marshal(orderListCancelResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	req := s.orderListCancelRequest
	response, err := s.orderListCancel.SyncDo(s.requestID, req)
	s.Require().NoError(err)
	s.Equal(s.orderListID, response.Result.OrderListId)
	s.Equal(req.symbol, response.Result.Symbol)
	s.Equal("ALL_DONE", response.Result.ListStatusType)
}

func (s *orderListCancelServiceWsTestSuite) TestOrderListCancelSync_CredentialErrors() {
	s.runSyncDoCredentialChecks(s.reset, func(requestID string) error {
		response, err := s.orderListCancel.SyncDo(requestID, s.orderListCancelRequest)
		s.Nil(response)
		return err
	})
}

func (s *orderListCancelServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.orderListCancel = &OrderListCancelWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
