import { ref } from 'vue'
import { useDashboardAssets } from './useDashboardAssets'
import { useDashboardProfile } from './useDashboardProfile'
import { useDashboardReminders } from './useDashboardReminders'
import { useDashboardTransactions } from './useDashboardTransactions'
import { createStatusState, clearStatus, setError } from './dashboardShared'

export function useDashboardData(userStore) {
    const currency = ref('try')

    const assets = useDashboardAssets(currency)
    const transactions = useDashboardTransactions(currency)
    const reminders = useDashboardReminders()
    const profile = useDashboardProfile(userStore)

    const pageStatus = createStatusState()
    const pageLoading = ref(false)

    async function bootstrap({ includeDashboardData = true } = {}) {
        try {
            pageLoading.value = true
            clearStatus(pageStatus)

            await profile.fetchProfile()

            if (includeDashboardData) {
                await Promise.all([
                    assets.fetchAssets(),
                    transactions.fetchTransactions(),
                    reminders.fetchReminders(),
                ])
            }
        } catch (err) {
            setError(pageStatus, err)
            if (err?.response?.status === 401) {
                userStore.logout()
            }
        } finally {
            pageLoading.value = false
        }
    }

    async function onCurrencyChange(nextCurrency) {
        currency.value = nextCurrency
        await Promise.all([assets.fetchAssets(), transactions.fetchTransactions()])
    }

    async function createTransaction() {
        await transactions.createTransaction({ onAfterCreate: assets.fetchAssets })
    }

    return {
        currency,
        assets,
        transactions,
        reminders,
        profile,
        pageLoading,
        pageStatus,
        bootstrap,
        onCurrencyChange,
        createTransaction,
    }
}