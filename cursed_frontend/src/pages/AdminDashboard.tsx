"use client"

import { useState, useEffect } from "react"
import { AdminNav } from "@/components/layout/AdminNav"
import { Card, CardContent } from "@/components/ui/Card"
import { useLanguageStore } from "@/stores/languageStore"
import type { Stats } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"
import { Users, Package, ShoppingBag, TrendingUp } from "lucide-react"

export function AdminDashboard() {
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
        <AdminNav />
        <div className="flex flex-1 items-center justify-center">
          <p className="text-muted-foreground">{t("common.loading")}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex h-screen">
      <AdminNav />
      <main className="flex-1 overflow-y-auto bg-background">
        <div className="p-8">
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-foreground">{t("admin.dashboard")}</h1>
            <p className="mt-2 text-muted-foreground">Monitor system performance and user activity</p>
          </div>

          {stats && (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
              <Card>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-muted-foreground">{t("admin.users")}</p>
                      <p className="mt-2 text-3xl font-bold text-foreground">{stats.users}</p>
                    </div>
                    <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary">
                      <Users className="h-6 w-6 text-white" />
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium text-muted-foreground">{t("admin.totalPets")}</p>
                      <p className="mt-2 text-3xl font-bold text-foreground">{stats.totalPets}</p>
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
                      <p className="text-sm font-medium text-muted-foreground">{t("admin.totalProducts")}</p>
                      <p className="mt-2 text-3xl font-bold text-foreground">{stats.totalProducts}</p>
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
                      <p className="text-sm font-medium text-muted-foreground">Store Items</p>
                      <p className="mt-2 text-3xl font-bold text-foreground">{stats.storePets + stats.storeProducts}</p>
                    </div>
                    <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-accent">
                      <TrendingUp className="h-6 w-6 text-white" />
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
