<template>
  <AppCard title="Portföy" subtitle="Varlık dağılımınızı grafik ve tabloyla takip edin">
    <template #actions>
      <button class="secondary" @click="$emit('download-pdf')">Portföy PDF</button>
    </template>

    <StatusBanner type="error" :message="status.error" />
    <StatusBanner type="ok" :message="status.ok" />

    <div class="toolbar">
      <label>Hedef Para Birimi</label>
      <select :value="currency" @change="$emit('currency-change', $event.target.value)">
        <option v-for="c in currencies" :key="c" :value="c">{{ c.toUpperCase() }}</option>
      </select>
    </div>

    <p class="summary" v-if="assetsAll">Toplam: {{ formatNumber(assetsAll.total_price) }} {{ assetsAll.currency?.toUpperCase() }}</p>

    <div class="charts" v-if="hasAssets">
      <div class="chart-box">
        <h3>Dağılım (Pasta)</h3>
        <Pie :data="pieData" :options="chartOptions" />
      </div>
      <div class="chart-box">
        <h3>Varlık Değerleri (Çubuk)</h3>
        <Bar :data="barData" :options="chartOptions" />
      </div>
    </div>

    <table v-if="hasAssets">
      <thead><tr><th>Varlık</th><th>Miktar</th><th>Birim</th><th>Toplam</th><th></th></tr></thead>
      <tbody>
      <tr v-for="asset in assetsAll.assets" :key="asset.id">
        <td>{{ asset.asset.toUpperCase() }}</td>
        <td>{{ formatNumber(asset.amount, 4) }}</td>
        <td>{{ formatNumber(asset.price) }}</td>
        <td>{{ formatNumber(asset.total_price_by_asset) }}</td>
        <td><button class="secondary" @click="$emit('show-asset', asset.asset)">Detay</button></td>
      </tr>
      </tbody>
    </table>
    <p v-else class="subtle">Henüz varlık kaydı yok.</p>

    <div class="asset-detail" v-if="singleAsset">
      <h3>{{ singleAsset.asset.toUpperCase() }} Detayı</h3>
      <p>{{ formatNumber(singleAsset.total_price_by_asset) }} {{ singleAsset.target_currency?.toUpperCase() }}</p>
      <ul class="list" v-if="singleAsset.transactions?.length">
        <li v-for="tx in singleAsset.transactions" :key="tx.id">
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
  assetsAll: { type: Object, default: null },
  singleAsset: { type: Object, default: null },
  currency: { type: String, required: true },
  currencies: { type: Array, required: true },
  status: { type: Object, default: () => ({ ok: '', error: '' }) },
})

defineEmits(['download-pdf', 'currency-change', 'show-asset'])

const chartAssets = computed(() => props.assetsAll?.assets || [])
const hasAssets = computed(() => chartAssets.value.length > 0)
const labels = computed(() => chartAssets.value.map((asset) => asset.asset.toUpperCase()))
const totals = computed(() => chartAssets.value.map((asset) => Number(asset.total_price_by_asset || 0)))
const colorPalette = ['#3b82f6', '#22d3ee', '#f59e0b', '#34d399', '#f472b6', '#a78bfa', '#f87171']

const pieData = computed(() => ({
  labels: labels.value,
  datasets: [{
    label: `Toplam (${props.assetsAll?.currency?.toUpperCase() || ''})`,
    data: totals.value,
    backgroundColor: labels.value.map((_, index) => colorPalette[index % colorPalette.length]),
    borderColor: '#0f172a',
    borderWidth: 1,
  }],
}))

const barData = computed(() => ({
  labels: labels.value,
  datasets: [{
    label: `Varlık Değeri (${props.assetsAll?.currency?.toUpperCase() || ''})`,
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

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { labels: { color: '#dbeafe' } },
    tooltip: {
      callbacks: {
        label(context) {
          const value = formatNumber(extractTooltipNumericValue(context))
          const code = props.assetsAll?.currency?.toUpperCase() || ''
          return `${context.label}: ${value} ${code}`
        },
      },
    },
  },
  scales: {
    x: { ticks: { color: '#cbd5e1' }, grid: { color: 'rgba(148, 163, 184, 0.15)' } },
    y: { ticks: { color: '#cbd5e1' }, grid: { color: 'rgba(148, 163, 184, 0.15)' } },
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
.subtle { color: var(--color-muted); }
.charts {
  display: grid;
  gap: .85rem;
  margin-bottom: .85rem;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
}
.chart-box {
  background: rgba(2, 6, 23, 0.45);
  border: 1px solid var(--color-surface-light);
  border-radius: 12px;
  padding: .6rem;
}
.chart-box h3 {
  margin: 0 0 .45rem;
  font-size: .9rem;
}
.chart-box :deep(canvas) {
  height: 220px !important;
}
.list li {
  display: grid;
  grid-template-columns: .8fr .7fr 1fr;
  gap: .5rem;
}
</style>