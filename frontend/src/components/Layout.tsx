import type { ReactNode } from "react";
import { Link } from "react-router-dom";
import { ConnectWallet } from "./ConnectWallet";

export function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-border px-6 py-4 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 text-xl font-bold text-foreground">
          <img src="/logo.png" alt="Merx" className="h-8 w-8 rounded-full" />
          Merx
        </Link>
        <ConnectWallet />
      </header>
      <main className="flex-1 p-6">{children}</main>
    </div>
  );
}
