package uniswap

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/RomainLafont/merx/uniswap-api/config"
)

// GetQuote fetches a swap quote from the Uniswap Trading API.
// It preserves the raw quote JSON for pass-through to /swap, and also
// parses it into the ClassicQuote struct for convenient field access.
func (c *Client) GetQuote(ctx context.Context, req *QuoteRequest) (*QuoteResponse, error) {
	var resp QuoteResponse
	if err := c.do(ctx, http.MethodPost, "/quote", req, &resp); err != nil {
		return nil, fmt.Errorf("get quote: %w", err)
	}
	if err := json.Unmarshal(resp.RawQuote, &resp.Quote); err != nil {
		return nil, fmt.Errorf("parse quote: %w", err)
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

// GetPriceInUSDC returns a quote for swapping the given token amount into USDC on the specified chain.
// amount is the raw amount in the token's smallest unit (e.g., wei for WETH).
func (c *Client) GetPriceInUSDC(ctx context.Context, amount string, tokenIn string, chain *config.Chain, swapper string) (*QuoteResponse, error) {
	req := &QuoteRequest{
		Type:            "EXACT_INPUT",
		Amount:          amount,
		TokenIn:         tokenIn,
		TokenOut:        chain.AddressUSDC,
		TokenInChainId:  chain.ChainID,
		TokenOutChainId: chain.ChainID,
		Swapper:         swapper,
	}

	return c.GetQuote(ctx, req)
}
