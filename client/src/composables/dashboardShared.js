import { reactive } from 'vue'
import { extractErrorMessage } from '../api'

export function createStatusState() {
    return reactive({ ok: '', error: '' })
}

export function clearStatus(status) {
    status.ok = ''
    status.error = ''
}

export function setError(status, err) {
    clearStatus(status)
    status.error = extractErrorMessage(err)
}

export function setOk(status, message) {
    clearStatus(status)
    status.ok = message
}

export function toTimestamp(value) {
    const parsed = new Date(value).getTime()
    return Number.isNaN(parsed) ? 0 : parsed
}