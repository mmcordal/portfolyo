const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api/v1'

function getAuthHeaders() {
    const token = localStorage.getItem('jwt')
    return token ? { Authorization: `Bearer ${token}` } : {}
}

export async function apiRequest(path, options = {}) {
    const { method = 'GET', body, headers = {}, responseType = 'json' } = options

    const requestInit = {
        method,
        headers: {
            ...getAuthHeaders(),
            ...headers,
        },
    }

    if (body !== undefined) {
        requestInit.body = JSON.stringify(body)
        requestInit.headers['Content-Type'] = 'application/json'
    }

    const response = await fetch(`${API_BASE}${path}`, requestInit)

    if (responseType === 'blob') {
        if (!response.ok) {
            const text = await response.text()
            throw new Error(text || 'Dosya indirme hatası')
        }
        return await response.blob()
    }

    const data = await response.json().catch(() => ({}))

    if (!response.ok) {
        const message = data?.error_message || data?.message || 'Beklenmeyen hata'
        const error = new Error(message)
        error.response = { status: response.status, data }
        throw error
    }

    return data
}

export function extractErrorMessage(error) {
    return error?.response?.data?.error_message || error?.response?.data?.message || error?.message || 'Beklenmeyen bir hata oluştu.'
}