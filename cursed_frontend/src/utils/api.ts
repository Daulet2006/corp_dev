// src/utils/api.ts
import axios, { AxiosInstance } from "axios";
import { useAuthStore } from "@/stores/authStore";

const api: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_URL || "http://localhost:8080/api",
    headers: { "Content-Type": "application/json" },
    withCredentials: true,  // Для CSRF cookie
});

let csrfToken: string = localStorage.getItem("csrf_token") || '';  // Persist, default empty string

// Fetch CSRF token (call after login or on init)
export const fetchCsrfToken = async (): Promise<string> => {
    try {
        console.log('Fetching CSRF...');
        const { data } = await api.get("/csrf-token");
        console.log('CSRF response:', { csrf_token: !!data.csrf_token });
        const newToken = data.csrf_token || '';
        csrfToken = newToken;
        localStorage.setItem("csrf_token", csrfToken);
        return csrfToken;
    } catch (error) {
        console.error("CSRF fetch failed:", error);
        csrfToken = '';
        localStorage.removeItem("csrf_token");
        return "";
    }
};

// Request interceptor: Add JWT + CSRF
api.interceptors.request.use(async (config) => {
    const { token } = useAuthStore.getState();  // From Zustand
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }

    // Skip CSRF for public paths (login/register)
    const publicPaths = ['/login', '/register'];
    const isPublicMutate = publicPaths.some(path => config.url?.includes(path)) &&
        ["post", "put", "patch", "delete"].includes(config.method || "");

    const needsCsrf = !isPublicMutate && ["post", "put", "patch", "delete"].includes(config.method || "");

    if (needsCsrf && !csrfToken) {
        const fetchedToken = await fetchCsrfToken();
        if (fetchedToken) {
            csrfToken = fetchedToken;
        }
    }

    if (needsCsrf && csrfToken) {
        config.headers["X-CSRF-Token"] = csrfToken;
    }
    return config;
});

// Response interceptor: Handle 401/403, refresh CSRF on fail
api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401) {
            useAuthStore.getState().logout();
            window.location.href = "/login";
            return Promise.reject(error);
        }

        // Auto-retry CSRF once, mutate originalRequest to avoid loops
        if (
            error.response?.status === 403 &&
            error.response?.data?.error?.includes("CSRF") &&
            !originalRequest?._retry &&
            !originalRequest.url?.includes('/login') &&
            !originalRequest.url?.includes('/register')
        ) {
            originalRequest._retry = true;
            try {
                const fetchedToken = await fetchCsrfToken();
                if (fetchedToken) {
                    csrfToken = fetchedToken;
                }
                return api(originalRequest);  // Retry original
            } catch (retryError) {
                console.error("CSRF retry failed:", retryError);
            }
        }

        if (error.response?.status === 429) {
            console.warn("Rate limited—retry in 1s");
            return new Promise((resolve) => setTimeout(() => resolve(api(originalRequest)), 1000));
        }

        return Promise.reject(error);
    }
);

export default api;