// src/components/layout/ManagerNav.tsx
"use client"

import { Link } from "react-router-dom"
import { Moon, Sun, Package, Settings, BarChart3 } from "lucide-react"
import { Button } from "@/components/ui/Button"
import { useThemeStore } from "@/stores/themeStore"
import { useLanguageStore } from "@/stores/languageStore"

export function ManagerNav() {
    const { theme, toggleTheme } = useThemeStore()
    const { language, setLanguage, t } = useLanguageStore()

    return (
        <aside className="w-64 border-r border-border bg-card">
            <div className="flex h-16 items-center gap-2 border-b border-border px-6">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-accent">
                    <Package className="h-5 w-5 text-accent-foreground" />
                </div>
                <div>
                    <div className="text-sm font-bold text-card-foreground">ZooPet Manager</div>
                    <div className="text-xs text-muted-foreground">Inventory System</div>
                </div>
            </div>
            <nav className="flex flex-col gap-1 p-4">
                <Link
                    to="/manager"
                    className="flex items-center gap-3 rounded-lg bg-accent px-4 py-3 text-sm font-medium text-accent-foreground"
                >
                    <BarChart3 className="h-5 w-5" />
                    {t("manager.dashboard")}
                </Link>
                <Link
                    to="/manager/inventory"
                    className="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                    <Package className="h-5 w-5" />
                    {t("manager.inventory")} {/* FIXED: "Pets & Products" */}
                </Link>
                <Link
                    to="/"
                    className="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                    <Settings className="h-5 w-5" />
                    {t("common.close")}
                </Link>
            </nav>
            <div className="absolute bottom-4 left-4 right-4 flex items-center justify-between">
                <select
                    value={language}
                    onChange={(e) => setLanguage(e.target.value as "en" | "ru" | "kz")}
                    className="rounded-lg border border-border bg-background px-3 py-1.5 text-sm text-foreground"
                >
                    <option value="en">EN</option>
                    <option value="ru">RU</option>
                    <option value="kz">KZ</option>
                </select>
                <Button variant="ghost" size="sm" onClick={toggleTheme}>
                    {theme === "light" ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5" />}
                </Button>
            </div>
        </aside>
    )
}