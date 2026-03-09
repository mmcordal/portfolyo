<template>
  <main class="page">
    <div class="hero">
      <div>
        <h1>Finans Paneli</h1>
        <p>Tüm varlıklarınızı tek grafikte izleyin, işlemleri ve hatırlatıcıları yönetin.</p>
      </div>
      <div class="hero-badges">
        <span>{{ transactions.length }} işlem</span>
        <span>{{ reminders.length }} hatırlatıcı</span>
      </div>
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
          :tx-type-filter="txTypeFilter"
          :tx-search="txSearch"
          :tx-date-from="txDateFrom"
          :tx-date-to="txDateTo"
          :tx-page="txPage"
          :tx-per-page="txPerPage"
          :total-tx-pages="totalTxPages"
          :total-filtered-count="filteredTransactions.length"
          :transactions="pagedTransactions"
          :action-types="ACTION_TYPES"
          :asset-types="ASSET_TYPES"
          :loading-tx="loading.tx"
          :status="status.tx"
          @create="createTransaction"
          @asset-filter-change="onAssetFilterChange"
          @type-filter-change="onTypeFilterChange"
          @search-change="onSearchChange"
          @date-from-change="onDateFromChange"
          @date-to-change="onDateToChange"
          @per-page-change="onPerPageChange"
          @page-change="onPageChange"
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
  txPage.value = 1
  fetchTransactions()
}

function onTypeFilterChange(next) {
  txTypeFilter.value = next
  txPage.value = 1
}

function onSearchChange(next) {
  txSearch.value = next
  txPage.value = 1
}

function onDateFromChange(next) {
  txDateFrom.value = next
  txPage.value = 1
}

function onDateToChange(next) {
  txDateTo.value = next
  txPage.value = 1
}

function onPerPageChange(next) {
  txPerPage.value = Math.min(12, Math.max(3, Number(next) || 6))
  txPage.value = 1
}

function onPageChange(next) {
  txPage.value = Math.min(totalTxPages.value, Math.max(1, Number(next) || 1))
}

onMounted(async () => {
  reminderForm.date = toISODateTimeLocal(new Date())
  await bootstrap()
})
</script>

<style scoped>
.hero {
  margin: .25rem 0 .2rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: .7rem;
  flex-wrap: wrap;
}
.hero h1 {
  margin: 0;
  font-size: 1.4rem;
}
.hero p {
  margin: .2rem 0 0;
  color: var(--color-muted);
}
.hero-badges {
  display: flex;
  gap: .45rem;
}
.hero-badges span {
  border: 1px solid rgba(148, 163, 184, 0.35);
  background: rgba(15, 23, 42, 0.55);
  border-radius: 999px;
  padding: .3rem .65rem;
  font-size: .78rem;
  color: #dbeafe;
}
.layout {
  display: grid;
  gap: 1rem;
}
</style>