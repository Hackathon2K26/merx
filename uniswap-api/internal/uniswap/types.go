package uniswap

import (
	"encoding/json"
	"fmt"
)

// --- Quote ---

type QuoteRequest struct {
	Type              string   `json:"type"`
	Amount            string   `json:"amount"`
	TokenIn           string   `json:"tokenIn"`
	TokenOut          string   `json:"tokenOut"`
	TokenInChainId    int      `json:"tokenInChainId"`
	TokenOutChainId   int      `json:"tokenOutChainId"`
	Swapper           string   `json:"swapper"`
	SlippageTolerance *float64 `json:"slippageTolerance,omitempty"`
	Urgency           string   `json:"urgency,omitempty"`
}

type QuoteResponse struct {
	RequestID         string              `json:"requestId"`
	RawQuote          json.RawMessage     `json:"quote"`
	Quote             ClassicQuote        `json:"-"` // parsed from RawQuote
	Routing           string              `json:"routing"`
	PermitData        *PermitData         `json:"permitData"`
	PermitTransaction *TransactionRequest `json:"permitTransaction,omitempty"`
}

type ClassicQuote struct {
	Input          TokenAmount `json:"input"`
	Output         TokenAmount `json:"output"`
	Swapper        string      `json:"swapper"`
	ChainID        int         `json:"chainId"`
	Slippage       float64     `json:"slippage"`
	TradeType      string      `json:"tradeType"`
	GasFee         string      `json:"gasFee"`
	GasFeeUSD      string      `json:"gasFeeUSD"`
	GasFeeQuote    string      `json:"gasFeeQuote"`
	GasUseEstimate string      `json:"gasUseEstimate"`
	RouteString    string      `json:"routeString"`
	QuoteID        string      `json:"quoteId"`
	PriceImpact    float64     `json:"priceImpact"`
	BlockNumber    string      `json:"blockNumber"`
	GasPrice       string      `json:"gasPrice"`
}

type TokenAmount struct {
	Token  string `json:"token"`
	Amount string `json:"amount"`
}

// --- Approval ---

type ApprovalRequest struct {
	WalletAddress string `json:"walletAddress"`
	Token         string `json:"token"`
	Amount        string `json:"amount"`
	ChainID       int    `json:"chainId"`
}

type ApprovalResponse struct {
	RequestID    string              `json:"requestId"`
	Approval     *TransactionRequest `json:"approval"`
	Cancel       *TransactionRequest `json:"cancel"`
	GasFee       string              `json:"gasFee"`
	CancelGasFee string              `json:"cancelGasFee"`
}

// --- Swap ---

// SwapRequest wraps the raw quote from the /quote response along with
// optional Permit2 signature data. Quote is passed through as raw JSON
// to preserve all fields the API expects.
type SwapRequest struct {
	Quote      json.RawMessage `json:"quote"`
	Signature  string          `json:"signature,omitempty"`
	PermitData *PermitData     `json:"permitData,omitempty"`
}

type SwapResponse struct {
	RequestID string             `json:"requestId"`
	Swap      TransactionRequest `json:"swap"`
	GasFee    string             `json:"gasFee"`
}

// --- Order (UniswapX gasless) ---

type OrderRequest struct {
	Signature string      `json:"signature"`
	Quote     any `json:"quote"`
	Routing   string      `json:"routing"`
}

type OrderResponse struct {
	RequestID   string `json:"requestId"`
	OrderID     string `json:"orderId"`
	OrderStatus string `json:"orderStatus"`
}

// --- Shared ---

type TransactionRequest struct {
	To                   string `json:"to"`
	From                 string `json:"from"`
	Data                 string `json:"data"`
	Value                string `json:"value"`
	ChainID              int    `json:"chainId"`
	GasLimit             string `json:"gasLimit,omitempty"`
	MaxFeePerGas         string `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas,omitempty"`
	GasPrice             string `json:"gasPrice,omitempty"`
}

type PermitData struct {
	Domain map[string]any `json:"domain"`
	Values map[string]any `json:"values"`
	Types  map[string]any `json:"types"`
}

// APIError represents an error response from the Uniswap API.
type APIError struct {
	ErrorCode  string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("uniswap API error %d: %s - %s", e.StatusCode, e.ErrorCode, e.Message)
}
