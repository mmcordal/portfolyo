import { computed, reactive, ref, watch } from 'vue'
import { extractErrorMessage } from '../api'
import { authService, assetService, downloadBlob, reminderService, transactionService } from '../services/portfolio'

function createStatusState() {
    return reactive({ ok: '', error: '' })
}

function clearStatus(status) {
    status.ok = ''
    status.error = ''
}

function setError(status, err) {
    clearStatus(status)
    status.error = extractErrorMessage(err)
}

function setOk(status, message) {
    clearStatus(status)
    status.ok = message
}

function toTimestamp(value) {
    const parsed = new Date(value).getTime()
    return Number.isNaN(parsed) ? 0 : parsed
}

export function useDashboardData(userStore) {
    const currency = ref('try')
    const txAssetFilter = ref('')
    const txTypeFilter = ref('')
    const txSearch = ref('')
    const txDateFrom = ref('')
    const txDateTo = ref('')
    const txPage = ref(1)
    const txPerPage = ref(6)

    const assetsAll = ref(null)
    const singleAsset = ref(null)
    const transactions = ref([])
    const reminders = ref([])

    const txForm = reactive({ type: 'add', asset: '', amount: null, transaction_date: '', description: '' })
    const reminderForm = reactive({ title: '', date: '' })
    const profileForm = reactive({ name: '', surname: '', email: '', password: '' })

    const loading = reactive({ profile: false, tx: false, reminder: false, page: false })
    const status = reactive({
        page: createStatusState(),
        assets: createStatusState(),
        tx: createStatusState(),
        reminders: createStatusState(),
        profile: createStatusState(),
    })

    const sortedTransactions = computed(() => {
        return [...transactions.value].sort((a, b) => toTimestamp(b.created_at) - toTimestamp(a.created_at))
    })

    const filteredTransactions = computed(() => {
        const fromTs = txDateFrom.value ? toTimestamp(`${txDateFrom.value}T00:00:00`) : null
        const toTs = txDateTo.value ? toTimestamp(`${txDateTo.value}T23:59:59`) : null
        const search = txSearch.value.trim().toLowerCase()

        return sortedTransactions.value.filter((tx) => {
            if (txTypeFilter.value && String(tx.type).toLowerCase() !== txTypeFilter.value.toLowerCase()) {
                return false
            }

            const createdAtTs = toTimestamp(tx.created_at)
            if (fromTs !== null && createdAtTs < fromTs) return false
            if (toTs !== null && createdAtTs > toTs) return false

            if (!search) return true
            const haystack = [
                tx.id,
                tx.type,
                tx.description,
                tx.user_asset?.asset,
                tx.transaction_date,
            ].join(' ').toLowerCase()
            return haystack.includes(search)
        })
    })

    const totalTxPages = computed(() => {
        return Math.max(1, Math.ceil(filteredTransactions.value.length / txPerPage.value))
    })

    const pagedTransactions = computed(() => {
        const start = (txPage.value - 1) * txPerPage.value
        return filteredTransactions.value.slice(start, start + txPerPage.value)
    })

    watch([filteredTransactions, txPerPage], () => {
        if (txPage.value > totalTxPages.value) {
            txPage.value = totalTxPages.value
        }
        if (txPage.value < 1) {
            txPage.value = 1
        }
    })

    function resetTxPagination() {
        txPage.value = 1
    }

    async function bootstrap() {
        try {
            loading.page = true
            clearStatus(status.page)
            const meRes = await authService.me()
            userStore.setProfile(meRes.data)
            profileForm.email = meRes.data?.email || ''
            await Promise.all([fetchAssets(), fetchTransactions(), fetchReminders()])
        } catch (err) {
            setError(status.page, err)
            if (err?.response?.status === 401) {
                userStore.logout()
            }
        } finally {
            loading.page = false
        }
    }

    async function fetchAssets() {
        try {
            const res = await assetService.getAll(currency.value)
            assetsAll.value = res.data
        } catch (err) {
            setError(status.assets, err)
        }
    }

    async function fetchSingleAsset(asset) {
        try {
            const res = await assetService.getByAsset(asset, currency.value)
            singleAsset.value = res.data
            setOk(status.assets, `${asset.toUpperCase()} detayı getirildi.`)
        } catch (err) {
            setError(status.assets, err)
        }
    }

    async function downloadAssetsPdf() {
        try {
            const blob = await assetService.downloadPDF(currency.value)
            downloadBlob(blob, `portfolio-${currency.value}.pdf`)
            setOk(status.assets, 'Portföy PDF indirildi.')
        } catch (err) {
            setError(status.assets, err)
        }
    }

    async function createTransaction() {
        try {
            loading.tx = true
            clearStatus(status.tx)
            await transactionService.create({
                ...txForm,
                transaction_date: txForm.transaction_date ? new Date(txForm.transaction_date).toISOString() : '',
            })
            txForm.asset = ''
            txForm.amount = null
            txForm.transaction_date = ''
            txForm.description = ''
            setOk(status.tx, 'İşlem eklendi.')
            await Promise.all([fetchAssets(), fetchTransactions()])
        } catch (err) {
            setError(status.tx, err)
        } finally {
            loading.tx = false
        }
    }

    async function fetchTransactions() {
        try {
            const res = txAssetFilter.value
                ? await transactionService.getByAsset(txAssetFilter.value, currency.value)
                : await transactionService.getAll(currency.value)
            transactions.value = res.data || []
            resetTxPagination()
        } catch (err) {
            setError(status.tx, err)
        }
    }

    async function downloadTxExcel() {
        try {
            const blob = await transactionService.downloadExcel(currency.value)
            downloadBlob(blob, `transactions-${currency.value}.xlsx`)
            setOk(status.tx, 'İşlemler Excel olarak indirildi.')
        } catch (err) {
            setError(status.tx, err)
        }
    }

    async function downloadTxPdf(txId) {
        try {
            const blob = await transactionService.downloadPDF(txId, currency.value)
            downloadBlob(blob, `transaction-${txId}.pdf`)
            setOk(status.tx, `İşlem #${txId} PDF indirildi.`)
        } catch (err) {
            setError(status.tx, err)
        }
    }

    async function fetchReminders() {
        try {
            const res = await reminderService.getAll()
            reminders.value = res.data || []
        } catch (err) {
            setError(status.reminders, err)
        }
    }

    async function createReminder() {
        try {
            loading.reminder = true
            clearStatus(status.reminders)
            await reminderService.create({
                title: reminderForm.title,
                date: reminderForm.date ? new Date(reminderForm.date).toISOString() : '',
            })
            reminderForm.title = ''
            reminderForm.date = ''
            setOk(status.reminders, 'Hatırlatıcı eklendi.')
            await fetchReminders()
        } catch (err) {
            setError(status.reminders, err)
        } finally {
            loading.reminder = false
        }
    }

    async function deleteReminder(id) {
        try {
            clearStatus(status.reminders)
            await reminderService.remove(id)
            setOk(status.reminders, 'Hatırlatıcı silindi.')
            await fetchReminders()
        } catch (err) {
            setError(status.reminders, err)
        }
    }

    async function updateProfile() {
        try {
            loading.profile = true
            clearStatus(status.profile)
            const payload = Object.fromEntries(Object.entries(profileForm).filter(([, v]) => String(v || '').trim() !== ''))
            if (!Object.keys(payload).length) {
                setOk(status.profile, 'Güncellenecek alan bulunamadı.')
                return
            }
            await authService.updateMe(payload)
            setOk(status.profile, 'Profil güncellendi.')
            Object.assign(profileForm, { name: '', surname: '', email: payload.email || profileForm.email, password: '' })
            const meRes = await authService.me()
            userStore.setProfile(meRes.data)
            profileForm.email = meRes.data?.email || ''
        } catch (err) {
            setError(status.profile, err)
        } finally {
            loading.profile = false
        }
    }

    async function deleteProfile() {
        if (!window.confirm('Hesabınızı silmek istediğinize emin misiniz?')) return
        try {
            loading.profile = true
            clearStatus(status.profile)
            await authService.deleteMe()
            userStore.logout()
        } catch (err) {
            setError(status.profile, err)
        } finally {
            loading.profile = false
        }
    }

    return {
        currency,
        txAssetFilter,
        txTypeFilter,
        txSearch,
        txDateFrom,
        txDateTo,
        txPage,
        txPerPage,
        assetsAll,
        singleAsset,
        transactions,
        filteredTransactions,
        pagedTransactions,
        totalTxPages,
        reminders,
        txForm,
        reminderForm,
        profileForm,
        loading,
        status,
        bootstrap,
        fetchAssets,
        fetchSingleAsset,
        downloadAssetsPdf,
        createTransaction,
        fetchTransactions,
        downloadTxExcel,
        downloadTxPdf,
        fetchReminders,
        createReminder,
        deleteReminder,
        updateProfile,
        deleteProfile,
    }
}