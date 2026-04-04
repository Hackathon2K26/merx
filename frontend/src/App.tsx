import { Routes, Route, Navigate } from "react-router-dom";
import { Layout } from "./components/Layout";
import { ShopPage } from "./pages/ShopPage";
import { CheckoutPage } from "./pages/CheckoutPage";

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<ShopPage />} />
        <Route path="/checkout/:productId" element={<CheckoutPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Layout>
  );
}
