// Updated api.ts — use navigate from react-router for redirect (inject via param or use global, but for simplicity, keep window.location; or import useNavigate if in hook)
import axios from "axios"
import { useAuthStore } from "@/stores/authStore"
import toast from "react-hot-toast"

const api = axios.create({
    baseURL: "http://localhost:8080/api",
    headers: {
        "Content-Type": "application/json",
    },
})

// Request interceptor
api.interceptors.request.use(
    (config) => {
        const token = useAuthStore.getState().token
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    },
)

// Response interceptor
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401 || error.response?.status === 403) {
            useAuthStore.getState().logout()
            toast.error("Session expired. Please login again.")
            // Better: if in router context, use navigate; for now, window.location
            window.location.href = "/login"
        }
        return Promise.reject(error)
    },
)

export default api