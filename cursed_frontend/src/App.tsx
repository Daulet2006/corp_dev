// src/App.tsx (modified: unified nav import paths assumed as AppNav, added HomeWrapper for redirects, restricted store/my to user role)
"use client"

import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import { Toaster } from "react-hot-toast"
import React, { useEffect } from "react"
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
import { ManagerInventory } from "./pages/ManagerInventory.tsx"
import { AdminUsers } from "./pages/AdminUsers.tsx"

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

// Wrapper for Home route with role-based redirects
function HomeWrapper() {
    const { isLoading, isAuth, user } = useAuthStore()

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-screen">
                <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
            </div>
        )
    }

    if (isAuth) {
        if (user?.role === "admin") {
            return <Navigate to="/admin" replace />
        }
        if (user?.role === "manager") {
            return <Navigate to="/manager" replace />
        }
        // For user, show Home
        return <Home />
    }

    // Not auth, show Home
    return <Home />
}

function App() {
    const { setTheme } = useThemeStore()
    const { setLanguage } = useLanguageStore()

    useEffect(() => {
        const initApp = async () => {
            // Theme/lang first
            const savedTheme = localStorage.getItem("theme") as "light" | "dark" | null;
            if (savedTheme) setTheme(savedTheme);

            const savedLanguage = localStorage.getItem("language") as "en" | "ru" | "kz" | null;
            if (savedLanguage && (savedLanguage === "en" || savedLanguage === "ru" || savedLanguage === "kz")) {
                setLanguage(savedLanguage);
            } else {
                setLanguage("en");
            }

            // Auth rehydrate handled in store hydration — no manual call here
            // CSRF lazy-fetched in interceptor for protected routes
        };
        initApp();
    }, [setTheme, setLanguage]);



    return (
        <BrowserRouter>
            <Toaster position="top-right" />
            <Routes>
                <Route path="/" element={<HomeWrapper />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route
                    path="/store/pets"
                    element={
                        <ProtectedRoute requiredRole="user">
                            <StorePets />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/store/products"
                    element={
                        <ProtectedRoute requiredRole="user">
                            <StoreProducts />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/my/pets"
                    element={
                        <ProtectedRoute requiredRole="user">
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
                        <ProtectedRoute requiredRole="user">
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
                            <AdminUsers />
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