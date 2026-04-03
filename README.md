# Merx

Modular trading toolkit. Each subdirectory is an independent Go module.

## Modules

### [uniswap-api](./uniswap-api/)

Go client for the [Uniswap Trading API](https://api-docs.uniswap.org/introduction). Get token prices in USDC and execute swaps on supported testnets (Ethereum Sepolia, Base Sepolia, Unichain Sepolia).

See [uniswap-api/README.md](./uniswap-api/README.md) for setup and usage.

## Progress

- [x] Uniswap API client (quote, swap, approval)
- [x] Multi-chain testnet support
- [x] Unit + integration tests
- [x] CLI tool
- [ ] Transaction signing and broadcasting
- [ ] UniswapX gasless order flow
