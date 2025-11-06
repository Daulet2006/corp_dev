"use client"

import { useState } from "react"
import { useNavigate, Link } from "react-router-dom"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import toast from "react-hot-toast"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { useAuthStore } from "@/stores/authStore"
import { useLanguageStore } from "@/stores/languageStore"
import api, { fetchCsrfToken } from "@/utils/api"  // Import fetchCsrfToken

const registerSchema = z
    .object({
        firstName: z
            .string()
            .min(2, "First name must be at least 2 characters")
            .max(50, "First name must be at most 50 characters")
            .regex(/^[A-Za-zА-Яа-яӘәӨөҰұҚқҒғІіЁёЪъЬь\s-]+$/, "Invalid characters in first name"),

        lastName: z
            .string()
            .min(2, "Last name must be at least 2 characters")
            .max(50, "Last name must be at most 50 characters")
            .regex(/^[A-Za-zА-Яа-яӘәӨөҰұҚқҒғІіЁёЪъЬь\s-]+$/, "Invalid characters in last name"),

        email: z
            .string()
            .email("Invalid email address")
            .max(100, "Email must be at most 100 characters"),

        password: z
            .string()
            .min(8, "Password must be at least 8 characters")
            .max(64, "Password must be at most 64 characters")
            .regex(/[A-Z]/, "Password must contain at least one uppercase letter")
            .regex(/[a-z]/, "Password must contain at least one lowercase letter")
            .regex(/\d/, "Password must contain at least one number")
            .regex(/[@$!%*?&]/, "Password must contain at least one special character (@, $, !, %, *, ?, &)"),
    })


type RegisterForm = z.infer<typeof registerSchema>

export function Register() {
    const navigate = useNavigate()
    const { login } = useAuthStore()
    const { t } = useLanguageStore()
    const [loading, setLoading] = useState(false)

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<RegisterForm>({
        resolver: zodResolver(registerSchema),
    })

    const onSubmit = async (data: RegisterForm) => {
        setLoading(true)
        try {
            const response = await api.post("/register", data)
            console.log("Register response:", response.data)  // Debug: check structure
            if (response.data.success) {
                const { token, user } = response.data.data  // Fix: data.data
                login(token, user)
                await fetchCsrfToken()  // Refresh CSRF after auth
                toast.success(t("auth.welcome"))
                navigate("/")
            } else {
                toast.error(response.data.message || "Registration failed")
            }
        } catch (error: any) {
            console.error("Register error:", error)  // Debug
            toast.error(error.response?.data?.error || "Registration failed")
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-primary/10 via-background to-accent/10 p-4">
            <Card className="w-full max-w-md">
                <CardHeader>
                    <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-primary">
                        <span className="text-3xl">🐾</span>
                    </div>
                    <CardTitle className="text-center text-2xl">{t("auth.register")}</CardTitle>
                </CardHeader>
                <CardContent>
                    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                        <div className="grid gap-4 sm:grid-cols-2">
                            <div>
                                <label className="mb-2 block text-sm font-medium text-foreground">{t("auth.firstName")}</label>
                                <Input {...register("firstName")} placeholder="John" />
                                {errors.firstName && <p className="mt-1 text-sm text-red-500">{errors.firstName.message}</p>}
                            </div>
                            <div>
                                <label className="mb-2 block text-sm font-medium text-foreground">{t("auth.lastName")}</label>
                                <Input {...register("lastName")} placeholder="Doe" />
                                {errors.lastName && <p className="mt-1 text-sm text-red-500">{errors.lastName.message}</p>}
                            </div>
                        </div>
                        <div>
                            <label className="mb-2 block text-sm font-medium text-foreground">{t("auth.email")}</label>
                            <Input {...register("email")} type="email" placeholder="you@example.com" />
                            {errors.email && <p className="mt-1 text-sm text-red-500">{errors.email.message}</p>}
                        </div>
                        <div>
                            <label className="mb-2 block text-sm font-medium text-foreground">{t("auth.password")}</label>
                            <Input {...register("password")} type="password" placeholder="••••••••" />
                            {errors.password && <p className="mt-1 text-sm text-red-500">{errors.password.message}</p>}
                        </div>
                        <Button type="submit" variant="default" className="w-full" disabled={loading}>
                            {loading ? t("common.loading") : t("auth.register")}
                        </Button>
                    </form>
                    <p className="mt-4 text-center text-sm text-muted-foreground">
                        Already have an account?{" "}
                        <Link to="/login" className="font-medium text-primary hover:underline">
                            {t("auth.login")}
                        </Link>
                    </p>
                </CardContent>
            </Card>
        </div>
    )
}