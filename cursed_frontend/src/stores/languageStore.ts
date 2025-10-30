import { create } from "zustand"
import { translations } from "@/utils/translations"

type Language = "en" | "ru" | "kz"

interface LanguageState {
  language: Language
  setLanguage: (lang: Language) => void
  t: (key: string) => string
}

export const useLanguageStore = create<LanguageState>((set, get) => ({
  language: "en",

  setLanguage: (language: Language) => {
    localStorage.setItem("language", language)
    set({ language })
  },

  t: (key: string) => {
    const { language } = get()
    return translations[language][key] || key
  },
}))
