package binance

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type sorOrderTestServiceWsTestSuite struct {
	wsTestScaffold

	symbol                 string
	side                   SideType
	orderType              OrderType
	quantity               string
	price                  string
	computeCommissionRates bool

	sorOrderTest        *SorOrderTestWsService
	sorOrderTestRequest *SorOrderTestWsRequest
}

func TestSorOrderTestServiceWsPlace(t *testing.T) {
	suite.Run(t, new(sorOrderTestServiceWsTestSuite))
}

func (s *sorOrderTestServiceWsTestSuite) SetupTest() {
	s.initScaffold("3a4437e2-41a3-4c19-897c-9cadc5dce8b6")

	s.symbol = "BTCUSDT"
	s.side = SideTypeBuy
	s.orderType = OrderTypeLimit
	s.quantity = "0.1"
	s.price = "0.1"
	s.computeCommissionRates = false

	s.sorOrderTest = &SorOrderTestWsService{
		c:         s.client,
		ApiKey:    s.apiKey,
		SecretKey: s.secretKey,
		KeyType:   s.signedKey,
	}

	s.sorOrderTestRequest = NewSorOrderTestWsRequest().
		Symbol(s.symbol).
		Side(s.side).
		Type(s.orderType).
		Quantity(s.quantity).
		Price(s.price).
		ComputeCommissionRates(s.computeCommissionRates)
}

func (s *sorOrderTestServiceWsTestSuite) TearDownTest() {
	s.finishScaffold()
}

func (s *sorOrderTestServiceWsTestSuite) reset(apiKey, secretKey, signKeyType string, timeOffset int64) {
	s.sorOrderTest = &SorOrderTestWsService{
		c:          s.client,
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		KeyType:    signKeyType,
		TimeOffset: timeOffset,
	}
}

// --- Do tests ---

func (s *sorOrderTestServiceWsTestSuite) TestSorOrderTest() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	s.client.EXPECT().Write(s.requestID, gomock.Any()).Return(nil).AnyTimes()

	err := s.sorOrderTest.Do(s.requestID, s.sorOrderTestRequest)
	s.NoError(err)
}

func (s *sorOrderTestServiceWsTestSuite) TestSorOrderTest_ErrorBranches() {
	s.assertDoErrorBranches(s.reset, func(reqID string) error {
		return s.sorOrderTest.Do(reqID, s.sorOrderTestRequest)
	})
}

// --- SyncDo tests ---

func (s *sorOrderTestServiceWsTestSuite) TestSorOrderTestSync_WithoutCommissionRates() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	sorOrderTestResponse := SorOrderTestWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: SorOrderTestResult{},
	}

	rawResponseData, err := json.Marshal(sorOrderTestResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.sorOrderTest.SyncDo(s.requestID, s.sorOrderTestRequest)
	s.Require().NoError(err)
	s.Equal(s.requestID, response.Id)
	s.Equal(200, response.Status)
}

func (s *sorOrderTestServiceWsTestSuite) TestSorOrderTestSync_WithCommissionRates() {
	s.reset(s.apiKey, s.secretKey, s.signedKey, s.timeOffset)

	// Update request to include commission rates
	s.sorOrderTestRequest.ComputeCommissionRates(true)

	sorOrderTestResponse := SorOrderTestWsResponse{
		Id:     s.requestID,
		Status: 200,
		Result: SorOrderTestResult{
			StandardCommissionForOrder: &struct {
				Maker string `json:"maker"`
				Taker string `json:"taker"`
			}{
				Maker: "0.00000112",
				Taker: "0.00000114",
			},
			TaxCommissionForOrder: &struct {
				Maker string `json:"maker"`
				Taker string `json:"taker"`
			}{
				Maker: "0.00000112",
				Taker: "0.00000114",
			},
			Discount: &struct {
				EnabledForAccount bool   `json:"enabledForAccount"`
				EnabledForSymbol  bool   `json:"enabledForSymbol"`
				DiscountAsset     string `json:"discountAsset"`
				Discount          string `json:"discount"`
			}{
				EnabledForAccount: true,
				EnabledForSymbol:  true,
				DiscountAsset:     "BNB",
				Discount:          "0.25",
			},
		},
	}

	rawResponseData, err := json.Marshal(sorOrderTestResponse)
	s.NoError(err)

	s.client.EXPECT().WriteSync(s.requestID, gomock.Any(), gomock.Any()).Return(rawResponseData, nil).Times(1)

	response, err := s.sorOrderTest.SyncDo(s.requestID, s.sorOrderTestRequest)
	s.Require().NoError(err)
	s.Equal(s.requestID, response.Id)
	s.Equal(200, response.Status)
	s.NotNil(response.Result.StandardCommissionForOrder)
	s.Equal("0.00000112", response.Result.StandardCommissionForOrder.Maker)
	s.Equal("BNB", response.Result.Discount.DiscountAsset)
}

func (s *sorOrderTestServiceWsTestSuite) TestSorOrderTestSync_ErrorBranches() {
	s.assertSyncDoErrorBranches(s.reset, func(reqID string) (interface{}, error) {
		return s.sorOrderTest.SyncDo(reqID, s.sorOrderTestRequest)
	})
}
