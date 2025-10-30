"use client"

import { useState, useEffect } from "react"
import { ManagerNav } from "@/components/layout/ManagerNav"
import { Card, CardContent } from "@/components/ui/Card"
import { useLanguageStore } from "@/stores/languageStore"
import type { Stats } from "@/types/types.ts"
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
      const response = await api.get("/stats")
      setStats(response.data)
    } catch (error) {
      toast.error("Failed to load stats")
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex h-screen">
        <ManagerNav />
        <div className="flex flex-1 items-center justify-center">
          <p className="text-muted-foreground">{t("common.loading")}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex h-screen">
      <ManagerNav />
      <main className="flex-1 overflow-y-auto bg-background">
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
