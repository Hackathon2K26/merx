package uniswap

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"merx/internal/config"
)

// GetQuote fetches a swap quote from the Uniswap Trading API.
func (c *Client) GetQuote(ctx context.Context, req *QuoteRequest) (*QuoteResponse, error) {
	var resp QuoteResponse
	if err := c.do(ctx, http.MethodPost, "/quote", req, &resp); err != nil {
		return nil, fmt.Errorf("get quote: %w", err)
	}
	return &resp, nil
}

// GetPriceUSDC returns a quote for swapping the given USDC amount to tokenOut on the specified chain.
// usdcAmount is in human-readable units (e.g., 100.0 for $100).
func (c *Client) GetPriceUSDC(ctx context.Context, usdcAmount float64, tokenOut string, chain *config.Chain, swapper string) (*QuoteResponse, error) {
	// Convert USDC amount to base units (6 decimals)
	amountBig := new(big.Int)
	amountBig.SetInt64(int64(usdcAmount * 1e6))

	req := &QuoteRequest{
		Type:            "EXACT_INPUT",
		Amount:          amountBig.String(),
		TokenIn:         chain.AddressUSDC,
		TokenOut:        tokenOut,
		TokenInChainId:  chain.ChainID,
		TokenOutChainId: chain.ChainID,
		Swapper:         swapper,
	}

	return c.GetQuote(ctx, req)
}
