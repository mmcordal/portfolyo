<template>
  <main class="page dashboard-page">
    <section class="dashboard-hero">
      <div class="hero-main">
        <p class="eyebrow">Portföy Komuta Merkezi</p>
        <h1>Finans Paneli</h1>
        <p>
          Portföy performansınızı tek bakışta görün, işlemleri yönetin ve hatırlatıcılarınızı aksatmadan takip edin.
        </p>

        <article class="hero-primary-metric">
          <span>Toplam Portföy</span>
          <strong>
            {{ formatNumber(assets.assetsAll.value?.total_price || 0) }}
            {{ assets.assetsAll.value?.currency?.toUpperCase() || '' }}
          </strong>
        </article>

        <div class="hero-metrics">
          <article>
            <span>Varlık Sayısı</span>
            <strong>{{ assets.assetsAll.value?.assets?.length || 0 }}</strong>
          </article>
          <article>
            <span>Toplam İşlem</span>
            <strong>{{ transactions.transactions.length }}</strong>
          </article>
          <article>
            <span>Yaklaşan Hatırlatıcı</span>
            <strong>{{ upcomingReminderCount }}</strong>
          </article>
        </div>
      </div>

      <div class="hero-actions-panel">
        <label>Hedef Para Birimi</label>
        <select :value="assets.currency.value" @change="onCurrencyChange($event.target.value)">
          <option v-for="c in CURRENCIES" :key="c" :value="c">{{ c.toUpperCase() }}</option>
        </select>

        <button class="secondary" @click="assets.downloadAssetsPdf">Portföy PDF</button>
        <button @click="isTxModalOpen = true">+ Yeni İşlem</button>
        <button class="secondary" @click="isReminderModalOpen = true">+ Hatırlatıcı</button>

        <div class="hero-badges">
          <span>{{ assets.assetsAll.value?.assets?.length || 0 }} varlık</span>
          <span>{{ reminders.reminders.value.length }} hatırlatıcı</span>
        </div>
      </div>
    </section>

    <section class="summary-cards" aria-label="Dashboard özetleri">
      <article class="summary-card">
        <p><span class="card-icon">🪙</span> Varlık Sayısı</p>
        <strong>{{ assets.assetsAll.value?.assets?.length || 0 }}</strong>
      </article>
      <article class="summary-card">
        <p><span class="card-icon">🧾</span> Toplam İşlem</p>
        <strong>{{ transactions.transactions.length }}</strong>
      </article>
      <article class="summary-card">
        <p><span class="card-icon">⏰</span> Yaklaşan Hatırlatıcı</p>
        <strong>{{ upcomingReminderCount }}</strong>
      </article>
      <article class="summary-card">
        <p><span class="card-icon">🌍</span> Aktif Para Birimi</p>
        <strong>{{ assets.currency.value.toUpperCase() }}</strong>
      </article>
    </section>

    <StatusBanner type="error" :message="pageStatus.error" />
    <StatusBanner type="ok" :message="pageStatus.ok" />

    <section v-if="pageLoading" class="page-loading" aria-live="polite">
      <strong>Dashboard hazırlanıyor...</strong>
      <p>Varlıklar, işlemler ve hatırlatıcılar yükleniyor.</p>
    </section>

    <div class="layout" v-else>
      <AssetsPanel :assets="assets" :currencies="CURRENCIES" @currency-change="onCurrencyChange" />

      <TransactionsPanel
          :transactions-domain="transactions"
          :action-types="ACTION_TYPES"
          :asset-types="ASSET_TYPES"
          :open-create="isTxModalOpen"
          @create="handleCreateTransaction"
          @close-create="isTxModalOpen = false"
      />

      <RemindersPanel
          :reminders-domain="reminders"
          :open-create="isReminderModalOpen"
          @close-create="isReminderModalOpen = false"
      />
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { ACTION_TYPES, ASSET_TYPES, CURRENCIES } from '../constants/assets'
import { useDashboardData } from '../composables/useDashboardData'
import { useUserStore } from '../stores/user'
import { formatNumber, toISODateTimeLocal } from '../utils/format'
import AssetsPanel from '../components/dashboard/AssetsPanel.vue'
import RemindersPanel from '../components/dashboard/RemindersPanel.vue'
import TransactionsPanel from '../components/dashboard/TransactionsPanel.vue'
import StatusBanner from '../components/ui/StatusBanner.vue'

