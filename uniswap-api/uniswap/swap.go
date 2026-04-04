package uniswap

import (
	"context"
	"fmt"
	"net/http"
)

// CreateSwap converts a CLASSIC quote into an unsigned transaction.
// Use this for routing types: CLASSIC, WRAP, UNWRAP, BRIDGE.
func (c *Client) CreateSwap(ctx context.Context, req *SwapRequest) (*SwapResponse, error) {
	var resp SwapResponse
	if err := c.do(ctx, http.MethodPost, "/swap", req, &resp); err != nil {
		return nil, fmt.Errorf("create swap: %w", err)
	}
	return &resp, nil
}

// CreateOrder submits a gasless UniswapX order.
// Use this for routing types: DUTCH_V2, DUTCH_V3, PRIORITY.
func (c *Client) CreateOrder(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	var resp OrderResponse
	if err := c.do(ctx, http.MethodPost, "/order", req, &resp); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	return &resp, nil
}
