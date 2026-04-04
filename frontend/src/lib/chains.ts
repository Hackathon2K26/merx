import { type Chain } from "viem";
import { sepolia, baseSepolia } from "viem/chains";

export const unichainSepolia: Chain = {
  id: 1301,
  name: "Unichain Sepolia",
  nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
  rpcUrls: {
    default: { http: ["https://sepolia.unichain.org"] },
  },
  testnet: true,
};

export const supportedChains = [sepolia, baseSepolia, unichainSepolia] as const;

export function chainName(chainId: number): string {
  const chain = supportedChains.find((c) => c.id === chainId);
  return chain?.name ?? `Chain ${chainId}`;
}
