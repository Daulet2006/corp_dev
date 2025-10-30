import { create } from "zustand"

interface ThemeState {
  theme: "light" | "dark"
  toggleTheme: () => void
  setTheme: (theme: "light" | "dark") => void
}

export const useThemeStore = create<ThemeState>((set) => ({
  theme: "light",

  toggleTheme: () =>
    set((state) => {
      const newTheme = state.theme === "light" ? "dark" : "light"
      document.documentElement.classList.toggle("dark", newTheme === "dark")
      localStorage.setItem("theme", newTheme)
      return { theme: newTheme }
    }),

  setTheme: (theme) => {
    document.documentElement.classList.toggle("dark", theme === "dark")
    localStorage.setItem("theme", theme)
    set({ theme })
  },
}))
