package binance

import (
	"encoding/json"

	"github.com/adshao/go-binance/v2/common"
	"github.com/adshao/go-binance/v2/common/websocket"
)

// OrderListCancelWsService cancels order list
type OrderListCancelWsService struct {
	spotWsApiService
}

// NewOrderListCancelWsService init OrderListCancelWsService
func NewOrderListCancelWsService(apiKey, secretKey string) (*OrderListCancelWsService, error) {
	base, err := newSpotWsApiService(apiKey, secretKey)
	if err != nil {
		return nil, err
	}

	return &OrderListCancelWsService{
		spotWsApiService: *base,
	}, nil
}

// OrderListCancelWsRequest parameters for 'orderList.cancel' websocket API
type OrderListCancelWsRequest struct {
	symbol            string
	orderListID       *int64
	listClientOrderID *string
	newClientOrderID  *string
	recvWindow        *uint16
}

// NewOrderListCancelWsRequest init OrderListCancelWsRequest
func NewOrderListCancelWsRequest() *OrderListCancelWsRequest {
	return &OrderListCancelWsRequest{}
}

func (s *OrderListCancelWsRequest) GetParams() map[string]any {
	return s.buildParams()
}

// buildParams builds params
func (s *OrderListCancelWsRequest) buildParams() params {
	m := params{
		"symbol": s.symbol,
	}
	if s.orderListID != nil {
		m["orderListId"] = *s.orderListID
	}
	if s.listClientOrderID != nil {
		m["listClientOrderId"] = *s.listClientOrderID
	}
	if s.newClientOrderID != nil {
		m["newClientOrderId"] = *s.newClientOrderID
	}
	if s.recvWindow != nil {
		m["recvWindow"] = *s.recvWindow
	}
	return m
}

// Do - sends 'orderList.cancel' request
func (s *OrderListCancelWsService) Do(requestID string, request *OrderListCancelWsRequest) error {
	return s.sendRequest(requestID, websocket.OrderListCancelSpotWsApiMethod, request.buildParams())
}

// SyncDo - sends 'orderList.cancel' request and receives response
func (s *OrderListCancelWsService) SyncDo(requestID string, request *OrderListCancelWsRequest) (*CancelOrderListWsResponse, error) {
	response, err := s.sendSyncRequest(requestID, websocket.OrderListCancelSpotWsApiMethod, request.buildParams())
	if err != nil {
		return nil, err
	}

	cancelOrderListWsResponse := &CancelOrderListWsResponse{}
	if err := json.Unmarshal(response, cancelOrderListWsResponse); err != nil {
		return nil, err
	}

	return cancelOrderListWsResponse, nil
}

// Symbol set symbol
func (s *OrderListCancelWsRequest) Symbol(symbol string) *OrderListCancelWsRequest {
	s.symbol = symbol
	return s
}

// OrderListID set orderListID
func (s *OrderListCancelWsRequest) OrderListID(orderListID int64) *OrderListCancelWsRequest {
	s.orderListID = &orderListID
	return s
}

// ListClientOrderID set listClientOrderID
func (s *OrderListCancelWsRequest) ListClientOrderID(listClientOrderID string) *OrderListCancelWsRequest {
	s.listClientOrderID = &listClientOrderID
	return s
}

// NewClientOrderID set newClientOrderID
func (s *OrderListCancelWsRequest) NewClientOrderID(newClientOrderID string) *OrderListCancelWsRequest {
	s.newClientOrderID = &newClientOrderID
	return s
}

// RecvWindow set recvWindow
func (s *OrderListCancelWsRequest) RecvWindow(recvWindow uint16) *OrderListCancelWsRequest {
	s.recvWindow = &recvWindow
	return s
}

// CancelOrderListResult define order list cancellation result
type CancelOrderListResult struct {
	OrderListId       int64  `json:"orderListId"`
	ContingencyType   string `json:"contingencyType"`
	ListStatusType    string `json:"listStatusType"`
	ListOrderStatus   string `json:"listOrderStatus"`
	ListClientOrderId string `json:"listClientOrderId"`
	TransactionTime   int64  `json:"transactionTime"`
	Symbol            string `json:"symbol"`
	Orders            []struct {
		Symbol        string `json:"symbol"`
		OrderId       int64  `json:"orderId"`
		ClientOrderId string `json:"clientOrderId"`
	} `json:"orders"`
	OrderReports []struct {
		Symbol                  string          `json:"symbol"`
		OrderId                 int64           `json:"orderId"`
		OrderListId             int64           `json:"orderListId"`
		ClientOrderId           string          `json:"clientOrderId"`
		TransactTime            int64           `json:"transactTime"`
		Price                   string          `json:"price"`
		OrigQty                 string          `json:"origQty"`
		ExecutedQty             string          `json:"executedQty"`
		CummulativeQuoteQty     string          `json:"cummulativeQuoteQty"`
		Status                  OrderStatusType `json:"status"`
		TimeInForce             TimeInForceType `json:"timeInForce"`
		Type                    OrderType       `json:"type"`
		Side                    SideType        `json:"side"`
		StopPrice               string          `json:"stopPrice"`
		SelfTradePreventionMode string          `json:"selfTradePreventionMode"`
	} `json:"orderReports"`
}

// CancelOrderListWsResponse define 'orderList.cancel' websocket API response
type CancelOrderListWsResponse struct {
	Id     string                `json:"id"`
	Status int                   `json:"status"`
	Result CancelOrderListResult `json:"result"`

	// error response
	Error *common.APIError `json:"error,omitempty"`
}
