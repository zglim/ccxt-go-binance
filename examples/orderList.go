/*
Package main provides comprehensive examples of Binance WebSocket API order list operations.

This file demonstrates how to use all the new WebSocket services for order lists:

1. OrderListPlaceOCO() - Creates OCO (One-Cancels-Other) orders using the current API
2. OrderListPlaceOTO() - Creates OTO (One-Triggers-the-Other) orders
3. OrderListPlaceOTOCO() - Creates OTOCO (One-Triggers-One-Cancels-the-Other) orders
4. OrderListCancel() - Cancels existing order lists
5. SorOrderPlace() - Places orders using Smart Order Routing (SOR)
6. SorOrderTest() - Tests SOR orders without execution

Usage:
1. Set your API credentials in each function (replace "your_api_key" and "your_secret_key")
2. Enable testnet for testing (binance.UseTestnet = true)
3. Run individual functions or call RunOrderListExamples() to run all examples

Important Notes:
- All examples use testnet by default for safety
- WebSocket services support both synchronous (SyncDo) and asynchronous (Do) operations
- Request IDs are automatically generated using common.GenerateSpotId()
- SOR orders provide information about whether Smart Order Routing was used
- Commission rate computation can be enabled for SOR test orders

Shared helpers:
The boilerplate that every example repeats (testnet setup, credential validation,
service-creation guards, request-ID generation, error reporting, and waiting for a
response) lives in a handful of small helpers so that each example below can focus on
the part that actually differs: how its request is built.

  - setupExampleClient() (in config.go) - testnet setup, validation, client creation
  - runSyncExample()                     - the synchronous SyncDo flow
  - runAsyncExample() + waitForWsResponse() - the asynchronous Do + channel flow

Each helper is intentionally tiny and self-contained, so copying any single example
function and tweaking its request is still all it takes to get started.

Example Usage:

	// Run all examples
	RunOrderListExamples()

	// Or run individual examples
	OrderListPlaceOTO()
	SorOrderPlace()
*/
package main

import (
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
)

// runSyncExample wires together the steps shared by every synchronous order-list
// example: it guards against a failed service creation, generates a request ID,
// invokes SyncDo through the supplied closure, and reports any error. On success it
// returns the typed response so the caller can print exactly what it cares about.
//
// The closure keeps the call site readable while letting Go infer the response type:
//
//	response, ok := runSyncExample("OTO", err, func(requestID string) (*binance.CreateOrderListWsResponse, error) {
//		return service.SyncDo(requestID, request)
//	})
func runSyncExample[Resp any](label string, serviceErr error, do func(requestID string) (Resp, error)) (Resp, bool) {
	var zero Resp
	if serviceErr != nil {
		fmt.Printf("Error creating %s service: %v\n", label, serviceErr)
		return zero, false
	}

	requestID := common.GenerateSpotId()
	response, err := do(requestID)
	if err != nil {
		fmt.Printf("Error during %s request: %v\n", label, err)
		return zero, false
	}
	return response, true
}

// runAsyncExample wires together the steps shared by an asynchronous order-list
// example: it guards against a failed service creation, generates a request ID, and
// invokes Do through the supplied closure, reporting any error. On success it returns
// the request ID so the caller can wait for the matching response.
func runAsyncExample(label string, serviceErr error, do func(requestID string) error) (string, bool) {
	if serviceErr != nil {
		fmt.Printf("Error creating %s service: %v\n", label, serviceErr)
		return "", false
	}

	requestID := common.GenerateSpotId()
	if err := do(requestID); err != nil {
		fmt.Printf("Error during %s request: %v\n", label, err)
		return "", false
	}
	return requestID, true
}

// asyncWsService is the subset of an asynchronous WebSocket service that
// waitForWsResponse needs in order to read a response and then shut down cleanly.
type asyncWsService interface {
	GetReadChannel() <-chan []byte
	GetReadErrorChannel() <-chan error
	ReceiveAllDataBeforeStop(timeout time.Duration)
}

