// src/stores/authStore.ts
import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import type { StateCreator } from "zustand";
import type { User } from "@/types/types.ts";

interface AuthState {
    user: User | null;
    token: string | null;
    isAuth: boolean;
    isLoading: boolean;  // Initial: false — no stuck
    login: (token: string, user: User) => void;
    logout: () => void;
    fetchUser: () => void;  // Keep for manual calls (e.g. after API refresh), but not in rehydrate
}

const createAuthSlice: StateCreator<AuthState> = (set, get) => ({
    user: null,
    token: null,
    isAuth: false,
    isLoading: false,  // ← FIX: false initially, no spinner flash

    login: (token: string, user: User) => {
        localStorage.setItem("token", token);
        localStorage.setItem("user", JSON.stringify(user));
        set({ token, user, isAuth: true, isLoading: false });
        console.log("Logged in, state updated");
    },

    logout: () => {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        localStorage.removeItem("auth-storage");  // Clear persist too
        set({ token: null, user: null, isAuth: false, isLoading: false });
        console.log("Logged out");
    },

    fetchUser: () => {
        const token = localStorage.getItem("token");
        const userStr = localStorage.getItem("user");
        console.log("fetchUser called, token exists:", !!token);
        if (token && userStr) {
            try {
                const user = JSON.parse(userStr) as User;
                set({ token, user, isAuth: true, isLoading: false });
                console.log("Rehydrated auth:", user.email);
            } catch (e) {
                console.error("Parse error:", e);
                get().logout();
            }
        } else {
            set({ isLoading: false });  // Redundant now, but safe
        }
    },
});

export const useAuthStore = create(
    persist(
        createAuthSlice,
        {
            name: "auth-storage",
            storage: createJSONStorage(() => localStorage),
            partialize: (state) => ({ token: state.token, user: state.user, isAuth: state.isAuth }),
            onRehydrateStorage: () => {
                return (state, error) => {
                    if (error) {
                        console.warn("Auth rehydrate error:", error);
                        // No fetchUser — persist already loaded partial state
                    } else {
                        console.log("Rehydrate complete");  // No fetchUser — duplicate & risky
                    }
                    // No manual set isLoading — initial false, persist handles
                };
            },
        }
    )
);