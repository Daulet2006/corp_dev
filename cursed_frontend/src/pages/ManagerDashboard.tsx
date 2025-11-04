// src/pages/ManagerDashboard.tsx (modified: unified nav, no sidebar, full-width content, adjusted loading)
"use client"

import { useState, useEffect } from "react"
import { AppNav } from "@/components/AppNav.tsx"  // Unified nav
import { Card, CardContent } from "@/components/ui/Card"
import { useLanguageStore } from "@/stores/languageStore"
import type { Stats } from "@/types/types.ts"
import type { ApiResponse } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"
import { Package, ShoppingBag, TrendingUp, AlertCircle } from "lucide-react"

export function ManagerDashboard() {
    const { t } = useLanguageStore()
    const [stats, setStats] = useState<Stats | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        fetchStats()
    }, [])

    const fetchStats = async () => {
        try {
            const response = await api.get<ApiResponse<Stats>>("/stats")
            setStats(response.data.data || null)
        } catch (error) {
            toast.error("Failed to load stats")
        } finally {
            setLoading(false)
        }
    }

    if (loading) {
        return (
            <div className="min-h-screen bg-background">
                <AppNav />
                <div className="flex h-[70vh] items-center justify-center">
                    <p className="text-muted-foreground">{t("common.loading")}</p>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-background">
            <AppNav />
            <main className="overflow-y-auto">
                <div className="p-8">
                    <div className="mb-8">
                        <h1 className="text-3xl font-bold text-foreground">{t("manager.dashboard")}</h1>
                        <p className="mt-2 text-muted-foreground">Manage your store inventory</p>
                    </div>

                    {stats && (
                        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
                            <Card>
                                <CardContent className="p-6">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="text-sm font-medium text-muted-foreground">{t("manager.storePets")}</p>
                                            <p className="mt-2 text-3xl font-bold text-foreground">{stats.storePets}</p>
                                        </div>
                                        <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-accent">
                                            <Package className="h-6 w-6 text-white" />
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardContent className="p-6">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="text-sm font-medium text-muted-foreground">{t("manager.storeProducts")}</p>
                                            <p className="mt-2 text-3xl font-bold text-foreground">{stats.storeProducts}</p>
                                        </div>
                                        <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary">
                                            <ShoppingBag className="h-6 w-6 text-white" />
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardContent className="p-6">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="text-sm font-medium text-muted-foreground">Total Inventory</p>
                                            <p className="mt-2 text-3xl font-bold text-foreground">{stats.storePets + stats.storeProducts}</p>
                                        </div>
                                        <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-accent">
                                            <TrendingUp className="h-6 w-6 text-white" />
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardContent className="p-6">
                                    <div className="flex items-center justify-between">
                                        <div>
                                            <p className="text-sm font-medium text-muted-foreground">Owned Items</p>
                                            <p className="mt-2 text-3xl font-bold text-foreground">{stats.ownedPets + stats.ownedProducts}</p>
                                        </div>
                                        <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary">
                                            <AlertCircle className="h-6 w-6 text-white" />
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        </div>
                    )}
                </div>
            </main>
        </div>
    )
}