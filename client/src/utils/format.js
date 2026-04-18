export function formatNumber(value, fraction = 2) {
    const num = Number(value || 0)
    return new Intl.NumberFormat('tr-TR', {
        minimumFractionDigits: fraction,
        maximumFractionDigits: fraction,
    }).format(num)
}

export function formatDate(value, options = {}) {
    if (!value) return '-'
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return value

    const {
        includeTime = true,
        includeSeconds = false,
    } = options

    return date.toLocaleString('tr-TR', {
        dateStyle: 'short',
        ...(includeTime
            ? {
                timeStyle: includeSeconds ? 'medium' : 'short',
            }
            : {}),
    })
}

export function toISODateTimeLocal(value) {
    if (!value) return ''
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return ''
    const offset = date.getTimezoneOffset()
    const local = new Date(date.getTime() - offset * 60000)
    return local.toISOString().slice(0, 16)
}