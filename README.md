# Merx

Merx is a chain-abstracted ebook shop built around Circle USDC on testnet. Customers can pay from 18 CCTP V2 source chains, settlement lands on Arc, and the merchant can refund to any supported chain or move idle USDC into Compound V3 on Ethereum Sepolia.

The current checkout flow no longer relies on `ShopPaymaster` or a gasless permit contract. Direct USDC checkout now uses Circle's CCTP Forwarding Service with a customer-side `approve + depositForBurnWithHook` flow. On Uniswap-enabled chains, customers can also swap supported tokens into USDC before bridging.

## Current Snapshot

- The live payment path uses Circle CCTP V2 + Forwarding Service, not `EIP-2612 + ShopPaymaster`.
- Chain and token metadata are registry-driven via [`registry.yaml`](/home/ari/Perso/git/merx/registry.yaml).
- The customer UI supports direct USDC checkout on all supported chains, plus token-to-USDC swaps on Ethereum Sepolia, Base Sepolia, and Unichain Sepolia.
- The merchant dashboard shows cross-chain USDC balances, native gas balances, invoice states, refunds, explorer links, and Compound V3 treasury actions.
- Ebook download access is tied to invoice state: paid invoices can download, refunded invoices cannot.

## Architecture

### Runtime Architecture

```text
┌────────────────────────────────────────────────────────────────────┐
│ Frontend (React + Vite + TypeScript)                              │
│                                                                    │
│ Shop        Checkout                          Merchant dashboard    │
│ /           /checkout/:productId              /dashboard            │
│                                                                    │
│ - browse ebooks                              - balances by chain    │
│ - connect wallet                             - invoices + refunds   │
│ - choose source chain                        - Compound supply      │
│ - pay in USDC or swap first                  - Compound withdraw    │
└───────────────────────────────┬────────────────────────────────────┘
                                │ /api/*
┌───────────────────────────────▼────────────────────────────────────┐
│ Go API server                                                     │
│                                                                    │
│ - loads chain/token registry                                       │
│ - builds CCTP calldata                                             │
│ - records invoices in invoices.json                                │
│ - polls Circle attestation / forwarding status                     │
│ - proxies Uniswap Trading API                                      │
│ - orchestrates refund and treasury flows                           │
└───────────────┬───────────────────────┬────────────────────────────┘
                │                       │
                │                       │
     ┌──────────▼──────────┐  ┌────────▼────────┐
     │ Circle CCTP V2      │  │ Uniswap API     │
     │ + Forwarding        │  │                 │
     │                     │  │ quotes + swaps  │
     └──────────┬──────────┘  └─────────────────┘
                │
     ┌──────────▼──────────┐
     │ Arc Testnet         │
     │ settlement hub      │
     └──────────┬──────────┘
                │
     ┌──────────▼──────────┐
     │ Compound V3         │
     │ Ethereum Sepolia    │
     │ relay + supply      │
     └─────────────────────┘
```

### Customer Flow

#### 1. Browse and select a product

1. The customer lands on `/` and browses the ebook catalog.
2. Each product links to `/checkout/:productId`.
3. Once connected, the UI loads the supported chain list from [`registry.yaml`](/home/ari/Perso/git/merx/registry.yaml) through `GET /api/chains`.

#### 2. Pay with USDC directly

1. The customer selects a source chain.
2. The frontend calls `GET /api/pay-tx` to fetch the `depositForBurnWithHook` calldata and the required approval target.
3. The wallet signs `approve(TokenMessengerV2, amount)`.
4. The wallet signs `depositForBurnWithHook(...)`.
5. Circle forwards the burn to Arc and mints USDC directly to the merchant wallet.
6. The frontend reports the source transaction with `POST /api/pay`.
7. The backend creates an invoice, polls forwarding completion, and the ebook becomes downloadable.

#### 3. Pay with another token on Uniswap-enabled chains

1. The customer picks a supported token with a positive wallet balance.
2. The frontend fetches a quote from `POST /api/uniswap/quote`.
3. If needed, the wallet approves Permit2 and signs the Uniswap permit payload.
4. The wallet executes the swap transaction.
5. Once the wallet holds USDC, the checkout continues with the same CCTP payment flow described above.

### Merchant Flow

#### 1. Operate the store from Arc

