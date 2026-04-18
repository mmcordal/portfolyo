<template>
  <main class="page">
    <div class="hero">
      <div>
        <h1>Finans Paneli</h1>
        <p>Tüm varlıklarınızı tek grafikte izleyin, işlemleri ve hatırlatıcıları yönetin.</p>
      </div>
      <div class="hero-badges">
        <span>{{ transactions.transactions.length }} işlem</span>
        <span>{{ reminders.reminders.length }} hatırlatıcı</span>
      </div>
    </div>

    <StatusBanner type="error" :message="pageStatus.error" />

    <div class="layout">
      <AssetsPanel :assets="assets" :currencies="CURRENCIES" @currency-change="onCurrencyChange" />

      <TransactionsPanel
          :transactions-domain="transactions"
          :action-types="ACTION_TYPES"
          :asset-types="ASSET_TYPES"
          @create="createTransaction"
      />

      <RemindersPanel :reminders-domain="reminders" />
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
  assets,
  transactions,
  reminders,
  pageStatus,
  bootstrap,
  onCurrencyChange,
  createTransaction,
} = useDashboardData(userStore)

onMounted(async () => {
  reminders.form.date = toISODateTimeLocal(new Date())
  await bootstrap()
})
</script>

<style scoped>
.hero {
  margin: .2rem 0 .1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: .85rem;
  flex-wrap: wrap;
}
.hero h1 {
  margin: 0;
  font-size: 1.46rem;
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
  border: 1px solid #cad8f2;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: var(--shadow-soft);
  border-radius: 999px;
  padding: .34rem .7rem;
  font-size: .78rem;
  color: #1e3a8a;
  font-weight: 600;
}
.layout {
  display: grid;
  gap: 1.05rem;
}
</style>