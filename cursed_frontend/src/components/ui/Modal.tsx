"use client"

import type { ReactNode } from "react"
import { X } from "lucide-react"
import { Button } from "./Button"
import { cn } from "@/lib/utils"

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title?: string
  children: ReactNode
  className?: string
}

export function Modal({ isOpen, onClose, title, children, className }: ModalProps) {
  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={onClose} />
      <div className={cn("relative z-50 w-full max-w-lg rounded-lg bg-card p-6 shadow-lg", className)}>
        <div className="mb-4 flex items-center justify-between">
          {title && <h2 className="text-xl font-semibold text-card-foreground">{title}</h2>}
          <Button variant="ghost" size="sm" onClick={onClose}>
            <X className="h-5 w-5" />
          </Button>
        </div>
        {children}
      </div>
    </div>
  )
}
