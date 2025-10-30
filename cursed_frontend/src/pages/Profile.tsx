import { CustomerNav } from "@/components/layout/CustomerNav"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { useAuthStore } from "@/stores/authStore"
import { useLanguageStore } from "@/stores/languageStore"
import { User } from "lucide-react"

export function Profile() {
  const { user } = useAuthStore()
  const { t } = useLanguageStore()

  if (!user) return null

  return (
    <div className="min-h-screen bg-background">
      <CustomerNav />
      <div className="mx-auto max-w-3xl px-4 py-12 sm:px-6 lg:px-8">
        <Card>
          <CardHeader>
            <CardTitle>{t("nav.profile")}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <div className="flex h-20 w-20 items-center justify-center rounded-full bg-primary text-primary-foreground">
                {user.image ? (
                  <img
                    src={user.image || "/placeholder.svg"}
                    alt={user.firstName}
                    className="h-full w-full rounded-full object-cover"
                  />
                ) : (
                  <User className="h-10 w-10" />
                )}
              </div>
              <div>
                <h2 className="text-2xl font-bold text-foreground">
                  {user.firstName} {user.lastName}
                </h2>
                <p className="text-muted-foreground">{user.email}</p>
                <p className="mt-1 text-sm text-muted-foreground">Role: {user.role}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
