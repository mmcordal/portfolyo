<template>
  <AppCard title="Portföy" subtitle="Varlık dağılımını güçlü görsel hiyerarşi ile izleyin">
    <StatusBanner type="error" :message="assets.status.error" />
    <StatusBanner type="ok" :message="assets.status.ok" />

    <p class="summary" v-if="assets.assetsAll.value">
      Toplam: {{ formatNumber(assets.assetsAll.value.total_price) }} {{ assets.assetsAll.value.currency?.toUpperCase() }}
    </p>

    <div class="charts" v-if="hasAssets">
      <div class="chart-box donut-box">
        <h3>Dağılım (Donut)</h3>
        <div class="donut-wrap">
          <Pie :data="pieData" :options="donutOptions" />
          <div class="donut-center">
            <span>Toplam</span>
            <strong>{{ formatNumber(assets.assetsAll.value?.total_price || 0) }}</strong>
            <small>{{ assets.assetsAll.value?.currency?.toUpperCase() }}</small>
          </div>
        </div>
      </div>
      <div class="chart-box">
        <h3>Varlık Değerleri</h3>
        <Bar :data="barData" :options="barOptions" />
      </div>
    </div>

    <section class="holding-cards" v-if="hasAssets">
      <article class="holding-card" v-for="asset in assets.assetsAll.value.assets" :key="asset.id">
        <div class="holding-head">
          <strong>{{ asset.asset.toUpperCase() }}</strong>
          <span>{{ formatRatio(asset) }}%</span>
        </div>
        <p>{{ formatNumber(asset.amount, 4) }} adet · {{ formatNumber(asset.price) }} birim</p>
        <p class="holding-total">{{ formatNumber(asset.total_price_by_asset) }} {{ assets.assetsAll.value.currency?.toUpperCase() }}</p>
        <div class="progress-track">
          <div class="progress-fill" :style="{ width: `${Math.min(100, formatRatio(asset))}%` }" />
        </div>
      </article>
    </section>

    <table v-if="hasAssets">
      <thead><tr><th>Varlık</th><th>Miktar</th><th>Birim</th><th>Toplam</th><th></th></tr></thead>
      <tbody>
      <tr v-for="asset in assets.assetsAll.value.assets" :key="asset.id">
        <td>{{ asset.asset.toUpperCase() }}</td>
        <td>{{ formatNumber(asset.amount, 4) }}</td>
        <td>{{ formatNumber(asset.price) }}</td>
        <td>{{ formatNumber(asset.total_price_by_asset) }}</td>
        <td><button class="secondary" @click="assets.fetchSingleAsset(asset.asset)">Detay</button></td>
      </tr>
      </tbody>
    </table>
    <div v-else class="empty-state">
      <h3>Henüz portföy verisi yok</h3>
      <p>Yeni bir işlem eklediğinizde varlık dağılımı burada görünecek.</p>
    </div>

    <div class="asset-detail" v-if="assets.singleAsset.value">
      <h3>{{ assets.singleAsset.value.asset.toUpperCase() }} Detayı</h3>
      <p>{{ formatNumber(assets.singleAsset.value.total_price_by_asset) }} {{ assets.singleAsset.value.target_currency?.toUpperCase() }}</p>
      <ul class="list" v-if="assets.singleAsset.value.transactions?.length">
        <li v-for="tx in assets.singleAsset.value.transactions" :key="tx.id">
          <span>{{ tx.type }}</span>
          <span>{{ formatNumber(tx.amount, 4) }}</span>
          <span>{{ formatDate(tx.transaction_date) }}</span>
        </li>
      </ul>
    </div>
  </AppCard>
</template>

<script setup>
import { computed } from 'vue'
import { Bar, Pie } from 'vue-chartjs'
import {
  ArcElement,
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  Title,
  Tooltip,
} from 'chart.js'
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'
import { formatDate, formatNumber } from '../../utils/format'

ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale, BarElement)

const props = defineProps({
  assets: { type: Object, required: true },
  currencies: { type: Array, required: true },
})

const chartAssets = computed(() => props.assets.assetsAll.value?.assets || [])
const hasAssets = computed(() => chartAssets.value.length > 0)
const labels = computed(() => chartAssets.value.map((asset) => asset.asset.toUpperCase()))
const totals = computed(() => chartAssets.value.map((asset) => Number(asset.total_price_by_asset || 0)))
const colorPalette = ['#3b82f6', '#38bdf8', '#f59e0b', '#34d399', '#f472b6', '#6366f1', '#fb7185']

