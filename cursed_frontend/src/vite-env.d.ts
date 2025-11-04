/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_API_URL: string
    readonly VITE_BASE_API_URL: string
    // можешь добавить сюда и другие переменные
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}
