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

    <section class="summary-cards" aria-label="Dashboard özetleri">
      <article class="summary-card">
        <p>Toplam Portföy</p>
        <strong>
          {{ formatNumber(assets.assetsAll.value?.total_price || 0) }}
          {{ assets.assetsAll.value?.currency?.toUpperCase() || '' }}
        </strong>
      </article>
      <article class="summary-card">
        <p>Varlık Sayısı</p>
        <strong>{{ assets.assetsAll.value?.assets?.length || 0 }}</strong>
      </article>
      <article class="summary-card">
        <p>Toplam İşlem</p>
        <strong>{{ transactions.transactions.length }}</strong>
      </article>
      <article class="summary-card">
        <p>Yaklaşan Hatırlatıcı</p>
        <strong>{{ upcomingReminderCount }}</strong>
      </article>
    </section>

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
import { computed, onMounted } from 'vue'
import { ACTION_TYPES, ASSET_TYPES, CURRENCIES } from '../constants/assets'
import { useDashboardData } from '../composables/useDashboardData'
import { useUserStore } from '../stores/user'
import { formatNumber, toISODateTimeLocal } from '../utils/format'
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

const upcomingReminderCount = computed(() => reminders.reminders.value.filter((item) => {
  const timestamp = new Date(item.date).getTime()
  return Number.isFinite(timestamp) && timestamp >= Date.now()
}).length)

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
.summary-cards {
  display: grid;
  gap: .65rem;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}
.summary-card {
  border: 1px solid #d7e4f8;
  background: linear-gradient(180deg, #ffffff, #f7fbff);
  border-radius: 12px;
  padding: .72rem .8rem;
  box-shadow: var(--shadow-soft);
}
.summary-card p {
  margin: 0;
  color: #64748b;
  font-size: .76rem;
  font-weight: 600;
  letter-spacing: .02em;
}
.summary-card strong {
  display: inline-block;
  margin-top: .28rem;
  font-size: 1.03rem;
  color: #0f2f73;
}
.layout {
  display: grid;
  gap: 1.05rem;
}
</style>