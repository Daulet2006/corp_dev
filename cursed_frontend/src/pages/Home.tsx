import { CustomerNav } from "@/components/layout/CustomerNav"
import { useLanguageStore } from "@/stores/languageStore"
import { Button } from "@/components/ui/Button"
import { Card, CardContent } from "@/components/ui/Card"
import { Heart, ShoppingBag, Shield, Truck } from "lucide-react"
import { Link } from "react-router-dom"

export function Home() {
  const { t } = useLanguageStore()

  const mockPets = [
    { id: 1, name: "Golden Retriever", breed: "Dog", age: 2, price: 500, image: null },
    { id: 2, name: "Persian Cat", breed: "Cat", age: 1, price: 300, image: null },
    { id: 3, name: "Parrot", breed: "Bird", age: 1, price: 200, image: null },
    { id: 4, name: "Rabbit", breed: "Rabbit", age: 1, price: 100, image: null },
  ]

  const mockProducts = [
    { id: 1, name: "Premium Dog Food", category: "Food", price: 45, image: null },
    { id: 2, name: "Cat Scratching Post", category: "Toys", price: 35, image: null },
    { id: 3, name: "Bird Cage", category: "Housing", price: 120, image: null },
    { id: 4, name: "Pet Carrier", category: "Accessories", price: 55, image: null },
  ]

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
              <p className="mt-6 text-pretty text-lg leading-relaxed text-muted-foreground">{t("hero.subtitle")}</p>
              <div className="mt-8 flex flex-wrap gap-4">
                <Link to="/store/pets">
                  <Button variant="primary" size="lg">
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
            <div className="flex flex-col items-center text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <Shield className="h-6 w-6" />
              </div>
              <h3 className="mt-4 font-semibold text-foreground">Verified Pets</h3>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                All pets are health-checked and certified
              </p>
            </div>
            <div className="flex flex-col items-center text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-accent text-accent-foreground">
                <Heart className="h-6 w-6" />
              </div>
              <h3 className="mt-4 font-semibold text-foreground">Lifetime Support</h3>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                Expert guidance for your pet's entire life
              </p>
            </div>
            <div className="flex flex-col items-center text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <Truck className="h-6 w-6" />
              </div>
              <h3 className="mt-4 font-semibold text-foreground">Safe Delivery</h3>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                Secure transport for pets and products
              </p>
            </div>
            <div className="flex flex-col items-center text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-accent text-accent-foreground">
                <ShoppingBag className="h-6 w-6" />
              </div>
              <h3 className="mt-4 font-semibold text-foreground">Premium Products</h3>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                Curated selection of quality supplies
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Available Pets */}
      <section className="py-16">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="mb-12 text-center">
            <h2 className="text-balance text-3xl font-bold text-foreground sm:text-4xl">{t("store.pets")}</h2>
            <p className="mt-4 text-pretty text-lg text-muted-foreground">Meet your new best friend</p>
          </div>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {mockPets.map((pet) => (
              <Card key={pet.id} className="overflow-hidden transition-shadow hover:shadow-lg">
                <div className="relative aspect-[4/3] overflow-hidden">
                  <img
                    src={`/.jpg?height=300&width=400&query=${encodeURIComponent(pet.name + " " + pet.breed)}`}
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
                    <Button variant="primary" size="sm">
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

      {/* Featured Products */}
      <section className="bg-muted/30 py-16">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="mb-12 text-center">
            <h2 className="text-balance text-3xl font-bold text-foreground sm:text-4xl">{t("store.products")}</h2>
            <p className="mt-4 text-pretty text-lg text-muted-foreground">Everything your pet needs</p>
          </div>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {mockProducts.map((product) => (
              <Card key={product.id} className="overflow-hidden transition-shadow hover:shadow-lg">
                <div className="relative aspect-square overflow-hidden">
                  <img
                    src={`/.jpg?height=400&width=400&query=${encodeURIComponent(product.name)}`}
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
          </div>
          <div className="mt-12 border-t border-border pt-8 text-center text-sm text-muted-foreground">
            © 2025 ZooPet. All rights reserved.
          </div>
        </div>
      </footer>
    </div>
  )
}
