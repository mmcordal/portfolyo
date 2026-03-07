import axios from 'axios'

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api/v1'

const api = axios.create({
    baseURL: API_BASE,
    timeout: 15000,
})

api.interceptors.request.use((config) => {
    const token = localStorage.getItem('jwt')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

export function extractErrorMessage(error) {
    return (
        error?.response?.data?.error_message ||
        error?.response?.data?.message ||
        error?.message ||
        'Beklenmeyen bir hata oluştu.'
    )
}

export default api