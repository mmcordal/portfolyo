<template>
  <main class="page">
    <div class="hero">
      <h1>Finans Paneli</h1>
      <p>Portföy, işlemler ve hatırlatıcılar tek ekranda.</p>
    </div>

    <StatusBanner type="error" :message="error" />
    <StatusBanner type="ok" :message="ok" />

    <ProfilePanel
        :full-name="userStore.fullName"
        :email="userStore.profile?.email"
        :form="profileForm"
        :loading="loading.profile"
        @update="updateProfile"
        @delete="deleteProfile"
    />

    <AssetsPanel
        :assets-all="assetsAll"
        :single-asset="singleAsset"
        :currency="currency"
        :currencies="CURRENCIES"
        @download-pdf="downloadAssetsPdf"
        @currency-change="onCurrencyChange"
        @show-asset="fetchSingleAsset"
    />

    <TransactionsPanel
        :tx-form="txForm"
        :tx-asset-filter="txAssetFilter"
        :transactions="transactions"
        :action-types="ACTION_TYPES"
        :asset-types="ASSET_TYPES"
        :loading-tx="loading.tx"
        @create="createTransaction"
        @asset-filter-change="onAssetFilterChange"
        @refresh="fetchTransactions"
        @download-excel="downloadTxExcel"
        @download-pdf="downloadTxPdf"
    />

    <RemindersPanel
        :reminder-form="reminderForm"
        :reminders="reminders"
        :loading-reminder="loading.reminder"
        @create="createReminder"
        @delete="deleteReminder"
    />
  </main>
</template>

<script setup>
import { onMounted } from 'vue'
import { ACTION_TYPES, ASSET_TYPES, CURRENCIES } from '../constants/assets'
import { useDashboardData } from '../composables/useDashboardData'
import { useUserStore } from '../stores/user'
import { toISODateTimeLocal } from '../utils/format'
import AssetsPanel from '../components/dashboard/AssetsPanel.vue'
import ProfilePanel from '../components/dashboard/ProfilePanel.vue'
import RemindersPanel from '../components/dashboard/RemindersPanel.vue'
import TransactionsPanel from '../components/dashboard/TransactionsPanel.vue'
import StatusBanner from '../components/ui/StatusBanner.vue'

const userStore = useUserStore()

const {
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
  bootstrap,
  fetchSingleAsset,
  downloadAssetsPdf,
  createTransaction,
  fetchTransactions,
  downloadTxExcel,
  downloadTxPdf,
  createReminder,
  deleteReminder,
  updateProfile,
  deleteProfile,
  fetchAssets,
} = useDashboardData(userStore)

function onCurrencyChange(next) {
  currency.value = next
  fetchAssets()
  fetchTransactions()
}

function onAssetFilterChange(next) {
  txAssetFilter.value = next
  fetchTransactions()
}

onMounted(async () => {
  reminderForm.date = toISODateTimeLocal(new Date())
  await bootstrap()
})
</script>

<style scoped>
.hero {
  margin: .25rem 0 .2rem;
}
.hero h1 {
  margin: 0;
  font-size: 1.4rem;
}
.hero p {
  margin: .2rem 0 0;
  color: var(--color-muted);
}
</style>