1. The merchant connects the configured merchant wallet and unlocks `/dashboard`.
2. The dashboard polls `GET /api/merchant/balances` and `GET /api/invoices`.
3. The merchant sees:
   - USDC balances on every supported source chain plus Arc
   - native gas balances
   - Compound V3 position and APY
   - invoice status and transaction links

#### 2. Refund a customer to any supported chain

1. The merchant opens an invoice and chooses the refund destination chain.
2. The backend estimates the forwarding fee with Circle's fee API.
3. The backend sends `depositForBurnWithHook` from Arc for `amount + fee`.
4. Circle forwards the refund to the customer's chosen chain.
5. The invoice transitions through `refunding -> refunded`, and ebook download access is blocked.

#### 3. Move idle treasury into Compound V3

1. `POST /api/supply` burns Arc USDC toward Ethereum Sepolia.
2. [`CompoundDepositor.sol`](/home/ari/Perso/git/merx/contracts/src/CompoundDepositor.sol) relays the CCTP message and supplies into Compound in one call.
3. `POST /api/withdraw` withdraws from Compound, approves `TokenMessengerV2`, and bridges USDC back to Arc.

### Repository Architecture

```text
cmd/server/
  main.go                    Go API server: CCTP flows, invoices, refunds, treasury, Uniswap proxy

frontend/
  src/pages/                 Shop, checkout, merchant dashboard
  src/components/            Payment flow, selectors, wallet UI, layout
  src/lib/                   API client, chain metadata, wagmi config, product catalog, formatting
  public/                    Chain icons, token icons, brand assets

contracts/
  src/CompoundDepositor.sol  CCTP relay + Compound V3 supply helper on Ethereum Sepolia
  script/                    Foundry deployment script
  test/                      Foundry tests

uniswap-api/
  config/                    Uniswap config loading and chain helpers
  uniswap/                   Quote, approval, swap client code
  cmd/merx/                  CLI for local testing of the Uniswap integration

ebooks/
  PDF assets served after a successful purchase

registry.yaml
  Source-of-truth registry for supported chains, tokens, RPCs, explorers, and CCTP domains

params.go
  Shared protocol addresses, Arc constants, RPC URLs, and deployed contract addresses
```

## Technical Stack

| Layer | Technology | Why it is used | Why it is useful here |
|-------|------------|----------------|-----------------------|
| Frontend | React 19 + TypeScript + Vite | Fast SPA development with strong typing and simple local DX | The checkout and dashboard have many async wallet states, and the codebase benefits from typed UI flows and fast iteration |
| Wallet / RPC | wagmi + viem | Multi-chain wallet connection, reads, writes, chain switching, typed calldata | The app needs reliable wallet UX across many testnets and direct contract interaction from the browser |
| Data fetching | TanStack React Query | Polling, cache invalidation, loading/error state handling | Balances, invoices, and payment status refresh continuously without custom polling code everywhere |
| Styling | Tailwind CSS v4 | Fast component styling with small surface area | Useful for a compact dashboard + checkout UI without adding a heavy design system |
| Backend | Go 1.25 + `net/http` | Single binary, low overhead, easy concurrency | The backend spends most of its time building calldata, calling RPCs, proxying APIs, and polling Circle status endpoints |
| EVM integration | `go-ethereum` | ABI packing, tx broadcasting, contract calls, RPC clients | Central to payment, refund, settlement, and treasury actions |
| Cross-chain settlement | Circle CCTP V2 + Forwarding Service | Native burn/mint USDC bridging without wrapped liquidity | Removes the need for custom source-chain payment contracts and keeps settlement in canonical USDC |
| Settlement hub | Arc Testnet | Fast settlement hub for merchant funds | Gives the merchant one main treasury location before refunding or deploying capital |
| Swap layer | Uniswap Trading API + Permit2 | Quotes and swap transaction building for non-USDC assets | Lets the customer pay from a broader token set without Merx implementing routing itself |
| Treasury | Compound V3 | Yield venue for idle USDC | Simple lending flow with on-chain APY data and reversible treasury operations |
| Smart contracts | Solidity + Foundry | Contract development, testing, deployment scripts | Only one focused helper contract remains in the live architecture: `CompoundDepositor` |
| Configuration | YAML registry (`registry.yaml`) | Central chain/token catalog | Keeps payment, refund, explorer, and UI metadata aligned across the app |

