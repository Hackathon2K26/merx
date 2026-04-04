# Merx

On-chain ebook shop where customers pay in USDC — or swap from any token via Uniswap in a single transaction using Permit2. The shop can refund cross-chain via Circle Gateway.

## Architecture

```
Browser (React SPA :5173)  →  Go API server (:8080)  →  Uniswap Trading API
                                                      →  Circle Gateway API
```

### Payment flow (customer → merchant)

**Direct USDC:** Customer sends an ERC-20 `transfer` to the merchant address.

**Swap from any token (single TX):**
1. Customer picks a token (WETH, UNI, etc.) and gets a quote
2. Signs a Permit2 message (gasless, off-chain)
3. Confirms one swap transaction — the Universal Router handles permit + swap atomically, sending USDC directly to the merchant

### Refund flow (merchant → customer)

```
POST /api/refund  →  allocate Gateway balances  →  sign burn intents
  →  Gateway Forwarding Service  →  customer receives USDC on target chain
```

## Project structure

```
registry.yaml              Token registry (chains, tokens, addresses)
gateway/                   Gateway client, types, EIP-712 signer, tests
uniswap-api/config/        Chain config, YAML loader
uniswap-api/uniswap/       Uniswap Trading API client (quote, swap, approval)
cmd/server/                API server (invoices, uniswap proxy, gateway, refund)
cmd/refund/                Admin refund CLI
frontend/                  React + Vite + TypeScript webshop
  src/pages/               ShopPage (catalog), CheckoutPage (payment)
  src/components/          PaymentFlow, ProductCard, Layout, ConnectWallet
  src/lib/                 API client, wagmi config, products, formatting
```

## Quick start

```bash
# 1. Start the API server (requires uniswap-api/config.yaml with API key)
go run cmd/server/main.go

# 2. Start the frontend
cd frontend && npm install && npm run dev
```

Open http://localhost:5173 — browse ebooks, connect MetaMask, pick a chain and token, pay.

### Configuration

**`uniswap-api/config.yaml`** — Uniswap Trading API key:
```yaml
uniswap_api_key: "your-key-from-developers.uniswap.org"
swapper_address: "0xYourAddress"
```

**`registry.yaml`** — Supported chains and tokens:
```yaml
chains:
  - name: Ethereum Sepolia
    chainId: 11155111
    tokens:
      - symbol: USDC
        decimals: 6
        address: "0x94a9d9ac8a22534e3faca9f4e7f2e2cf85d5e4c8"
      - symbol: WETH
        decimals: 18
        address: "0xfff9976782d46cc05630d1f6ebab18b2324d6b14"
```

## API server

```bash
go run cmd/server/main.go
go run cmd/server/main.go --port 3001 --registry registry.yaml
```

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/chains` | Supported chains and tokens (from registry.yaml) |
| POST | `/api/invoices` | Create invoice `{merchantAddress, amount, chainId, description}` |
| GET | `/api/invoices` | List invoices (optional `?merchant=0x...`) |
| GET | `/api/invoices/{id}` | Get invoice by ID |
| POST | `/api/invoices/{id}/pay` | Mark invoice paid `{txHash}` |
| POST | `/api/uniswap/quote` | Get swap quote `{tokenIn, tokenInChainId, amount, swapper}` |
| POST | `/api/uniswap/approval` | Check token approval |
| POST | `/api/uniswap/swap` | Build swap TX `{quote, signature, permitData}` |
| GET | `/api/info` | Gateway domains and contracts |
| GET | `/api/balances` | Shop's Gateway USDC balances |
| POST | `/api/refund` | Start cross-chain refund `{to, chain, amount}` |
| GET | `/api/refund/{id}` | Poll refund status |

## Testing

```bash
# Unit tests
go test ./...

# Frontend type-check + build
cd frontend && npx tsc --noEmit && npx vite build

# Integration tests (hits live APIs)
INTEGRATION=1 go test ./gateway/ -run TestSmoke -v
INTEGRATION=1 go test ./cmd/server/ -v

# On-chain smoke tests (requires funded wallets)
SMOKE=1 go test ./gateway/ -run TestSelfmintFull -v -timeout 35m
```

## Supported chains (testnet)

| Chain | Chain ID | USDC | Swap support |
|-------|----------|------|--------------|
| Ethereum Sepolia | 11155111 | `0x94a9...e4c8` | WETH, UNI |
| Base Sepolia | 84532 | `0x036c...f7e` | WETH (limited liquidity) |
| Unichain Sepolia | 1301 | `0x31d0...68f` | WETH |
