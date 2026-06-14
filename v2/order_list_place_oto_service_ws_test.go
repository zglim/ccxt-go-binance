package binance

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func (s *orderListPlaceOtoServiceWsTestSuite) SetupTest() {
	s.setup("1712544395950")

	s.symbol = "LTCBNB"
	s.workingType = OrderTypeLimit
	s.workingSide = SideTypeSell
	s.workingPrice = "1.0"
	s.workingQuantity = "1"
	s.pendingType = OrderTypeMarket
	s.pendingSide = SideTypeBuy
	s.pendingQuantity = "1"
	s.listClientOrderID = "testOTOList"

	s.orderListPlaceOto = &OrderListPlaceOtoWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.orderListPlaceOtoRequest = NewOrderListPlaceOtoWsRequest().
		Symbol(s.symbol).
		WorkingType(s.workingType).
		WorkingSide(s.workingSide).
		WorkingPrice(s.workingPrice).
		WorkingQuantity(s.workingQuantity).
		PendingType(s.pendingType).
		PendingSide(s.pendingSide).
		PendingQuantity(s.pendingQuantity).
		ListClientOrderID(s.listClientOrderID).
		NewOrderRespType(NewOrderRespTypeRESULT)
}

type orderListPlaceOtoServiceWsTestSuite struct {
	baseOrderWsServiceTestSuite

	symbol            string
	workingType       OrderType
	workingSide       SideType
	workingPrice      string
	workingQuantity   string
	pendingType       OrderType
	pendingSide       SideType
	pendingQuantity   string
	listClientOrderID string

	orderListPlaceOto        *OrderListPlaceOtoWsService
	orderListPlaceOtoRequest *OrderListPlaceOtoWsRequest
}

func TestOrderListPlaceOtoServiceWsPlace(t *testing.T) {
	suite.Run(t, new(orderListPlaceOtoServiceWsTestSuite))
}

func (s *orderListPlaceOtoServiceWsTestSuite) TestOrderListPlaceOto() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.orderListPlaceOto.Do(s.requestID, s.orderListPlaceOtoRequest)
	s.NoError(err)
}

func (s *orderListPlaceOtoServiceWsTestSuite) TestOrderListPlaceOto_CredentialErrors() {
	s.runDoCredentialChecks(s.reset, func(requestID string) error {
		return s.orderListPlaceOto.Do(requestID, s.orderListPlaceOtoRequest)
	})
}

func (s *orderListPlaceOtoServiceWsTestSuite) TestOrderListPlaceOtoSync() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	orderListPlaceOtoResponse := CreateOrderListWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: CreateOrderListResult{
			OrderListId:       626,
			ContingencyType:   "OTO",
			ListStatusType:    "EXEC_STARTED",
			ListOrderStatus:   "EXECUTING",
			ListClientOrderId: s.listClientOrderID,
			TransactionTime:   1712544395981,
			Symbol:            s.symbol,
			Orders: []struct {
				Symbol        string `json:"symbol"`
				OrderId       int64  `json:"orderId"`
				ClientOrderId string `json:"clientOrderId"`
			}{
				{Symbol: s.symbol, OrderId: 13, ClientOrderId: "YiAUtM9yJjl1a2jXHSp9Ny"},
				{Symbol: s.symbol, OrderId: 14, ClientOrderId: "9MxJSE1TYkmyx5lbGLve7R"},
			},
		},
	}

	rawResponseData, err := json.Marshal(orderListPlaceOtoResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	req := s.orderListPlaceOtoRequest
	response, err := s.orderListPlaceOto.SyncDo(s.requestID, req)
	s.Require().NoError(err)
	s.Equal(*req.listClientOrderID, response.Result.ListClientOrderId)
	s.Equal(req.symbol, response.Result.Symbol)
	s.Equal("OTO", response.Result.ContingencyType)
}

func (s *orderListPlaceOtoServiceWsTestSuite) TestOrderListPlaceOtoSync_CredentialErrors() {
	s.runSyncDoCredentialChecks(s.reset, func(requestID string) error {
		response, err := s.orderListPlaceOto.SyncDo(requestID, s.orderListPlaceOtoRequest)
		s.Nil(response)
		return err
	})
}

func (s *orderListPlaceOtoServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.orderListPlaceOto = &OrderListPlaceOtoWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
