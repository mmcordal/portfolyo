import api from '../api'

export const authService = {
    register(payload) {
        return api.post('/auth/register', payload)
    },
    login(payload) {
        return api.post('/auth/login', payload)
    },
    me() {
        return api.get('/users/me')
    },
    updateMe(payload) {
        return api.put('/users/me', payload)
    },
    deleteMe() {
        return api.delete('/users/me')
    },
}

export const assetService = {
    getAll(currency) {
        return api.get('/assets/all', { headers: { 'X-Currency': currency } })
    },
    getByAsset(asset, currency) {
        return api.get(`/assets/${asset}`, { headers: { 'X-Currency': currency } })
    },
    downloadPDF(currency) {
        return api.get('/assets/pdf', {
            headers: { 'X-Currency': currency },
            responseType: 'blob',
        })
    },
}

export const transactionService = {
    create(payload) {
        return api.post('/transactions/', payload)
    },
    getAll(currency) {
        return api.get('/transactions/all', { headers: { 'X-Currency': currency } })
    },
    getByAsset(asset, currency) {
        return api.get(`/transactions/${asset}`, { headers: { 'X-Currency': currency } })
    },
    downloadExcel(currency) {
        return api.get('/transactions/excel', {
            headers: { 'X-Currency': currency },
            responseType: 'blob',
        })
    },
    downloadPDF(txId, currency) {
        return api.get(`/transactions/pdf/${txId}`, {
            headers: { 'X-Currency': currency },
            responseType: 'blob',
        })
    },
}

export const reminderService = {
    create(payload) {
        return api.post('/reminders/', payload)
    },
    getAll() {
        return api.get('/reminders/')
    },
    remove(id) {
        return api.delete(`/reminders/${id}`)
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
