// Updated App.tsx — ProtectedRoute now handles isLoading
"use client"

import type React from "react"

import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { Toaster } from "react-hot-toast"
import { useEffect } from "react"
import { useAuthStore } from "./stores/authStore"
import { useThemeStore } from "./stores/themeStore"
import { useLanguageStore } from "./stores/languageStore"
import { Home } from "./pages/Home"
import { Login } from "./pages/Login"
import { Register } from "./pages/Register"
import { StorePets } from "./pages/StorePets"
import { StoreProducts } from "./pages/StoreProducts"
import { MyPets } from "./pages/MyPets"
import { MyProducts } from "./pages/MyProducts"
import { Profile } from "./pages/Profile"
import { AdminDashboard } from "./pages/AdminDashboard"
import { ManagerDashboard } from "./pages/ManagerDashboard"
import {ManagerInventory} from "@/pages/ManagerInvertory.tsx";
import {AdminUsers} from "@/pages/AdminUsers.tsx";

function ProtectedRoute({ children, requiredRole }: { children: React.ReactNode; requiredRole?: string }) {
    const { isLoading, isAuth, user } = useAuthStore()

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
            </div>
        )
    }

    if (!isAuth) {
        return <Navigate to="/login" replace />
    }

    if (requiredRole && user?.role !== requiredRole && !(requiredRole === "manager" && user?.role === "admin")) {
        return <Navigate to="/" replace />
    }

    return <>{children}</>
}

function App() {
    const { fetchUser } = useAuthStore()
    const { setTheme } = useThemeStore()
    const { setLanguage } = useLanguageStore()

    useEffect(() => {
        fetchUser()  // Sync, sets isLoading false immediately
        const savedTheme = localStorage.getItem("theme") as "light" | "dark" | null
        if (savedTheme) setTheme(savedTheme)

        const savedLanguage = localStorage.getItem("language") as "en" | "ru" | "kz" | null
        if (savedLanguage) setLanguage(savedLanguage)
    }, [fetchUser, setTheme, setLanguage])

    return (
        <BrowserRouter>
            <Toaster position="top-right" />
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/store/pets" element={<StorePets />} />
                <Route path="/store/products" element={<StoreProducts />} />
                <Route
                    path="/my/pets"
                    element={
                        <ProtectedRoute>
                            <MyPets />
                        </ProtectedRoute>
                    }
                />
                <Route
                path="/manager/inventory"
                element={
                    <ProtectedRoute requiredRole="manager">
                        <ManagerInventory />
                    </ProtectedRoute>
                }
            />
                <Route
                    path="/my/products"
                    element={
                        <ProtectedRoute>
                            <MyProducts />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/profile"
                    element={
                        <ProtectedRoute>
                            <Profile />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/admin"
                    element={
                        <ProtectedRoute requiredRole="admin">
                            <AdminDashboard />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/admin/users"
                    element={
                        <ProtectedRoute requiredRole="admin">
                            <AdminUsers/>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/manager"
                    element={
                        <ProtectedRoute requiredRole="manager">
                            <ManagerDashboard />
                        </ProtectedRoute>
                    }
                />
            </Routes>
        </BrowserRouter>
    )
}

export default App