## Supported Networks

Merx currently supports 18 CCTP-enabled source testnets:

- Ethereum Sepolia
- Avalanche Fuji
- OP Sepolia
- Arbitrum Sepolia
- Base Sepolia
- Polygon Amoy
- Unichain Sepolia
- Sonic Blaze
- Worldchain Sepolia
- Sei Atlantic
- Linea Sepolia
- Codex Testnet
- Monad Testnet
- HyperEVM Testnet
- Ink Sepolia
- Plume Testnet
- EDGE Testnet
- Morph Hoodi

Settlement hub:

- Arc Testnet (`chainId 5042002`, `domain 26`)

Uniswap swap flow is currently enabled on:

- Ethereum Sepolia
- Base Sepolia
- Unichain Sepolia

## API Surface

| Method | Endpoint | Role |
|--------|----------|------|
| GET | `/api/config` | Return the active merchant address exposed by the backend |
| GET | `/api/chains` | Return the chain/token registry consumed by the frontend |
| GET | `/api/pay-tx` | Build direct checkout calldata for `depositForBurnWithHook` |
| POST | `/api/pay` | Record a customer payment and start settlement tracking |
| GET | `/api/balances` | Read the merchant USDC balance on Arc |
| GET | `/api/merchant/balances` | Aggregate balances, native gas, Compound balance, and APY |
| POST | `/api/refund` | Start a refund from Arc to any supported destination chain |
| POST | `/api/supply` | Bridge Arc USDC to Ethereum Sepolia and supply to Compound |
| POST | `/api/withdraw` | Withdraw from Compound and bridge the funds back to Arc |
| POST | `/api/invoices` | Create a pending invoice manually |
| GET | `/api/invoices` | List invoices |
| GET | `/api/invoices/{id}` | Read a single invoice |
| GET | `/api/ebooks/{invoiceId}` | Download the purchased PDF if the invoice is still valid |
| POST | `/api/uniswap/quote` | Get a swap quote |
| POST | `/api/uniswap/approval` | Check swap approval requirements |
| POST | `/api/uniswap/swap` | Build a swap transaction |

## Quick Start

### Prerequisites

- Go `1.25+`
- Node.js `20+`
- npm
- Foundry, if you want to run Solidity tests or deploy the helper contract
- A private key funded on the testnets you want to operate from

### 1. Configure Uniswap (optional, required for swap checkout)

```bash
cp uniswap-api/config.yaml.example uniswap-api/config.yaml
```

Then edit `uniswap-api/config.yaml`:

```yaml
uniswap_api_key: "your-api-key"
swapper_address: "0xYourWallet"
```

If this file is missing, the backend still starts, but the `/api/uniswap/*` routes are effectively disabled and only direct USDC checkout works.

### 2. Start the backend

```bash
PRIVATE_KEY=0x... go run ./cmd/server
```

Optional flags:

```bash
go run ./cmd/server --port 8080 --registry registry.yaml --uniswap-config uniswap-api/config.yaml
```

The backend persists invoice state to `invoices.json` at the repository root.

### 3. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

Open:

- Shop: `http://localhost:5173/`
- Merchant dashboard: `http://localhost:5173/dashboard`

The Vite dev server proxies `/api/*` to `http://localhost:8080`.

## Testing

```bash
# Go tests
go test ./...

# Solidity tests
cd contracts && forge test -vvv

# Frontend production build
cd frontend && npm run build
```

## Deployed Components

| Component | Network | Address / note |
|-----------|---------|----------------|
| `CompoundDepositor` | Ethereum Sepolia | `0x832705f381957C8218d7ae8B20A10d510B5AFB75` |
| `TokenMessengerV2` | CCTP testnets | `0x8FE6B999Dc680CcFDD5Bf7EB0974218be2542DAA` |
| `MessageTransmitter` | CCTP testnets | `0xE737e5cEBEEBa77EFE34D4aa090756590b1CE275` |

There is no source-chain `ShopPaymaster` in the current live checkout flow.

## Invoice Lifecycle

```text
paid -> bridging -> attesting -> settled
                             |
                             -> refunding -> refunded
```

The dashboard refreshes invoice and balance data every 10 seconds.
