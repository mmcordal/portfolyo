import { apiRequest } from '../api'

export const authService = {
    register(payload) {
        return apiRequest('/auth/register', { method: 'POST', body: payload })
    },
    login(payload) {
        return apiRequest('/auth/login', { method: 'POST', body: payload })
    },
    me() {
        return apiRequest('/users/me')
    },
    updateMe(payload) {
        return apiRequest('/users/me', { method: 'PUT', body: payload })
    },
    deleteMe() {
        return apiRequest('/users/me', { method: 'DELETE' })
    },
}

export const assetService = {
    getAll(currency) {
        return apiRequest('/assets/all', { headers: { 'X-Currency': currency } })
    },
    getByAsset(asset, currency) {
        return apiRequest(`/assets/${asset}`, { headers: { 'X-Currency': currency } })
    },
    downloadPDF(currency) {
        return apiRequest('/assets/pdf', { headers: { 'X-Currency': currency }, responseType: 'blob' })
    },
}

export const transactionService = {
    create(payload) {
        return apiRequest('/transactions/', { method: 'POST', body: payload })
    },
    getAll(currency) {
        return apiRequest('/transactions/all', { headers: { 'X-Currency': currency } })
    },
    getByAsset(asset, currency) {
        return apiRequest(`/transactions/${asset}`, { headers: { 'X-Currency': currency } })
    },
    downloadExcel(currency) {
        return apiRequest('/transactions/excel', { headers: { 'X-Currency': currency }, responseType: 'blob' })
    },
    downloadPDF(txId, currency) {
        return apiRequest(`/transactions/pdf/${txId}`, { headers: { 'X-Currency': currency }, responseType: 'blob' })
    },
}

export const reminderService = {
    create(payload) {
        return apiRequest('/reminders/', { method: 'POST', body: payload })
    },
    getAll() {
        return apiRequest('/reminders/')
    },
    remove(id) {
        return apiRequest(`/reminders/${id}`, { method: 'DELETE' })
    },
}

export function downloadBlob(blob, fileName) {
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = fileName
    document.body.appendChild(a)
    a.click()
    a.remove()
    window.URL.revokeObjectURL(url)
}