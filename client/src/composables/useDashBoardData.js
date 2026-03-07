import { reactive, ref } from 'vue'
import { extractErrorMessage } from '../api'
import { authService, assetService, downloadBlob, reminderService, transactionService } from '../services/portfolio'

export function useDashboardData(userStore) {
    const currency = ref('try')
    const txAssetFilter = ref('')

    const assetsAll = ref(null)
    const singleAsset = ref(null)
    const transactions = ref([])
    const reminders = ref([])

    const txForm = reactive({ type: 'add', asset: '', amount: null, transaction_date: '', description: '' })
    const reminderForm = reactive({ title: '', date: '' })
    const profileForm = reactive({ name: '', surname: '', email: '', password: '' })

    const loading = reactive({ profile: false, tx: false, reminder: false, page: false })
    const error = ref('')
    const ok = ref('')

    function resetMessage() {
        error.value = ''
        ok.value = ''
    }

    async function bootstrap() {
        try {
            loading.page = true
            const meRes = await authService.me()
            userStore.setProfile(meRes.data)
            profileForm.email = meRes.data?.email || ''
            await Promise.all([fetchAssets(), fetchTransactions(), fetchReminders()])
        } catch (err) {
            error.value = extractErrorMessage(err)
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
            error.value = extractErrorMessage(err)
        }
    }

    async function fetchSingleAsset(asset) {
        try {
            const res = await assetService.getByAsset(asset, currency.value)
            singleAsset.value = res.data
        } catch (err) {
            error.value = extractErrorMessage(err)
        }
    }

    async function downloadAssetsPdf() {
        try {
            const blob = await assetService.downloadPDF(currency.value)
            downloadBlob(blob, `portfolio-${currency.value}.pdf`)
            ok.value = 'Portföy PDF indirildi.'
        } catch (err) {
            error.value = extractErrorMessage(err)
        }
    }

    async function createTransaction() {
        try {
            loading.tx = true
            resetMessage()
            await transactionService.create({
                ...txForm,
                transaction_date: txForm.transaction_date ? new Date(txForm.transaction_date).toISOString() : '',
            })
            txForm.asset = ''
            txForm.amount = null
            txForm.transaction_date = ''
            txForm.description = ''
            ok.value = 'İşlem eklendi.'
            await Promise.all([fetchAssets(), fetchTransactions()])
        } catch (err) {
            error.value = extractErrorMessage(err)
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
        } catch (err) {
            error.value = extractErrorMessage(err)
        }
    }

    async function downloadTxExcel() {
        try {
            const blob = await transactionService.downloadExcel(currency.value)
            downloadBlob(blob, `transactions-${currency.value}.xlsx`)
            ok.value = 'İşlemler Excel olarak indirildi.'
        } catch (err) {
            error.value = extractErrorMessage(err)
        }
    }

    async function downloadTxPdf(txId) {
        try {
            const blob = await transactionService.downloadPDF(txId, currency.value)
            downloadBlob(blob, `transaction-${txId}.pdf`)
            ok.value = `İşlem #${txId} PDF indirildi.`
        } catch (err) {
            error.value = extractErrorMessage(err)
        }
    }

    async function fetchReminders() {
        try {
            const res = await reminderService.getAll()
            reminders.value = res.data || []
        } catch (err) {
            error.value = extractErrorMessage(err)
        }
    }

    async function createReminder() {
        try {
            loading.reminder = true
            resetMessage()
            await reminderService.create({
                title: reminderForm.title,
                date: reminderForm.date ? new Date(reminderForm.date).toISOString() : '',
            })
            reminderForm.title = ''
            reminderForm.date = ''
            ok.value = 'Hatırlatıcı eklendi.'
            await fetchReminders()
        } catch (err) {
            error.value = extractErrorMessage(err)
        } finally {
            loading.reminder = false
        }
    }

    async function deleteReminder(id) {
        try {
            resetMessage()
            await reminderService.remove(id)
            ok.value = 'Hatırlatıcı silindi.'
            await fetchReminders()
        } catch (err) {
            error.value = extractErrorMessage(err)
        }
    }

    async function updateProfile() {
        try {
            loading.profile = true
            resetMessage()
            const payload = Object.fromEntries(Object.entries(profileForm).filter(([, v]) => String(v || '').trim() !== ''))
            if (!Object.keys(payload).length) {
                ok.value = 'Güncellenecek alan bulunamadı.'
                return
            }
            await authService.updateMe(payload)
            ok.value = 'Profil güncellendi.'
            Object.assign(profileForm, { name: '', surname: '', email: payload.email || profileForm.email, password: '' })
            const meRes = await authService.me()
            userStore.setProfile(meRes.data)
            profileForm.email = meRes.data?.email || ''
        } catch (err) {
            error.value = extractErrorMessage(err)
        } finally {
            loading.profile = false
        }
    }

    async function deleteProfile() {
        if (!window.confirm('Hesabınızı silmek istediğinize emin misiniz?')) return
        try {
            loading.profile = true
            resetMessage()
            await authService.deleteMe()
            userStore.logout()
        } catch (err) {
            error.value = extractErrorMessage(err)
        } finally {
            loading.profile = false
        }
    }

    return {
        currency,
        txAssetFilter,
        assetsAll,
        singleAsset,
        transactions,
        reminders,
        txForm,
        reminderForm,
        profileForm,
        loading,
        error,
        ok,
        resetMessage,
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