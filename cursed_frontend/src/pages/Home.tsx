// src/pages/Home.tsx
"use client"

import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import { CustomerNav } from "@/components/layout/CustomerNav"
import { Button } from "@/components/ui/Button"
import { Card, CardContent } from "@/components/ui/Card"
import { Heart, ShoppingBag, Shield, Truck } from "lucide-react"
import { useLanguageStore } from "@/stores/languageStore"
import type { Pet, Product } from "@/types/types.ts"
import api from "@/utils/api"
import toast from "react-hot-toast"

export function Home() {
    const { t } = useLanguageStore()
    const [pets, setPets] = useState<Pet[]>([])
    const [products, setProducts] = useState<Product[]>([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        const fetchData = async () => {
            try {
                console.log("Fetching pets/products..."); // Debug
                const [petsRes, productsRes] = await Promise.all([
                    api.get("/pets").catch(e => { console.error("Pets fetch error:", e); throw e; }),
                    api.get("/products").catch(e => { console.error("Products fetch error:", e); throw e; })
                ]);
                console.log("Fetched:", petsRes.data, productsRes.data); // Debug
                setPets(petsRes.data.data || []);
                setProducts(productsRes.data.data || []);
            } catch (error: any) {
                console.error("Failed to load data:", error);
                toast.error("Failed to load pets or products");
            } finally {
                setLoading(false); // Always!
            }
        };

        fetchData()
    }, [])

    const resolveImage = (url?: string | null, fallback: string = "/placeholder-pet.jpg") => {
        if (!url) return fallback
        if (url.startsWith("http")) return url
        return `${import.meta.env.VITE_API_URL}${url}`
    }

    if (loading) {
        return (
            <div className="min-h-screen bg-background">
                <CustomerNav />
                <div className="flex h-[70vh] items-center justify-center">
                    <p className="text-lg text-muted-foreground">{t("common.loading")}</p>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen bg-background">
            <CustomerNav />

            {/* Hero Section */}
            <section className="relative overflow-hidden bg-gradient-to-br from-primary/10 via-background to-accent/10">
                <div className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8 lg:py-32">
                    <div className="grid gap-12 lg:grid-cols-2 lg:gap-16">
                        <div className="flex flex-col justify-center">
                            <h1 className="text-balance text-4xl font-bold tracking-tight text-foreground sm:text-5xl lg:text-6xl">
                                {t("hero.title")}
                            </h1>
                            <p className="mt-6 text-pretty text-lg leading-relaxed text-muted-foreground">
                                {t("hero.subtitle")}
                            </p>
                            <div className="mt-8 flex flex-wrap gap-4">
                                <Link to="/store/pets">
                                    <Button variant="default" size="lg">
                                        {t("hero.shopPets")}
                                    </Button>
                                </Link>
                                <Link to="/store/products">
                                    <Button variant="outline" size="lg">
                                        {t("hero.shopProducts")}
                                    </Button>
                                </Link>
                            </div>
                        </div>
                        <div className="relative">
                            <img
                                src="/happy-pets-together.jpg"
                                alt="Happy pets"
                                className="h-full w-full rounded-2xl object-cover shadow-2xl"
                            />
                        </div>
                    </div>
                </div>
            </section>

            {/* Features */}
            <section className="border-y border-border bg-muted/30 py-16">
                <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
                        <Feature icon={<Shield className="h-6 w-6" />} title="Verified Pets" text="All pets are health-checked and certified" />
                        <Feature icon={<Heart className="h-6 w-6" />} title="Lifetime Support" text="Expert guidance for your pet's entire life" />
                        <Feature icon={<Truck className="h-6 w-6" />} title="Safe Delivery" text="Secure transport for pets and products" />
                        <Feature icon={<ShoppingBag className="h-6 w-6" />} title="Premium Products" text="Curated selection of quality supplies" />
                    </div>
                </div>
            </section>

            {/* Pets Section */}
            <section className="py-16">
                <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <SectionHeader title={t("store.pets")} subtitle="Meet your new best friend" />
                    <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
                        {pets.slice(0, 4).map((pet) => (
                            <Card key={pet.id} className="overflow-hidden transition-shadow hover:shadow-lg">
                                <div className="relative aspect-[4/3] overflow-hidden">
                                    <img
                                        src={resolveImage(pet.image)}
                                        alt={pet.name}
                                        className="h-full w-full object-cover"
                                    />
                                </div>
                                <CardContent className="p-4">
                                    <h3 className="font-semibold text-card-foreground">{pet.name}</h3>
                                    <p className="mt-1 text-sm text-muted-foreground">
                                        {pet.breed} • {pet.age} years
                                    </p>
                                    <div className="mt-4 flex items-center justify-between">
                                        <span className="text-lg font-bold text-primary">${pet.price}</span>
                                        <Button variant="default" size="sm">
                                            {t("store.buy")}
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                    <div className="mt-8 text-center">
                        <Link to="/store/pets">
                            <Button variant="outline" size="lg">
                                View All Pets
                            </Button>
                        </Link>
                    </div>
                </div>
            </section>

            {/* Products Section */}
            <section className="bg-muted/30 py-16">
                <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <SectionHeader title={t("store.products")} subtitle="Everything your pet needs" />
                    <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
                        {products.slice(0, 4).map((product) => (
                            <Card key={product.id} className="overflow-hidden transition-shadow hover:shadow-lg">
                                <div className="relative aspect-square overflow-hidden">
                                    <img
                                        src={resolveImage(product.image, "/placeholder-product.jpg")}
                                        alt={product.name}
                                        className="h-full w-full object-cover"
                                    />
                                </div>
                                <CardContent className="p-4">
                                    <h3 className="font-semibold text-card-foreground">{product.name}</h3>
                                    <p className="mt-1 text-sm text-muted-foreground">{product.category}</p>
                                    <div className="mt-4 flex items-center justify-between">
                                        <span className="text-lg font-bold text-accent">${product.price}</span>
                                        <Button variant="outline" size="sm">
                                            {t("store.addToCart")}
                                        </Button>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                    <div className="mt-8 text-center">
                        <Link to="/store/products">
                            <Button variant="outline" size="lg">
                                View All Products
                            </Button>
                        </Link>
                    </div>
                </div>
            </section>

            {/* Footer */}
            <footer className="border-t border-border bg-card py-12">
                <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
                    <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
                        <FooterSection />
                    </div>
                    <div className="mt-12 border-t border-border pt-8 text-center text-sm text-muted-foreground">
                        © 2025 ZooPet. All rights reserved.
                    </div>
                </div>
            </footer>
        </div>
    )
}

// === Subcomponents ===
function Feature({ icon, title, text }: { icon: React.ReactNode; title: string; text: string }) {
    return (
        <div className="flex flex-col items-center text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                {icon}
            </div>
            <h3 className="mt-4 font-semibold text-foreground">{title}</h3>
            <p className="mt-2 text-sm leading-relaxed text-muted-foreground">{text}</p>
        </div>
    )
}

function SectionHeader({ title, subtitle }: { title: string; subtitle: string }) {
    return (
        <div className="mb-12 text-center">
            <h2 className="text-balance text-3xl font-bold text-foreground sm:text-4xl">{title}</h2>
            <p className="mt-4 text-pretty text-lg text-muted-foreground">{subtitle}</p>
        </div>
    )
}

function FooterSection() {
    return (
        <>
            <div>
                <div className="flex items-center gap-2">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
                        <span className="text-xl font-bold text-primary-foreground">🐾</span>
                    </div>
                    <span className="text-xl font-bold text-card-foreground">ZooPet</span>
                </div>
                <p className="mt-4 text-sm leading-relaxed text-muted-foreground">
                    Your trusted partner in pet adoption and care.
                </p>
            </div>
            <div>
                <h4 className="font-semibold text-card-foreground">Shop</h4>
                <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
                    <li>
                        <Link to="/store/pets" className="hover:text-primary">
                            Pets
                        </Link>
                    </li>
                    <li>
                        <Link to="/store/products" className="hover:text-primary">
                            Products
                        </Link>
                    </li>
                </ul>
            </div>
            <div>
                <h4 className="font-semibold text-card-foreground">Company</h4>
                <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
                    <li>
                        <a href="#" className="hover:text-primary">
                            About Us
                        </a>
                    </li>
                    <li>
                        <a href="#" className="hover:text-primary">
                            Contact
                        </a>
                    </li>
                </ul>
            </div>
            <div>
                <h4 className="font-semibold text-card-foreground">Support</h4>
                <ul className="mt-4 space-y-2 text-sm text-muted-foreground">
                    <li>
                        <a href="#" className="hover:text-primary">
                            Help Center
                        </a>
                    </li>
                    <li>
                        <a href="#" className="hover:text-primary">
                            Shipping Info
                        </a>
                    </li>
                </ul>
            </div>
        </>
    )
}