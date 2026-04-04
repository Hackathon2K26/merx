package config

import "fmt"

// Chain holds network-specific configuration for a supported blockchain.
type Chain struct {
	Name         string
	ChainID      int
	AddressUSDC  string
	AddressWETH  string
	DecimalsUSDC int
	// Tokens maps a human-readable symbol to its contract address on this chain.
	Tokens map[string]string
}

// Ethereum Mainnet token addresses.
const (
	MainnetUSDC = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	MainnetWETH = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"
	MainnetUSDT = "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	MainnetDAI  = "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	MainnetWBTC = "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599"
)

var (
	// EthereumMainnet is used for forked tests only — not exposed in the CLI.
	EthereumMainnet = &Chain{
		Name:         "ethereum-mainnet",
		ChainID:      1,
		AddressUSDC:  MainnetUSDC,
		AddressWETH:  MainnetWETH,
		DecimalsUSDC: 6,
		Tokens: map[string]string{
			"USDC": MainnetUSDC,
			"WETH": MainnetWETH,
			"USDT": MainnetUSDT,
			"DAI":  MainnetDAI,
			"WBTC": MainnetWBTC,
		},
	}

	EthereumSepolia = &Chain{
		Name:         "ethereum-sepolia",
		ChainID:      11155111,
		AddressUSDC:  "0x94a9d9ac8a22534e3faca9f4e7f2e2cf85d5e4c8",
		AddressWETH:  "0xfff9976782d46cc05630d1f6ebab18b2324d6b14",
		DecimalsUSDC: 6,
	}

	BaseSepolia = &Chain{
		Name:         "base-sepolia",
		ChainID:      84532,
		AddressUSDC:  "0x036cbd53842c5426634e7929541ec2318f3dcf7e",
		AddressWETH:  "0x4200000000000000000000000000000000000006",
		DecimalsUSDC: 6,
	}

	UnichainSepolia = &Chain{
		Name:         "unichain-sepolia",
		ChainID:      1301,
		AddressUSDC:  "0x31d0220469e10c4e71834a79b1f276d740d3768f",
		AddressWETH:  "0x4200000000000000000000000000000000000006",
		DecimalsUSDC: 6,
	}

	// supportedChains are the chains exposed in the CLI (testnets only).
	supportedChains = map[string]*Chain{
		"ethereum-sepolia": EthereumSepolia,
		"base-sepolia":     BaseSepolia,
		"unichain-sepolia": UnichainSepolia,
	}
)

// ChainByName returns a chain configuration by its name.
func ChainByName(name string) (*Chain, error) {
	chain, ok := supportedChains[name]
	if !ok {
		return nil, fmt.Errorf("unsupported chain: %s", name)
	}
	return chain, nil
}

// SupportedChainNames returns all supported chain names.
func SupportedChainNames() []string {
	names := make([]string, 0, len(supportedChains))
	for name := range supportedChains {
		names = append(names, name)
	}
	return names
}
