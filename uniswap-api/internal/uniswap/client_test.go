package uniswap

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"merx/internal/config"
)

func TestGetQuote_MockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/quote" {
			t.Errorf("expected /quote, got %s", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("missing or wrong x-api-key header")
		}

		var req QuoteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.TokenIn != config.EthereumSepolia.AddressUSDC {
			t.Errorf("expected tokenIn=%s, got %s", config.EthereumSepolia.AddressUSDC, req.TokenIn)
		}
		if req.Type != "EXACT_INPUT" {
			t.Errorf("expected type=EXACT_INPUT, got %s", req.Type)
		}
		if req.Amount != "100000000" {
			t.Errorf("expected amount=100000000, got %s", req.Amount)
		}
		if req.TokenInChainId != config.EthereumSepolia.ChainID {
			t.Errorf("expected chainId=%d, got %d", config.EthereumSepolia.ChainID, req.TokenInChainId)
		}

		resp := QuoteResponse{
			RequestID: "test-request-id",
			Routing:   "CLASSIC",
			Quote: ClassicQuote{
				Input:       TokenAmount{Token: config.EthereumSepolia.AddressUSDC, Amount: "100000000"},
				Output:      TokenAmount{Token: config.EthereumSepolia.AddressWETH, Amount: "55000000000000000"},
				ChainID:     config.EthereumSepolia.ChainID,
				TradeType:   "EXACT_INPUT",
				RouteString: "USDC -> WETH",
				QuoteID:     "test-quote-id",
				GasFeeUSD:   "2.50",
				PriceImpact: 0.01,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")

	resp, err := client.GetPriceUSDC(context.Background(), 100.0, config.EthereumSepolia.AddressWETH, config.EthereumSepolia, "0x0000000000000000000000000000000000000001")
	if err != nil {
		t.Fatalf("GetPriceUSDC failed: %v", err)
	}

	if resp.Routing != "CLASSIC" {
		t.Errorf("expected routing=CLASSIC, got %s", resp.Routing)
	}
	if resp.Quote.Output.Token != config.EthereumSepolia.AddressWETH {
		t.Errorf("expected output token=%s, got %s", config.EthereumSepolia.AddressWETH, resp.Quote.Output.Token)
	}
	if resp.Quote.Output.Amount != "55000000000000000" {
		t.Errorf("expected output amount=55000000000000000, got %s", resp.Quote.Output.Amount)
	}
	if resp.Quote.GasFeeUSD != "2.50" {
		t.Errorf("expected gasFeeUSD=2.50, got %s", resp.Quote.GasFeeUSD)
	}
}

func TestGetQuote_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "VALIDATION_ERROR",
			"message": "Invalid token address",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	_, err := client.GetQuote(context.Background(), &QuoteRequest{
		Type:            "EXACT_INPUT",
		Amount:          "1000000",
		TokenIn:         "invalid",
		TokenOut:        config.EthereumSepolia.AddressWETH,
		TokenInChainId:  config.EthereumSepolia.ChainID,
		TokenOutChainId: config.EthereumSepolia.ChainID,
		Swapper:         "0x0000000000000000000000000000000000000001",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError in error chain, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", apiErr.StatusCode)
	}
}

func TestGetPriceUSDC_MultipleChains(t *testing.T) {
	chains := []*config.Chain{
		config.EthereumSepolia,
		config.BaseSepolia,
		config.UnichainSepolia,
	}

	for _, chain := range chains {
		t.Run(chain.Name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req QuoteRequest
				json.NewDecoder(r.Body).Decode(&req)

				if req.TokenIn != chain.AddressUSDC {
					t.Errorf("expected tokenIn=%s for chain %s, got %s", chain.AddressUSDC, chain.Name, req.TokenIn)
				}
				if req.TokenInChainId != chain.ChainID {
					t.Errorf("expected chainId=%d for chain %s, got %d", chain.ChainID, chain.Name, req.TokenInChainId)
				}

				resp := QuoteResponse{
					RequestID: "test-" + chain.Name,
					Routing:   "CLASSIC",
					Quote: ClassicQuote{
						Input:  TokenAmount{Token: chain.AddressUSDC, Amount: "50000000"},
						Output: TokenAmount{Token: chain.AddressWETH, Amount: "25000000000000000"},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			client := NewClient(server.URL, "test-key")
			resp, err := client.GetPriceUSDC(context.Background(), 50.0, chain.AddressWETH, chain, "0x0000000000000000000000000000000000000001")
			if err != nil {
				t.Fatalf("GetPriceUSDC failed for chain %s: %v", chain.Name, err)
			}

			if resp.Quote.Output.Amount == "" {
				t.Errorf("expected non-empty output amount for chain %s", chain.Name)
			}
		})
	}
}

// TestGetQuote_Integration hits the real Uniswap API.
// Requires a config.yaml at the project root. Run from the project root:
//
//	go test ./internal/uniswap/ -run TestGetQuote_Integration -v
func TestGetQuote_Integration(t *testing.T) {
	cfg, err := config.Load("../../config.yaml")
	if err != nil {
		t.Skipf("config.yaml not found or invalid, skipping integration test: %v", err)
	}

	client := NewClient(cfg.BaseURL, cfg.APIKey)

	resp, err := client.GetPriceUSDC(context.Background(), 100.0, config.EthereumSepolia.AddressWETH, config.EthereumSepolia, cfg.SwapperAddress)
	if err != nil {
		t.Fatalf("GetPriceUSDC failed: %v", err)
	}

	t.Logf("Routing: %s", resp.Routing)
	t.Logf("Quote ID: %s", resp.Quote.QuoteID)
	t.Logf("Input: %s %s", resp.Quote.Input.Amount, resp.Quote.Input.Token)
	t.Logf("Output: %s %s", resp.Quote.Output.Amount, resp.Quote.Output.Token)
	t.Logf("Gas fee (USD): %s", resp.Quote.GasFeeUSD)
	t.Logf("Price impact: %.4f%%", resp.Quote.PriceImpact)

	if resp.Quote.Output.Amount == "" || resp.Quote.Output.Amount == "0" {
		t.Error("expected non-zero output amount")
	}
}
