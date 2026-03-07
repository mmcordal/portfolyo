<template>
  <main class="page">
    <div class="hero">
      <h1>Finans Paneli</h1>
      <p>Tüm varlıklarınızı tek grafikte izleyin, işlemleri ve hatırlatıcıları yönetin.</p>
    </div>

    <StatusBanner type="error" :message="status.page.error" />

    <div class="layout">
      <AssetsPanel
          :assets-all="assetsAll"
          :single-asset="singleAsset"
          :currency="currency"
          :currencies="CURRENCIES"
          :status="status.assets"
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
          :status="status.tx"
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
          :status="status.reminders"
          @create="createReminder"
          @delete="deleteReminder"
      />
    </div>
  </main>
</template>

<script setup>
import { onMounted } from 'vue'
import { ACTION_TYPES, ASSET_TYPES, CURRENCIES } from '../constants/assets'
import { useDashboardData } from '../composables/useDashboardData'
import { useUserStore } from '../stores/user'
import { toISODateTimeLocal } from '../utils/format'
import AssetsPanel from '../components/dashboard/AssetsPanel.vue'
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
  loading,
  status,
  bootstrap,
  fetchSingleAsset,
  downloadAssetsPdf,
  createTransaction,
  fetchTransactions,
  downloadTxExcel,
  downloadTxPdf,
  createReminder,
  deleteReminder,
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
.layout {
  display: grid;
  gap: 1rem;
}
</style>