import { computed, reactive, ref, watch } from 'vue'
import { transactionService, downloadBlob } from '../services/portfolio'
import { createStatusState, clearStatus, setError, setOk, toTimestamp } from './dashboardShared'

export function useDashboardTransactions(currency) {
    const transactions = ref([])
    const loading = ref(false)
    const status = createStatusState()

    const txForm = reactive({
        type: 'add',
        asset: '',
        amount: null,
        transaction_date: '',
        description: '',
    })

    const filters = reactive({
        asset: '',
        type: '',
        search: '',
        dateFrom: '',
        dateTo: '',
        page: 1,
        perPage: 6,
    })

    const sortedTransactions = computed(() => {
        return [...transactions.value].sort((a, b) => toTimestamp(b.created_at) - toTimestamp(a.created_at))
    })

    const filteredTransactions = computed(() => {
        const fromTs = filters.dateFrom ? toTimestamp(`${filters.dateFrom}T00:00:00`) : null
        const toTs = filters.dateTo ? toTimestamp(`${filters.dateTo}T23:59:59`) : null
        const search = filters.search.trim().toLowerCase()

        return sortedTransactions.value.filter((tx) => {
            if (filters.type && String(tx.type).toLowerCase() !== filters.type.toLowerCase()) {
                return false
            }

            const createdAtTs = toTimestamp(tx.created_at)
            if (fromTs !== null && createdAtTs < fromTs) return false
            if (toTs !== null && createdAtTs > toTs) return false

            if (!search) return true
            const haystack = [tx.id, tx.type, tx.description, tx.user_asset?.asset, tx.transaction_date]
                .join(' ')
                .toLowerCase()
            return haystack.includes(search)
        })
    })

    const totalPages = computed(() => {
        return Math.max(1, Math.ceil(filteredTransactions.value.length / filters.perPage))
    })

    const pagedTransactions = computed(() => {
        const start = (filters.page - 1) * filters.perPage
        return filteredTransactions.value.slice(start, start + filters.perPage)
    })

    watch([filteredTransactions, () => filters.perPage], () => {
        if (filters.page > totalPages.value) {
            filters.page = totalPages.value
        }
        if (filters.page < 1) {
            filters.page = 1
        }
    })

    watch(
        () => filters.asset,
        async () => {
            filters.page = 1
            await fetchTransactions()
        },
    )

    watch(
        [() => filters.type, () => filters.search, () => filters.dateFrom, () => filters.dateTo],
        () => {
            filters.page = 1
        },
    )

    function setPerPage(value) {
        filters.perPage = Math.min(12, Math.max(3, Number(value) || 6))
        filters.page = 1
    }

    function setPage(value) {
        filters.page = Math.min(totalPages.value, Math.max(1, Number(value) || 1))
    }

    async function fetchTransactions() {
        try {
            const res = filters.asset
                ? await transactionService.getByAsset(filters.asset, currency.value)
                : await transactionService.getAll(currency.value)
            transactions.value = res.data || []
            filters.page = 1
        } catch (err) {
            setError(status, err)
        }
    }

    async function createTransaction({ onAfterCreate } = {}) {
        try {
            loading.value = true
            clearStatus(status)
            await transactionService.create({
                ...txForm,
                transaction_date: txForm.transaction_date ? new Date(txForm.transaction_date).toISOString() : '',
            })

            txForm.asset = ''
            txForm.amount = null
            txForm.transaction_date = ''
            txForm.description = ''
            setOk(status, 'İşlem eklendi.')

            await fetchTransactions()
            if (onAfterCreate) {
                await onAfterCreate()
            }
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    async function downloadTxExcel() {
        try {
            const blob = await transactionService.downloadExcel(currency.value)
            downloadBlob(blob, `transactions-${currency.value}.xlsx`)
            setOk(status, 'İşlemler Excel olarak indirildi.')
        } catch (err) {
            setError(status, err)
        }
    }

    async function downloadTxPdf(txId) {
        try {
            const blob = await transactionService.downloadPDF(txId, currency.value)
            downloadBlob(blob, `transaction-${txId}.pdf`)
            setOk(status, `İşlem #${txId} PDF indirildi.`)
        } catch (err) {
            setError(status, err)
        }
    }

    return {
        transactions,
        filteredTransactions,
        pagedTransactions,
        totalPages,
        txForm,
        filters,
        loading,
        status,
        setPerPage,
        setPage,
        fetchTransactions,
        createTransaction,
        downloadTxExcel,
        downloadTxPdf,
    }
}