import { Link } from "react-router-dom";
import type { Product } from "@/lib/products";

export function ProductCard({ product }: { product: Product }) {
  return (
    <div className="rounded-lg border border-border bg-card p-5 flex flex-col gap-3 hover:border-primary/40 transition-colors">
      <div className="text-6xl text-center py-4">{product.cover}</div>
      <div className="flex-1 space-y-1">
        <h3 className="font-semibold text-foreground leading-tight">
          {product.title}
        </h3>
        <p className="text-sm text-muted-foreground">{product.author}</p>
        <p className="text-xs text-muted-foreground leading-relaxed mt-2">
          {product.description}
        </p>
      </div>
      <div className="flex items-center justify-between pt-2">
        <span className="text-lg font-bold text-primary">
          {product.price} USDC
        </span>
        <Link
          to={`/checkout/${product.id}`}
          className="rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          Buy Now
        </Link>
      </div>
    </div>
  );
}