const pieData = computed(() => ({
  labels: labels.value,
  datasets: [{
    label: `Toplam (${props.assets.assetsAll.value?.currency?.toUpperCase() || ''})`,
    data: totals.value,
    backgroundColor: labels.value.map((_, index) => colorPalette[index % colorPalette.length]),
    borderColor: '#ffffff',
    borderWidth: 2,
    hoverOffset: 6,
  }],
}))

const barData = computed(() => ({
  labels: labels.value,
  datasets: [{
    label: `Varlık Değeri (${props.assets.assetsAll.value?.currency?.toUpperCase() || ''})`,
    data: totals.value,
    backgroundColor: 'rgba(59, 130, 246, 0.75)',
    borderRadius: 8,
  }],
}))

function extractTooltipNumericValue(context) {
  if (typeof context.parsed === 'number') return context.parsed
  if (typeof context.parsed?.y === 'number') return context.parsed.y
  if (typeof context.raw === 'number') return context.raw
  return 0
}

function tooltipLabel(context) {
  const value = formatNumber(extractTooltipNumericValue(context))
  const code = props.assets.assetsAll.value?.currency?.toUpperCase() || ''
  return `${context.label}: ${value} ${code}`
}

function formatRatio(asset) {
  const total = Number(props.assets.assetsAll.value?.total_price || 0)
  if (total <= 0) return 0
  return Number(((Number(asset.total_price_by_asset || 0) / total) * 100).toFixed(2))
}

const donutOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  cutout: '67%',
  layout: {
    padding: {
      top: 4,
      right: 8,
      bottom: 4,
      left: 8,
    },
  },
  plugins: {
    legend: {
      position: 'bottom',
      labels: {
        color: '#334155',
        boxWidth: 12,
        boxHeight: 12,
        padding: 14,
      },
    },
    tooltip: {
      callbacks: {
        label: tooltipLabel,
      },
    },
  },
  scales: {
    x: { display: false, grid: { display: false }, ticks: { display: false } },
    y: { display: false, grid: { display: false }, ticks: { display: false } },
  },
}))

const barOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom',
      labels: {
        color: '#334155',
        boxWidth: 12,
        boxHeight: 12,
        padding: 10,
      },
    },
    tooltip: {
      callbacks: {
        label: tooltipLabel,
      },
    },
  },
  scales: {
    x: { ticks: { color: '#64748b' }, grid: { color: 'rgba(148, 163, 184, 0.16)' } },
    y: { ticks: { color: '#64748b' }, grid: { color: 'rgba(148, 163, 184, 0.22)' } },
  },
}))
</script>

<style scoped>
.asset-detail {
  margin-top: .8rem;
  padding-top: .8rem;
  border-top: 1px dashed var(--color-surface-light);
}
.asset-detail h3 { margin: 0 0 .4rem; }
.asset-detail p { margin: 0 0 .6rem; }
.charts {
  display: grid;
  gap: .85rem;
  margin-bottom: .85rem;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
}
.chart-box {
  background: #ffffff;
  border: 1px solid #dde7f7;
  border-radius: 12px;
  padding: .68rem .68rem .55rem;
}
.donut-box {
  display: grid;
  align-content: start;
}
.donut-wrap {
  position: relative;
}
.donut-center {
  position: absolute;
  inset: 0;
  display: grid;
  place-content: center;
  text-align: center;
  pointer-events: none;
}
.donut-center span {
  color: #64748b;
  font-size: .74rem;
}
.donut-center strong {
  color: #0f2f73;
  font-size: .95rem;
}
.donut-center small {
  color: #1e40af;
  font-weight: 700;
}
.chart-box h3 {
  margin: 0 0 .45rem;
  font-size: .9rem;
}
.chart-box :deep(canvas) {
  height: 220px !important;
}
.holding-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: .6rem;
  margin-bottom: .85rem;
}
.holding-card {
  border: 1px solid #dce7f9;
  border-radius: 12px;
  background: #ffffff;
  padding: .65rem;
}
.holding-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: .2rem;
}
.holding-card p {
  color: #64748b;
  font-size: .82rem;
  margin-top: .18rem;
}
.holding-total {
  font-size: .92rem !important;
  color: #0f2f73 !important;
  font-weight: 700;
}
.progress-track {
  margin-top: .45rem;
  width: 100%;
  height: 8px;
  border-radius: 999px;
  background: #e7efff;
  overflow: hidden;
}
.progress-fill {
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, #2563eb, #60a5fa);
}
.empty-state {
  border: 1px dashed #c9d9f4;
  border-radius: 12px;
  padding: 1rem;
  background: #f9fbff;
}
.empty-state h3 {
  margin: 0;
  font-size: 1rem;
}
.empty-state p {
  margin-top: .3rem;
  color: #64748b;
}
.list li {
  display: grid;
  grid-template-columns: .8fr .7fr 1fr;
  gap: .5rem;
}
</style>