const userStore = useUserStore()
const isTxModalOpen = ref(false)
const isReminderModalOpen = ref(false)

const {
  assets,
  transactions,
  reminders,
  pageStatus,
  pageLoading,
  bootstrap,
  onCurrencyChange,
  createTransaction,
} = useDashboardData(userStore)

const upcomingReminderCount = computed(() => reminders.reminders.value.filter((item) => {
  const timestamp = new Date(item.date).getTime()
  return Number.isFinite(timestamp) && timestamp >= Date.now()
}).length)

async function handleCreateTransaction(onSuccessClose) {
  await createTransaction()
  if (!transactions.status.error && typeof onSuccessClose === 'function') {
    onSuccessClose()
  }
}

onMounted(async () => {
  reminders.form.date = toISODateTimeLocal(new Date())
  await bootstrap()
})
</script>

<style scoped>
.dashboard-page {
  gap: 1rem;
}
.dashboard-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(260px, 1fr);
  gap: .9rem;
  border: 1px solid #d5e3fb;
  border-radius: 18px;
  padding: 1rem;
  background: linear-gradient(145deg, #ffffff, #f3f8ff);
  box-shadow: var(--shadow-card);
}
.eyebrow {
  display: inline-flex;
  margin-bottom: .3rem;
  padding: .2rem .55rem;
  border-radius: 999px;
  background: #e5efff;
  color: #1e40af;
  font-size: .72rem;
  font-weight: 700;
}
.hero-main h1 {
  margin: 0;
  font-size: 1.62rem;
}
.hero-main > p {
  margin-top: .28rem;
  color: var(--color-muted);
  max-width: 65ch;
}
.hero-primary-metric {
  margin-top: .9rem;
  border: 1px solid #d6e3fa;
  border-radius: 14px;
  padding: .72rem;
  background: #ffffff;
}
.hero-primary-metric span {
  font-size: .76rem;
  color: #64748b;
}
.hero-primary-metric strong {
  display: block;
  margin-top: .24rem;
  color: #0f2f73;
  font-size: 1.2rem;
}
.hero-metrics {
  margin-top: .65rem;
  display: grid;
  gap: .58rem;
  grid-template-columns: repeat(auto-fit, minmax(165px, 1fr));
}
.hero-metrics article {
  border: 1px solid #d6e3fa;
  border-radius: 12px;
  padding: .62rem .7rem;
  background: #ffffff;
}
.hero-metrics span {
  display: block;
  color: #64748b;
  font-size: .76rem;
}
.hero-metrics strong {
  display: block;
  margin-top: .24rem;
  color: #0f2f73;
  font-size: 1.02rem;
}
.hero-actions-panel {
  border: 1px solid #d7e3f8;
  border-radius: 14px;
  padding: .75rem;
  background: #ffffff;
  display: grid;
  gap: .52rem;
  align-content: start;
}
.hero-actions-panel .hero-badges {
  display: flex;
  flex-wrap: wrap;
  gap: .4rem;
  margin-top: .1rem;
}
.hero-badges span {
  border: 1px solid #cad8f2;
  background: rgba(255, 255, 255, 0.86);
  border-radius: 999px;
  padding: .3rem .62rem;
  font-size: .76rem;
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
  display: flex;
  gap: .35rem;
  align-items: center;
}
.card-icon { font-size: .95rem; }
.summary-card strong {
  display: inline-block;
  margin-top: .28rem;
  font-size: 1.03rem;
  color: #0f2f73;
}
.page-loading {
  border: 1px solid #cddbf5;
  border-radius: 12px;
  padding: .9rem;
  background: #f8fbff;
}
.page-loading p {
  margin-top: .25rem;
  color: var(--color-muted);
}
.layout {
  display: grid;
  gap: 1.05rem;
}
@media (max-width: 900px) {
  .dashboard-hero {
    grid-template-columns: 1fr;
  }
}
</style>