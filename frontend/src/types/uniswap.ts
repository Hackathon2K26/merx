export interface QuoteRequest {
  tokenIn: string;
  tokenInChainId: number;
  amount: string; // USDC in base units
  swapper: string;
}

export interface QuoteResponse {
  requestId: string;
  quote: RawQuote;
  routing: string;
  permitData?: PermitData;
}

export interface RawQuote {
  input: TokenAmount;
  output: TokenAmount & { recipient?: string };
  swapper: string;
  chainId: number;
  slippage: number;
  tradeType: string;
  gasFee: string;
  gasFeeUSD: string;
  gasFeeQuote: string;
  gasUseEstimate: string;
  routeString: string;
  quoteId: string;
  priceImpact: number;
  blockNumber: string;
  gasPrice: string;
}

export interface TokenAmount {
  token: string;
  amount: string;
}

export interface PermitData {
  domain: Record<string, unknown>;
  values: Record<string, unknown>;
  types: Record<string, unknown>;
}

export interface TransactionRequest {
  to: string;
  from: string;
  data: string;
  value: string;
  chainId: number;
  gasLimit?: string;
  maxFeePerGas?: string;
  maxPriorityFeePerGas?: string;
  gasPrice?: string;
}

export interface SwapResponse {
  requestId: string;
  swap: TransactionRequest;
  gasFee: string;
}