// waitForWsResponse listens for the first message (or error) on an asynchronous
// service's channels, prints it, and then drains any remaining data before stopping.
// It blocks for up to wait while the listener runs.
func waitForWsResponse(label string, service asyncWsService, wait time.Duration) {
	go func() {
		select {
		case response := <-service.GetReadChannel():
			fmt.Printf("%s Response: %s\n", label, string(response))
		case err := <-service.GetReadErrorChannel():
			fmt.Printf("%s Error: %v\n", label, err)
		}
	}()

	time.Sleep(wait)
	service.ReceiveAllDataBeforeStop(2 * time.Second)
}

// OrderListPlaceOCO demonstrates creating an OCO order using WebSocket API
func OrderListPlaceOCO() {
	client, ok := setupExampleClient()
	if !ok {
		return
	}

	// Create OCO WebSocket service
	service, err := client.NewOrderListCreateWsService()

	// Create OCO order request
	request := binance.NewOrderListCreateWsRequest().
		Symbol("BTCUSDT").
		Side(binance.SideTypeSell).
		Quantity("0.001").
		AboveType(binance.OrderTypeTakeProfitLimit).
		AbovePrice("115000").
		AboveStopPrice("115000").
		AboveTimeInForce(binance.TimeInForceTypeGTC).
		BelowType(binance.OrderTypeStopLossLimit).
		BelowPrice("110000").
		BelowStopPrice("110000").
		BelowTimeInForce(binance.TimeInForceTypeGTC).
		NewOrderRespType(binance.NewOrderRespTypeFULL)

	// Send async request
	requestID, ok := runAsyncExample("OCO", err, func(requestID string) error {
		return service.Do(requestID, request)
	})
	if !ok {
		return
	}

	fmt.Printf("OCO order sent with request ID: %s\n", requestID)

	// Listen for the response, then stop cleanly
	waitForWsResponse("OCO", service, 5*time.Second)
}

// OrderListPlaceOTO demonstrates creating an OTO order using WebSocket API
func OrderListPlaceOTO() {
	client, ok := setupExampleClient()
	if !ok {
		return
	}

	// Create OTO WebSocket service
	service, err := client.NewOrderListPlaceOtoWsService()

	// Create OTO order request
	request := binance.NewOrderListPlaceOtoWsRequest().
		Symbol("BTCUSDT").
		WorkingType(binance.OrderTypeLimit).
		WorkingSide(binance.SideTypeBuy).
		WorkingPrice("30000").
		WorkingQuantity("0.001").
		PendingType(binance.OrderTypeLimit).
		PendingSide(binance.SideTypeSell).
		WorkingTimeInForce(binance.TimeInForceTypeGTC).
		PendingTimeInForce(binance.TimeInForceTypeGTC).
		PendingPrice("32000").
		PendingQuantity("0.001")

	// Send synchronous request
	response, ok := runSyncExample("OTO", err, func(requestID string) (*binance.CreateOrderListWsResponse, error) {
		return service.SyncDo(requestID, request)
	})
	if !ok {
		return
	}

	fmt.Printf("OTO Order Response: %+v\n", response)
}

// OrderListPlaceOTOCO demonstrates creating an OTOCO order using WebSocket API
func OrderListPlaceOTOCO() {
	client, ok := setupExampleClient()
	if !ok {
		return
	}

	// Create OTOCO WebSocket service
	service, err := client.NewOrderListPlaceOtocoWsService()

	// Create OTOCO order request
	request := binance.NewOrderListPlaceOtocoWsRequest().
		Symbol("BTCUSDT").
		WorkingType(binance.OrderTypeLimit).
		WorkingSide(binance.SideTypeBuy).
		WorkingPrice("30000").
		WorkingQuantity("0.001").
		WorkingTimeInForce(binance.TimeInForceTypeGTC).
		PendingSide(binance.SideTypeSell).
		PendingQuantity("0.001").
		PendingAboveType(binance.OrderTypeLimitMaker).
		PendingAbovePrice("32000").
		PendingBelowType(binance.OrderTypeStopLoss).
		PendingBelowStopPrice("28000").
		ListClientOrderID("testOTOCOList")

	// Send synchronous request
	response, ok := runSyncExample("OTOCO", err, func(requestID string) (*binance.CreateOrderListWsResponse, error) {
		return service.SyncDo(requestID, request)
	})
	if !ok {
		return
	}

	fmt.Printf("OTOCO Order Response: %+v\n", response)
}

