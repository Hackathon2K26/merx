import { http, createConfig, fallback } from "wagmi";
import { sepolia, baseSepolia } from "viem/chains";
import { injected } from "wagmi/connectors";
import { unichainSepolia } from "./chains";

export const config = createConfig({
  chains: [sepolia, baseSepolia, unichainSepolia],
  connectors: [injected()],
  transports: {
    [sepolia.id]: fallback([
      http("https://ethereum-sepolia-rpc.publicnode.com"),
      http("https://rpc.sepolia.org"),
      http("https://rpc2.sepolia.org"),
    ]),
    [baseSepolia.id]: fallback([
      http("https://base-sepolia-rpc.publicnode.com"),
      http("https://sepolia.base.org"),
    ]),
    [unichainSepolia.id]: http("https://sepolia.unichain.org"),
  },
});
