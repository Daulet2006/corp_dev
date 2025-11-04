// src/pages/StoreProducts.tsx
"use client"

import { useState, useEffect } from "react"
import { CustomerNav } from "@/components/layout/CustomerNav"
import { ItemCard } from "@/components/ItemCard"
import { useLanguageStore } from "@/stores/languageStore"
import { useAuthStore } from "@/stores/authStore"
import type { Product } from "@/types/types.ts"
import type { ApiResponse } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"
import { Modal } from "@/components/ui/Modal"
import { Button } from "@/components/ui/Button"

export function StoreProducts() {
    const { t } = useLanguageStore()
    const [products, setProducts] = useState<Product[]>([])
    const [loading, setLoading] = useState(true)
    const [showLoginModal, setShowLoginModal] = useState(false)

    useEffect(() => {
        fetchProducts()
    }, [])

    const fetchProducts = async () => {
        try {
            const response = await api.get<ApiResponse<Product[]>>("/products?owner_id=0")
            if (response.data.success) {
                setProducts(response.data.data || [])
            } else {
                toast.error(response.data.message || "Failed to load products")
                setProducts([])
            }
        } catch (error: any) {
            toast.error(error.response?.data?.message || "Failed to load products")
            setProducts([])
        } finally {
            setLoading(false)
        }
    }

    const handleBuy = async (id: number) => {
        const { isAuth, token } = useAuthStore.getState();
        if (!isAuth || !token) {
            setShowLoginModal(true);
            return;
        }

        try {
            // Interceptor сам retry'ит CSRF — не нужно manual loop
            const response = await api.post<ApiResponse<Product>>(`/products/${id}/buy`);
            if (response.data.success) {
                toast.success("Purchased!");
                fetchProducts();
            } else {
                toast.error(response.data.message || "Failed");
            }
        } catch (error: any) {
            // Только финальная ошибка — без CSRF spam
            toast.error(error.response?.data?.message || "Purchase failed");
        }
    };

    if (loading) {
        return (
            <div className="min-h-screen bg-background">
                <CustomerNav />
                <div className="flex h-96 items-center justify-center">
                    <p className="text-muted-foreground">{t("common.loading")}</p>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-background">
            <CustomerNav />
            <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-foreground">{t("store.products")}</h1>
                    <p className="mt-2 text-muted-foreground">Everything your pet needs</p>
                </div>
                <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {products.map((product) => (
                        <ItemCard key={product.id} item={product} type="product" onBuy={() => handleBuy(product.id)} />
                    ))}
                </div>
                {products.length === 0 && (
                    <div className="py-12 text-center">
                        <p className="text-muted-foreground">No products available at the moment</p>
                    </div>
                )}
            </div>

            <Modal isOpen={showLoginModal} onClose={() => setShowLoginModal(false)} title={t("auth.loginPrompt")}>
                <p className="mb-4 text-muted-foreground">Please login to purchase products</p>
                <Button variant="default" className="w-full" onClick={() => (window.location.href = "/login")}>
                    {t("auth.login")}
                </Button>
            </Modal>
        </div>
    )
}