// OrderListCancel demonstrates canceling an order list using WebSocket API
func OrderListCancel() {
	client, ok := setupExampleClient()
	if !ok {
		return
	}

	// Create order list cancel WebSocket service
	service, err := client.NewOrderListCancelWsService()

	// Create cancel request
	request := binance.NewOrderListCancelWsRequest().
		Symbol("BTCUSDT").
		OrderListID(123456789) // Replace with actual order list ID

	// Send synchronous request
	response, ok := runSyncExample("Cancel Order List", err, func(requestID string) (*binance.CancelOrderListWsResponse, error) {
		return service.SyncDo(requestID, request)
	})
	if !ok {
		return
	}

	fmt.Printf("Cancel Order List Response: %+v\n", response)
}

// SorOrderPlace demonstrates placing a SOR order using WebSocket API
func SorOrderPlace() {
	client, ok := setupExampleClient()
	if !ok {
		return
	}

	// Create SOR order placement WebSocket service
	service, err := client.NewSorOrderPlaceWsService()

	// Create SOR order request - using ETHUSDT as it has SOR support
	request := binance.NewSorOrderPlaceWsRequest().
		Symbol("ETHUSDT").
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity("0.1").
		Price("2000").
		TimeInForce(binance.TimeInForceTypeGTC).
		NewClientOrderID("sBI1KM6nNtOfj5tccZSKly").
		NewOrderRespType(binance.NewOrderRespTypeFULL)

	// Send synchronous request
	response, ok := runSyncExample("SOR Place", err, func(requestID string) (*binance.SorOrderPlaceWsResponse, error) {
		return service.SyncDo(requestID, request)
	})
	if !ok {
		return
	}

	fmt.Printf("SOR Order Response: %+v\n", response)

	// Check if result array is not empty before accessing
	if len(response.Result) > 0 {
		fmt.Printf("Used SOR: %v\n", response.Result[0].UsedSor)
	} else {
		fmt.Printf("No order results returned\n")
	}
}

// SorOrderTest demonstrates testing a SOR order using WebSocket API
func SorOrderTest() {
	client, ok := setupExampleClient()
	if !ok {
		return
	}

	// Create SOR order test WebSocket service
	service, err := client.NewSorOrderTestWsService()

	// Create SOR test request with commission rates - using ETHUSDT for SOR support
	request := binance.NewSorOrderTestWsRequest().
		Symbol("ETHUSDT").
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity("0.1").
		Price("2000").
		TimeInForce(binance.TimeInForceTypeGTC).
		ComputeCommissionRates(true)

	// Send synchronous request
	response, ok := runSyncExample("SOR Test", err, func(requestID string) (*binance.SorOrderTestWsResponse, error) {
		return service.SyncDo(requestID, request)
	})
	if !ok {
		return
	}

	fmt.Printf("SOR Test Response: %+v\n", response)
	if response.Result.StandardCommissionForOrder != nil {
		fmt.Printf("Standard Commission - Maker: %s, Taker: %s\n",
			response.Result.StandardCommissionForOrder.Maker,
			response.Result.StandardCommissionForOrder.Taker)
	}
}

func RunOrderListExamples() {
	fmt.Println("=== Binance Order List WebSocket API Examples ===")

	fmt.Println("1. OCO Order (Current)")
	OrderListPlaceOCO()
	time.Sleep(2 * time.Second)

	fmt.Println("\n2. OTO Order")
	OrderListPlaceOTO()
	time.Sleep(2 * time.Second)

	fmt.Println("\n3. OTOCO Order")
	OrderListPlaceOTOCO()
	time.Sleep(2 * time.Second)

	fmt.Println("\n4. Cancel Order List")
	// OrderListCancel() // Uncomment and provide valid order list ID

	fmt.Println("\n=== Examples completed ===")
}
