"use client";

import { Link } from "react-router-dom";
import { Moon, Sun, ShoppingCart, User, Menu, LogOut } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/Button";
import { useThemeStore } from "@/stores/themeStore";
import { useLanguageStore } from "@/stores/languageStore";
import { useAuthStore } from "@/stores/authStore";

export function CustomerNav() {
  const { theme, toggleTheme } = useThemeStore();
  const { language, setLanguage, t } = useLanguageStore();
  const { isAuth, user, logout } = useAuthStore();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <nav className="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <Link to="/" className="flex items-center gap-2">
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary">
              <span className="text-xl font-bold text-primary-foreground">
                🐾
              </span>
            </div>
            <span className="text-xl font-bold text-foreground">ZooPet</span>
          </Link>

          <div className="hidden items-center gap-6 md:flex">
            <Link
              to="/"
              className="text-sm font-medium text-foreground transition-colors hover:text-primary"
            >
              {t("nav.home")}
            </Link>
            <Link
              to="/store/pets"
              className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
            >
              {t("nav.pets")}
            </Link>
            <Link
              to="/store/products"
              className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
            >
              {t("nav.products")}
            </Link>
            {isAuth && (
              <>
                <Link
                  to="/profile"
                  className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
                >
                  {t("nav.profile")}
                </Link>
                <Link
                  to="/my/pets"
                  className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
                >
                    {t("nav.myPets")}
                </Link>
                 <Link
                  to="/my/products"
                  className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
                >
                    {t("nav.myPets")}
                </Link>
                {user?.role === "admin" && (
                  <Link
                    to="/admin"
                    className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
                  >
                    {t("nav.admin")}
                  </Link>
                )}
                {(user?.role === "manager" || user?.role === "admin") && (
                  <Link
                    to="/manager"
                    className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
                  >
                    {t("nav.manager")}
                  </Link>
                )}
              </>
            )}
          </div>

          <div className="flex items-center gap-2">
            <select
              value={language}
              onChange={(e) =>
                setLanguage(e.target.value as "en" | "ru" | "kz")
              }
              className="rounded-lg border border-border bg-background px-3 py-1.5 text-sm text-foreground"
            >
              <option value="en">EN</option>
              <option value="ru">RU</option>
              <option value="kz">KZ</option>
            </select>

            <Button variant="ghost" size="sm" onClick={toggleTheme}>
              {theme === "light" ? (
                <Moon className="h-5 w-5" />
              ) : (
                <Sun className="h-5 w-5" />
              )}
            </Button>

            {isAuth ? (
              <>
                <Button variant="ghost" size="sm">
                  <ShoppingCart className="h-5 w-5" />
                </Button>
                <Link to="/profile">
                  <Button variant="ghost" size="sm">
                    <User className="h-5 w-5" />
                  </Button>
                </Link>
                <Button variant="ghost" size="sm" onClick={logout}>
                  <LogOut className="h-5 w-5" />
                </Button>
              </>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="outline" size="sm">
                    {t("nav.login")}
                  </Button>
                </Link>
                <Link to="/register">
                  <Button variant="primary" size="sm">
                    {t("nav.register")}
                  </Button>
                </Link>
              </>
            )}

            <Button
              variant="ghost"
              size="sm"
              className="md:hidden"
              onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            >
              <Menu className="h-5 w-5" />
            </Button>
          </div>
        </div>

        {mobileMenuOpen && (
          <div className="border-t border-border py-4 md:hidden">
            <div className="flex flex-col gap-4">
              <Link to="/" className="text-sm font-medium text-foreground">
                {t("nav.home")}
              </Link>
              <Link
                to="/store/pets"
                className="text-sm font-medium text-muted-foreground"
              >
                {t("nav.pets")}
              </Link>
              <Link
                to="/store/products"
                className="text-sm font-medium text-muted-foreground"
              >
                {t("nav.products")}
              </Link>
              {isAuth && (
                <>
                  <Link
                    to="/my/pets"
                    className="text-sm font-medium text-muted-foreground"
                  >
                    {t("nav.myPets")}
                  </Link>
                  <Link
                    to="/my/products"
                    className="text-sm font-medium text-muted-foreground"
                  >
                    {t("nav.myProducts")}
                  </Link>
                </>
              )}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
}
