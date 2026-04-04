export interface Invoice {
  id: string;
  merchantAddress: string;
  amount: string; // base units (6 decimals)
  amountHuman: string; // e.g. "100.00"
  chainId: number;
  description: string;
  status: "pending" | "paid";
  txHash?: string;
  createdAt: string;
  paidAt?: string;
}

export interface CreateInvoiceRequest {
  merchantAddress: string;
  amount: string; // human-readable USDC
  chainId: number;
  description: string;
}
