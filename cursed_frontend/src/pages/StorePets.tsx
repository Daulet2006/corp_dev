"use client"

import { useState, useEffect } from "react"
import { CustomerNav } from "@/components/layout/CustomerNav"
import { ItemCard } from "@/components/ItemCard"
import { useLanguageStore } from "@/stores/languageStore"
import { useAuthStore } from "@/stores/authStore"
import type { Pet } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"
import { Modal } from "@/components/ui/Modal"
import { Button } from "@/components/ui/Button"

export function StorePets() {
  const { t } = useLanguageStore()
  const { isAuth } = useAuthStore()
  const [pets, setPets] = useState<Pet[]>([])
  const [loading, setLoading] = useState(true)
  const [showLoginModal, setShowLoginModal] = useState(false)

  useEffect(() => {
    fetchPets()
  }, [])

  const fetchPets = async () => {
    try {
      const response = await api.get("/pets?owner_id=0")
      setPets(response.data)
    } catch (error) {
      toast.error("Failed to load pets")
    } finally {
      setLoading(false)
    }
  }

  const handleBuy = async (petId: number) => {
    if (!isAuth) {
      setShowLoginModal(true)
      return
    }

    try {
      await api.post(`/pets/${petId}/buy`)
      toast.success("Pet purchased successfully!")
      fetchPets()
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Purchase failed")
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
          <h1 className="text-3xl font-bold text-foreground">{t("store.pets")}</h1>
          <p className="mt-2 text-muted-foreground">Find your perfect companion</p>
        </div>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {pets.map((pet) => (
            <ItemCard key={pet.id} item={pet} type="pet" onBuy={() => handleBuy(pet.id)} />
          ))}
        </div>
        {pets.length === 0 && (
          <div className="py-12 text-center">
            <p className="text-muted-foreground">No pets available at the moment</p>
          </div>
        )}
      </div>

      <Modal isOpen={showLoginModal} onClose={() => setShowLoginModal(false)} title={t("auth.loginPrompt")}>
        <p className="mb-4 text-muted-foreground">Please login to purchase pets</p>
        <Button variant="primary" className="w-full" onClick={() => (window.location.href = "/login")}>
          {t("auth.login")}
        </Button>
      </Modal>
    </div>
  )
}
