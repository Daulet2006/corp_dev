"use client"

import { Button } from "./ui/Button"
import { Card, CardContent } from "./ui/Card"
import { Edit, Trash2 } from "lucide-react"
import type { Pet, Product } from "@/types/types.ts"

interface ItemCardProps {
  item: Pet | Product
  type: "pet" | "product"
  onBuy?: () => void
  onEdit?: () => void
  onDelete?: () => void
  showActions?: boolean
}

export function ItemCard({ item, type, onBuy, onEdit, onDelete, showActions }: ItemCardProps) {
  const isPet = type === "pet"
  const pet = isPet ? (item as Pet) : null
  const product = !isPet ? (item as Product) : null

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-lg">
      <div className={`relative ${isPet ? "aspect-[4/3]" : "aspect-square"} overflow-hidden`}>
        <img
          src={item.image || `/placeholder.svg?height=300&width=400&query=${encodeURIComponent(item.name)}`}
          alt={item.name}
          className="h-full w-full object-cover"
        />
      </div>
      <CardContent className="p-4">
        <h3 className="font-semibold text-card-foreground">{item.name}</h3>
        {isPet && pet && (
          <p className="mt-1 text-sm text-muted-foreground">
            {pet.breed} • {pet.age} years • {pet.gender}
          </p>
        )}
        {!isPet && product && (
          <p className="mt-1 text-sm text-muted-foreground">
            {product.category} • Stock: {product.stock}
          </p>
        )}
        {item.description && (
          <p className="mt-2 line-clamp-2 text-sm leading-relaxed text-muted-foreground">{item.description}</p>
        )}
        <div className="mt-4 flex items-center justify-between">
          <span className="text-lg font-bold text-primary">${item.price.toFixed(2)}</span>
          {showActions ? (
            <div className="flex gap-2">
              {onEdit && (
                <Button variant="outline" size="sm" onClick={onEdit}>
                  <Edit className="h-4 w-4" />
                </Button>
              )}
              {onDelete && (
                <Button variant="outline" size="sm" onClick={onDelete}>
                  <Trash2 className="h-4 w-4" />
                </Button>
              )}
            </div>
          ) : (
            onBuy && (
              <Button variant="primary" size="sm" onClick={onBuy}>
                Buy Now
              </Button>
            )
          )}
        </div>
      </CardContent>
    </Card>
  )
}
