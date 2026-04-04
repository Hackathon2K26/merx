import { products } from "@/lib/products";
import { ProductCard } from "@/components/ProductCard";

export function ShopPage() {
  return (
    <div className="max-w-5xl mx-auto space-y-6">
      <div className="text-center space-y-2">
        <h1 className="text-3xl font-bold">Merx Bookshop</h1>
        <p className="text-muted-foreground">
          Premium crypto & Web3 ebooks — pay with USDC or swap from any asset
        </p>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {products.map((p) => (
          <ProductCard key={p.id} product={p} />
        ))}
      </div>
    </div>
  );
}
