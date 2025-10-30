// Updated authStore.ts — add isLoading for sync rehydration
import { create } from "zustand"
import type { User } from "@/types/types.ts"

interface AuthState {
    user: User | null
    token: string | null
    isAuth: boolean
    isLoading: boolean  // New: loading during fetchUser
    login: (token: string, user: User) => void
    logout: () => void
    fetchUser: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
    user: null,
    token: null,
    isAuth: false,
    isLoading: true,  // Start as loading

    login: (token: string, user: User) => {
        localStorage.setItem("token", token)
        localStorage.setItem("user", JSON.stringify(user))
        set({ token, user, isAuth: true, isLoading: false })
    },

    logout: () => {
        localStorage.removeItem("token")
        localStorage.removeItem("user")
        set({ token: null, user: null, isAuth: false, isLoading: false })
    },

    fetchUser: () => {
        const token = localStorage.getItem("token")
        const userStr = localStorage.getItem("user")
        if (token && userStr) {
            try {
                const user = JSON.parse(userStr) as User
                set({ token, user, isAuth: true, isLoading: false })
            } catch (e) {
                // Invalid JSON — treat as logout
                localStorage.removeItem("token")
                localStorage.removeItem("user")
                set({ token: null, user: null, isAuth: false, isLoading: false })
            }
        } else {
            set({ isLoading: false })
        }
    },
}))