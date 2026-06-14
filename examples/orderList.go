/*
Package main provides comprehensive examples of Binance WebSocket API order list operations.

This file demonstrates how to use all the new WebSocket services for order lists:

1. OrderListPlaceOCO() - Creates OCO (One-Cancels-Other) orders using the current API
2. OrderListPlaceOTO() - Creates OTO (One-Triggers-the-Other) orders
3. OrderListPlaceOTOCO() - Creates OTOCO (One-Triggers-One-Cancels-the-Other) orders
4. OrderListCancel() - Cancels existing order lists
5. SorOrderPlace() - Places orders using Smart Order Routing (SOR)
6. SorOrderTest() - Tests SOR orders without execution
7. OrderListPlaceDeprecated() - Uses the deprecated OCO endpoint

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

// ---------------------------------------------------------------------------
// Example helpers – keep boilerplate out of each demo function
// ---------------------------------------------------------------------------

// newClient sets up testnet config, validates credentials, and returns a ready
// Binance client.  Every example starts with these three steps, so we fold
// them into one call.  Returns nil when setup fails (error is already printed).
func newClient(label string) *binance.Client {
	AppConfig.SetupTestnet()
	if err := AppConfig.Validate(); err != nil {
		fmt.Printf("[%s] Configuration error: %v\n", label, err)
		return nil
	}
	return AppConfig.GetClient()
}

// newRequestID generates a unique spot request ID and prints it for tracing.
func newRequestID(label string) string {
	id := common.GenerateSpotId()
	fmt.Printf("[%s] Request ID: %s\n", label, id)
	return id
}

// doSync wraps the synchronous SyncDo pattern used by most examples:
// generate a request ID, call SyncDo, and handle errors uniformly.
// The caller only needs to supply a label (for log prefix) and a closure that
// performs the actual SyncDo call, keeping request construction visible.
//
// Usage:
//
//	doSync("OTO", func(requestID string) (any, error) {
//	    return service.SyncDo(requestID, request)
//	})
func doSync(label string, fn func(requestID string) (any, error)) {
	requestID := newRequestID(label)
	resp, err := fn(requestID)
	if err != nil {
		fmt.Printf("[%s] Error: %v\n", label, err)
		return
	}
	fmt.Printf("[%s] Response: %+v\n", label, resp)
}

// doAsync wraps the async Do + channel-listen pattern used by examples that
// need to wait for a WebSocket push response.  It sends the request, then
// blocks until a response or error arrives (or timeout), so callers don't
// have to write goroutines and time.Sleep by hand.
//
// Usage:
//
//	doAsync("OCO", func(requestID string) (readCh <-chan []byte, errCh <-chan error, err error) {
//	    err = service.Do(requestID, request)
//	    return service.GetReadChannel(), service.GetReadErrorChannel(), err
//	})
func doAsync(label string, timeout time.Duration, fn func(requestID string) (<-chan []byte, <-chan error, error)) {
	requestID := newRequestID(label)
	readCh, errCh, err := fn(requestID)
	if err != nil {
		fmt.Printf("[%s] Error sending request: %v\n", label, err)
		return
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case resp := <-readCh:
		fmt.Printf("[%s] Response: %s\n", label, string(resp))
	case e := <-errCh:
		fmt.Printf("[%s] Error: %v\n", label, e)
	case <-timer.C:
		fmt.Printf("[%s] Timed out waiting for response\n", label)
	}
}

// ---------------------------------------------------------------------------
// Example functions – each one focuses on its own request construction
// ---------------------------------------------------------------------------

// OrderListPlaceOCO demonstrates creating an OCO order using WebSocket API.
// This example uses the async pattern because OCO responses arrive via push.
func OrderListPlaceOCO() {
	client := newClient("OCO")
	if client == nil {
		return
	}

	service, err := client.NewOrderListCreateWsService()
	if err != nil {
		fmt.Printf("[OCO] Error creating service: %v\n", err)
		return
	}

	// --- OCO-specific request ------------------------------------------------
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
	// --------------------------------------------------------------------------

	doAsync("OCO", 5*time.Second, func(requestID string) (<-chan []byte, <-chan error, error) {
		err := service.Do(requestID, request)
		return service.GetReadChannel(), service.GetReadErrorChannel(), err
	})

	service.ReceiveAllDataBeforeStop(2 * time.Second)
}

// OrderListPlaceOTO demonstrates creating an OTO order using WebSocket API.
// Uses the synchronous request/response pattern.
func OrderListPlaceOTO() {
	client := newClient("OTO")
	if client == nil {
		return
	}

	service, err := client.NewOrderListPlaceOtoWsService()
	if err != nil {
		fmt.Printf("[OTO] Error creating service: %v\n", err)
		return
	}

	// --- OTO-specific request ------------------------------------------------
	request := binance.NewOrderListPlaceOtoWsRequest().
		Symbol("BTCUSDT").
		WorkingType(binance.OrderTypeLimit).
		WorkingSide(binance.SideTypeBuy).
		WorkingPrice("30000").
		WorkingQuantity("0.001").
		WorkingTimeInForce(binance.TimeInForceTypeGTC).
		PendingType(binance.OrderTypeLimit).
		PendingSide(binance.SideTypeSell).
		PendingTimeInForce(binance.TimeInForceTypeGTC).
		PendingPrice("32000").
		PendingQuantity("0.001")
	// --------------------------------------------------------------------------

	doSync("OTO", func(requestID string) (any, error) {
		return service.SyncDo(requestID, request)
	})
}

// OrderListPlaceOTOCO demonstrates creating an OTOCO order using WebSocket API.
// Uses the synchronous request/response pattern.
func OrderListPlaceOTOCO() {
	client := newClient("OTOCO")
	if client == nil {
		return
	}

	service, err := client.NewOrderListPlaceOtocoWsService()
	if err != nil {
		fmt.Printf("[OTOCO] Error creating service: %v\n", err)
		return
	}

	// --- OTOCO-specific request ----------------------------------------------
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
	// --------------------------------------------------------------------------

	doSync("OTOCO", func(requestID string) (any, error) {
		return service.SyncDo(requestID, request)
	})
}

// OrderListCancel demonstrates canceling an order list using WebSocket API.
// Uses the synchronous request/response pattern.
func OrderListCancel() {
	client := newClient("Cancel")
	if client == nil {
		return
	}

	service, err := client.NewOrderListCancelWsService()
	if err != nil {
		fmt.Printf("[Cancel] Error creating service: %v\n", err)
		return
	}

	// --- Cancel-specific request ---------------------------------------------
	request := binance.NewOrderListCancelWsRequest().
		Symbol("BTCUSDT").
		OrderListID(123456789) // Replace with actual order list ID
	// --------------------------------------------------------------------------

	doSync("Cancel", func(requestID string) (any, error) {
		return service.SyncDo(requestID, request)
	})
}

// SorOrderPlace demonstrates placing a SOR order using WebSocket API.
// Uses the synchronous pattern and inspects the SOR-specific result fields.
func SorOrderPlace() {
	client := newClient("SOR-Place")
	if client == nil {
		return
	}

	service, err := client.NewSorOrderPlaceWsService()
	if err != nil {
		fmt.Printf("[SOR-Place] Error creating service: %v\n", err)
		return
	}

	// --- SOR Place-specific request (ETHUSDT has SOR support) ----------------
	request := binance.NewSorOrderPlaceWsRequest().
		Symbol("ETHUSDT").
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity("0.1").
		Price("2000").
		TimeInForce(binance.TimeInForceTypeGTC).
		NewClientOrderID("sBI1KM6nNtOfj5tccZSKly").
		NewOrderRespType(binance.NewOrderRespTypeFULL)
	// --------------------------------------------------------------------------

	requestID := newRequestID("SOR-Place")
	response, err := service.SyncDo(requestID, request)
	if err != nil {
		fmt.Printf("[SOR-Place] Error: %v\n", err)
		return
	}
	fmt.Printf("[SOR-Place] Response: %+v\n", response)

	// SOR-specific: check whether Smart Order Routing was used
	if len(response.Result) > 0 {
		fmt.Printf("[SOR-Place] Used SOR: %v\n", response.Result[0].UsedSor)
	} else {
		fmt.Printf("[SOR-Place] No order results returned\n")
	}
}

// SorOrderTest demonstrates testing a SOR order (dry-run) using WebSocket API.
// Uses the synchronous pattern and prints commission rate details.
func SorOrderTest() {
	client := newClient("SOR-Test")
	if client == nil {
		return
	}

	service, err := client.NewSorOrderTestWsService()
	if err != nil {
		fmt.Printf("[SOR-Test] Error creating service: %v\n", err)
		return
	}

	// --- SOR Test-specific request (ETHUSDT for SOR support) -----------------
	request := binance.NewSorOrderTestWsRequest().
		Symbol("ETHUSDT").
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		Quantity("0.1").
		Price("2000").
		TimeInForce(binance.TimeInForceTypeGTC).
		ComputeCommissionRates(true)
	// --------------------------------------------------------------------------

	requestID := newRequestID("SOR-Test")
	response, err := service.SyncDo(requestID, request)
	if err != nil {
		fmt.Printf("[SOR-Test] Error: %v\n", err)
		return
	}
	fmt.Printf("[SOR-Test] Response: %+v\n", response)

	// SOR Test-specific: show commission rates when available
	if response.Result.StandardCommissionForOrder != nil {
		fmt.Printf("[SOR-Test] Standard Commission - Maker: %s, Taker: %s\n",
			response.Result.StandardCommissionForOrder.Maker,
			response.Result.StandardCommissionForOrder.Taker)
	}
}

// RunOrderListExamples runs all order list examples sequentially.
func RunOrderListExamples() {
	fmt.Println("=== Binance Order List WebSocket API Examples ===")

	fmt.Println("\n1. OCO Order (Current)")
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
