// Fixed MyPets.tsx — extract data from wrapped APIResponse
"use client"

import { useState, useEffect } from "react"
import { AppNav } from "@/components/AppNav.tsx"
import { ItemCard } from "@/components/ItemCard"
import { useLanguageStore } from "@/stores/languageStore"
import type { Pet } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"

export function MyPets() {
    const { t } = useLanguageStore()
    const [pets, setPets] = useState<Pet[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        fetchMyPets()
    }, [])

    const fetchMyPets = async () => {
        try {
            const response = await api.get("/my/pets")
            // Fixed: Extract from wrapped {success, data} response
            const apiData = response.data as { success: boolean; data: Pet[] }
            if (apiData.success && Array.isArray(apiData.data)) {
                setPets(apiData.data)
            } else {
                throw new Error("Invalid response format")
            }
        } catch (error: any) {
            console.error("Fetch my pets error:", error)
            toast.error(error.response?.data?.message || error.response?.data?.error || "Failed to load your pets")
            setPets([])
        } finally {
            setLoading(false)
        }
    }

    const handleDelete = async (petId: number) => {
        if (!confirm("Are you sure you want to delete this pet?")) return

        try {
            await api.delete(`/pets/${petId}`)
            toast.success("Pet deleted successfully")
            fetchMyPets()
        } catch (error: any) {
            toast.error(error.response?.data?.error || "Delete failed")
        }
    }

    if (loading) {
        return (
            <div className="min-h-screen bg-background">
                <AppNav />
                <div className="flex h-96 items-center justify-center">
                    <p className="text-muted-foreground">{t("common.loading")}</p>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-background">
            <AppNav />
            <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold text-foreground">{t("my.pets")}</h1>
                    <p className="mt-2 text-muted-foreground">Manage your pets</p>
                </div>
                <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {pets.map((pet) => (
                        <ItemCard key={pet.id} item={pet} type="pet" showActions onDelete={() => handleDelete(pet.id)} />
                    ))}
                </div>
                {pets.length === 0 && (
                    <div className="py-12 text-center">
                        <p className="text-muted-foreground">You don't have any pets yet</p>
                    </div>
                )}
            </div>
        </div>
    )
}