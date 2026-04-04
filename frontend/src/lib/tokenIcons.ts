const icons: Record<string, string> = {
  USDC:   "/tokens/usdc.png",
  WETH:   "/tokens/weth.png",
  UNI:    "/tokens/uni.png",
  LINK:   "/tokens/link.png",
  DAI:    "/tokens/dai.png",
  USDT:   "/tokens/usdt.png",
  WBTC:   "/tokens/wbtc.png",
  wstETH: "/tokens/wsteth.png",
  POL:    "/tokens/pol.png",
};

export function tokenIcon(symbol: string): string | undefined {
  return icons[symbol];
}
