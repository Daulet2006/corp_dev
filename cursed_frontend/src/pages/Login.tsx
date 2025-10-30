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
import api from "@/utils/api"

const loginSchema = z.object({
  email: z.string().email("Invalid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
})

type LoginForm = z.infer<typeof loginSchema>

export function Login() {
  const navigate = useNavigate()
  const { login } = useAuthStore()
  const { t } = useLanguageStore()
  const [loading, setLoading] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data: LoginForm) => {
    setLoading(true)
    try {
      const response = await api.post("/login", data)
      login(response.data.token, response.data.user)
      toast.success(t("auth.welcome"))
      navigate("/")
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Login failed")
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
          <CardTitle className="text-center text-2xl">{t("auth.login")}</CardTitle>
          <p className="text-center text-sm text-muted-foreground">{t("auth.loginPrompt")}</p>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
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
            <Button type="submit" variant="primary" className="w-full" disabled={loading}>
              {loading ? t("common.loading") : t("auth.login")}
            </Button>
          </form>
          <p className="mt-4 text-center text-sm text-muted-foreground">
            Don't have an account?{" "}
            <Link to="/register" className="font-medium text-primary hover:underline">
              {t("auth.register")}
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  )
}
