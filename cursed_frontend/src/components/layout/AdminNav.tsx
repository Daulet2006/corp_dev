// src/components/layout/AdminNav.tsx
"use client"

import { Link } from "react-router-dom"
import { Moon, Sun, Shield, Users, BarChart3, Settings } from "lucide-react"
import { Button } from "@/components/ui/Button"
import { useThemeStore } from "@/stores/themeStore"
import { useLanguageStore } from "@/stores/languageStore"

export function AdminNav() {
    const { theme, toggleTheme } = useThemeStore()
    const { language, setLanguage, t } = useLanguageStore()

    return (
        <aside className="w-64 border-r border-border bg-card">
            <div className="flex h-16 items-center gap-2 border-b border-border px-6">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
                    <Shield className="h-5 w-5 text-primary-foreground" />
                </div>
                <div>
                    <div className="text-sm font-bold text-card-foreground">ZooPet Admin</div>
                    <div className="text-xs text-muted-foreground">Control Panel</div>
                </div>
            </div>
            <nav className="flex flex-col gap-1 p-4">
                <Link
                    to="/admin"
                    className="flex items-center gap-3 rounded-lg bg-primary px-4 py-3 text-sm font-medium text-primary-foreground"
                >
                    <BarChart3 className="h-5 w-5" />
                    {t("admin.dashboard")}
                </Link>
                <Link
                    to="/admin/users"
                    className="flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                    <Users className="h-5 w-5" />
                    {t("admin.users")}
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