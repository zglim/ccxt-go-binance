package binance

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func (s *sorOrderPlaceServiceWsTestSuite) SetupTest() {
	s.setup("3a4437e2-41a3-4c19-897c-9cadc5dce8b6")

	s.symbol = "BTCUSDT"
	s.side = SideTypeBuy
	s.orderType = OrderTypeLimit
	s.quantity = "0.5"
	s.price = "31000"
	s.newClientOrderID = "sBI1KM6nNtOfj5tccZSKly"

	s.sorOrderPlace = &SorOrderPlaceWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.sorOrderPlaceRequest = NewSorOrderPlaceWsRequest().
		Symbol(s.symbol).
		Side(s.side).
		Type(s.orderType).
		Quantity(s.quantity).
		Price(s.price).
		NewClientOrderID(s.newClientOrderID).
		NewOrderRespType(NewOrderRespTypeFULL)
}

type sorOrderPlaceServiceWsTestSuite struct {
	baseOrderWsServiceTestSuite

	symbol           string
	side             SideType
	orderType        OrderType
	quantity         string
	price            string
	newClientOrderID string

	sorOrderPlace        *SorOrderPlaceWsService
	sorOrderPlaceRequest *SorOrderPlaceWsRequest
}

func TestSorOrderPlaceServiceWsPlace(t *testing.T) {
	suite.Run(t, new(sorOrderPlaceServiceWsTestSuite))
}

func (s *sorOrderPlaceServiceWsTestSuite) TestSorOrderPlace() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.sorOrderPlace.Do(s.requestID, s.sorOrderPlaceRequest)
	s.NoError(err)
}

func (s *sorOrderPlaceServiceWsTestSuite) TestSorOrderPlace_CredentialErrors() {
	s.runDoCredentialChecks(s.reset, func(requestID string) error {
		return s.sorOrderPlace.Do(requestID, s.sorOrderPlaceRequest)
	})
}

func (s *sorOrderPlaceServiceWsTestSuite) TestSorOrderPlaceSync() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	sorOrderPlaceResponse := SorOrderPlaceWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: []SorOrderPlaceResult{
			{
				Symbol:                  s.symbol,
				OrderId:                 2,
				OrderListId:             -1,
				ClientOrderId:           s.newClientOrderID,
				TransactTime:            1689149087774,
				Price:                   s.price,
				OrigQty:                 s.quantity,
				ExecutedQty:             s.quantity,
				CummulativeQuoteQty:     "14000.00000000",
				Status:                  OrderStatusTypeFilled,
				TimeInForce:             TimeInForceTypeGTC,
				Type:                    s.orderType,
				Side:                    s.side,
				WorkingTime:             1689149087774,
				WorkingFloor:            "SOR",
				SelfTradePreventionMode: "NONE",
				UsedSor:                 true,
				Fills: []struct {
					MatchType       string `json:"matchType"`
					Price           string `json:"price"`
					Qty             string `json:"qty"`
					Commission      string `json:"commission"`
					CommissionAsset string `json:"commissionAsset"`
					TradeId         int64  `json:"tradeId"`
					AllocId         int64  `json:"allocId"`
				}{
					{
						MatchType:       "ONE_PARTY_TRADE_REPORT",
						Price:           "28000.00000000",
						Qty:             s.quantity,
						Commission:      "0.00000000",
						CommissionAsset: "BTC",
						TradeId:         -1,
						AllocId:         0,
					},
				},
			},
		},
	}

	rawResponseData, err := json.Marshal(sorOrderPlaceResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	req := s.sorOrderPlaceRequest
	response, err := s.sorOrderPlace.SyncDo(s.requestID, req)
	s.Require().NoError(err)
	s.Equal(*req.newClientOrderID, response.Result[0].ClientOrderId)
	s.Equal(req.symbol, response.Result[0].Symbol)
	s.Equal(true, response.Result[0].UsedSor)
}

func (s *sorOrderPlaceServiceWsTestSuite) TestSorOrderPlaceSync_CredentialErrors() {
	s.runSyncDoCredentialChecks(s.reset, func(requestID string) error {
		response, err := s.sorOrderPlace.SyncDo(requestID, s.sorOrderPlaceRequest)
		s.Nil(response)
		return err
	})
}

func (s *sorOrderPlaceServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.sorOrderPlace = &SorOrderPlaceWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}
