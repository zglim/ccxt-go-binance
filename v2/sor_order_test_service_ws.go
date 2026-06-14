package binance

import (
	"encoding/json"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// SorOrderTestWsService tests order using SOR
type SorOrderTestWsService struct {
	spotWsApiService
}

// NewSorOrderTestWsService init SorOrderTestWsService
func NewSorOrderTestWsService(apiKey, secretKey string) (*SorOrderTestWsService, error) {
	base, err := newSpotWsApiService(apiKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &SorOrderTestWsService{
		spotWsApiService: *base,
	}, nil
}

// SorOrderTestWsRequest parameters for 'sor.order.test' websocket API
type SorOrderTestWsRequest struct {
	symbol                  string
	side                    SideType
	orderType               OrderType
	timeInForce             *TimeInForceType
	price                   *string
	quantity                string
	newClientOrderID        *string
	icebergQty              *string
	strategyId              *int64
	strategyType            *int32
	selfTradePreventionMode *SelfTradePreventionMode
	computeCommissionRates  *bool
	recvWindow              *uint16
}

// NewSorOrderTestWsRequest init SorOrderTestWsRequest
func NewSorOrderTestWsRequest() *SorOrderTestWsRequest {
	return &SorOrderTestWsRequest{}
}

func (s *SorOrderTestWsRequest) GetParams() map[string]any {
	return s.buildParams()
}

// buildParams builds params
func (s *SorOrderTestWsRequest) buildParams() params {
	m := params{
		"symbol":   s.symbol,
		"side":     s.side,
		"type":     s.orderType,
		"quantity": s.quantity,
	}
	if s.timeInForce != nil {
		m["timeInForce"] = *s.timeInForce
	}
	if s.price != nil {
		m["price"] = *s.price
	}
	if s.newClientOrderID != nil {
		m["newClientOrderId"] = *s.newClientOrderID
	}
	if s.icebergQty != nil {
		m["icebergQty"] = *s.icebergQty
	}
	if s.strategyId != nil {
		m["strategyId"] = *s.strategyId
	}
	if s.strategyType != nil {
		m["strategyType"] = *s.strategyType
	}
	if s.selfTradePreventionMode != nil {
		m["selfTradePreventionMode"] = *s.selfTradePreventionMode
	}
	if s.computeCommissionRates != nil {
		m["computeCommissionRates"] = *s.computeCommissionRates
	}
	if s.recvWindow != nil {
		m["recvWindow"] = *s.recvWindow
	}
	return m
}

// Do - sends 'sor.order.test' request
func (s *SorOrderTestWsService) Do(requestID string, request *SorOrderTestWsRequest) error {
	return s.sendRequest(requestID, websocket.SorOrderTestSpotWsApiMethod, request.buildParams())
}

// SyncDo - sends 'sor.order.test' request and receives response
func (s *SorOrderTestWsService) SyncDo(requestID string, request *SorOrderTestWsRequest) (*SorOrderTestWsResponse, error) {
	response, err := s.sendSyncRequest(requestID, websocket.SorOrderTestSpotWsApiMethod, request.buildParams())
	if err != nil {
		return nil, err
	}

	sorOrderTestWsResponse := &SorOrderTestWsResponse{}
	if err := json.Unmarshal(response, sorOrderTestWsResponse); err != nil {
		return nil, err
	}

	return sorOrderTestWsResponse, nil
}

// Symbol set symbol
func (s *SorOrderTestWsRequest) Symbol(symbol string) *SorOrderTestWsRequest {
	s.symbol = symbol
	return s
}

// Side set side
func (s *SorOrderTestWsRequest) Side(side SideType) *SorOrderTestWsRequest {
	s.side = side
	return s
}

// Type set orderType
func (s *SorOrderTestWsRequest) Type(orderType OrderType) *SorOrderTestWsRequest {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *SorOrderTestWsRequest) TimeInForce(timeInForce TimeInForceType) *SorOrderTestWsRequest {
	s.timeInForce = &timeInForce
	return s
}

// Price set price
func (s *SorOrderTestWsRequest) Price(price string) *SorOrderTestWsRequest {
	s.price = &price
	return s
}

// Quantity set quantity
func (s *SorOrderTestWsRequest) Quantity(quantity string) *SorOrderTestWsRequest {
	s.quantity = quantity
	return s
}

// NewClientOrderID set newClientOrderID
func (s *SorOrderTestWsRequest) NewClientOrderID(newClientOrderID string) *SorOrderTestWsRequest {
	s.newClientOrderID = &newClientOrderID
	return s
}

// IcebergQty set icebergQty
func (s *SorOrderTestWsRequest) IcebergQty(icebergQty string) *SorOrderTestWsRequest {
	s.icebergQty = &icebergQty
	return s
}

// StrategyId set strategyId
func (s *SorOrderTestWsRequest) StrategyId(strategyId int64) *SorOrderTestWsRequest {
	s.strategyId = &strategyId
	return s
}

// StrategyType set strategyType
func (s *SorOrderTestWsRequest) StrategyType(strategyType int32) *SorOrderTestWsRequest {
	s.strategyType = &strategyType
	return s
}

// SelfTradePreventionMode set selfTradePreventionMode
func (s *SorOrderTestWsRequest) SelfTradePreventionMode(selfTradePreventionMode SelfTradePreventionMode) *SorOrderTestWsRequest {
	s.selfTradePreventionMode = &selfTradePreventionMode
	return s
}

// ComputeCommissionRates set computeCommissionRates
func (s *SorOrderTestWsRequest) ComputeCommissionRates(computeCommissionRates bool) *SorOrderTestWsRequest {
	s.computeCommissionRates = &computeCommissionRates
	return s
}

// RecvWindow set recvWindow
func (s *SorOrderTestWsRequest) RecvWindow(recvWindow uint16) *SorOrderTestWsRequest {
	s.recvWindow = &recvWindow
	return s
}

// SorOrderTestResult define SOR order test result
type SorOrderTestResult struct {
	StandardCommissionForOrder *struct {
		Maker string `json:"maker"`
		Taker string `json:"taker"`
	} `json:"standardCommissionForOrder,omitempty"`
	TaxCommissionForOrder *struct {
		Maker string `json:"maker"`
		Taker string `json:"taker"`
	} `json:"taxCommissionForOrder,omitempty"`
	Discount *struct {
		EnabledForAccount bool   `json:"enabledForAccount"`
		EnabledForSymbol  bool   `json:"enabledForSymbol"`
		DiscountAsset     string `json:"discountAsset"`
		Discount          string `json:"discount"`
	} `json:"discount,omitempty"`
}

// SorOrderTestWsResponse define 'sor.order.test' websocket API response
type SorOrderTestWsResponse struct {
	Id     string             `json:"id"`
	Status int                `json:"status"`
	Result SorOrderTestResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
