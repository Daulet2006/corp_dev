"use client"

import { useState, useEffect } from "react"
import { CustomerNav } from "@/components/layout/CustomerNav"
import { ItemCard } from "@/components/ItemCard"
import { useLanguageStore } from "@/stores/languageStore"
import { useAuthStore } from "@/stores/authStore"  // Added missing import
import { useNavigate } from "react-router-dom"  // Для redirect если !auth
import type { Product } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"

export function MyProducts() {
    const navigate = useNavigate()
    const { t } = useLanguageStore()
    const [products, setProducts] = useState<Product[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        if (!useAuthStore.getState().isAuth) {
            navigate("/login")
            return
        }
        fetchMyProducts()
    }, [])

    const fetchMyProducts = async () => {
        try {
            const response = await api.get("/my/products")
            // Fixed: Extract from wrapped {success, data} response
            const apiData = response.data as { success: boolean; data: Product[]; message?: string }
            if (apiData.success && Array.isArray(apiData.data)) {
                setProducts(apiData.data)
            } else {
                throw new Error(apiData.message || "Invalid response format")
            }
        } catch (error: any) {
            console.error("Fetch my products error:", error)
            toast.error(error.response?.data?.message || "Failed to load your products")
            setProducts([])
        } finally {
            setLoading(false)
        }
    }

    const handleDelete = async (productId: number) => {
        if (!confirm("Are you sure you want to delete this product?")) return

        try {
            const response = await api.delete(`/products/${productId}`)
            if (response.data.success) {
                toast.success("Product deleted successfully")
                fetchMyProducts()
            } else {
                throw new Error(response.data.message)
            }
        } catch (error: any) {
            toast.error(error.response?.data?.message || "Delete failed")
        }
    }

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
                    <h1 className="text-3xl font-bold text-foreground">{t("my.products")}</h1>
                    <p className="mt-2 text-muted-foreground">Manage your products</p>
                </div>
                <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {products.map((product) => (
                        <ItemCard
                            key={product.id}
                            item={product}
                            type="product"
                            showActions
                            onDelete={() => handleDelete(product.id)}
                        />
                    ))}
                </div>
                {products.length === 0 && (
                    <div className="py-12 text-center">
                        <p className="text-muted-foreground">You don't have any products yet</p>
                    </div>
                )}
            </div>
        </div>
    )
}