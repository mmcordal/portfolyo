<template>
  <AppCard title="İşlemler" subtitle="Ekleme, filtreleme ve raporlama">
    <StatusBanner type="error" :message="status.error" />
    <StatusBanner type="ok" :message="status.ok" />
    <div class="grid two-col">
      <form class="inline-form" @submit.prevent="$emit('create')">
        <h3>Yeni İşlem</h3>
        <select v-model="txForm.type" required>
          <option v-for="a in actionTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
        </select>
        <select v-model="txForm.asset" required>
          <option value="">Varlık seçin</option>
          <option v-for="a in assetTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
        </select>
        <input v-model.number="txForm.amount" type="number" min="0.0001" step="0.0001" placeholder="Miktar" required />
        <input v-model="txForm.transaction_date" type="datetime-local" />
        <input v-model="txForm.description" type="text" placeholder="Açıklama" />
        <button :disabled="loadingTx">{{ loadingTx ? 'Kaydediliyor...' : 'Kaydet' }}</button>
      </form>

      <div>
        <h3>Filtreler & Dışa Aktar</h3>
        <div class="toolbar">
          <select :value="txAssetFilter" @change="$emit('asset-filter-change', $event.target.value)">
            <option value="">Tüm Varlıklar</option>
            <option v-for="a in assetTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>
          <button @click="$emit('refresh')">İşlemleri Yenile</button>
          <button @click="$emit('download-excel')">Excel İndir</button>
        </div>
      </div>
    </div>

    <table v-if="transactions.length">
      <thead><tr><th>ID</th><th>Tip</th><th>Miktar</th><th>Tarih</th><th>Kur</th><th>Toplam</th><th></th></tr></thead>
      <tbody>
      <tr v-for="tx in transactions" :key="tx.id">
        <td>{{ tx.id }}</td>
        <td>{{ tx.type }}</td>
        <td>{{ formatNumber(tx.amount, 4) }}</td>
        <td>{{ formatDate(tx.transaction_date) }}</td>
        <td>{{ formatNumber(tx.target_currency_price) }}</td>
        <td>{{ formatNumber(tx.target_currency_total_price) }} {{ tx.now_target_currency?.toUpperCase() }}</td>
        <td><button @click="$emit('download-pdf', tx.id)">PDF</button></td>
      </tr>
      </tbody>
    </table>
    <p v-else class="subtle">İşlem bulunamadı.</p>
  </AppCard>
</template>

<script setup>
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'
import { formatDate, formatNumber } from '../../utils/format'

defineProps({
  txForm: { type: Object, required: true },
  txAssetFilter: { type: String, required: true },
  transactions: { type: Array, required: true },
  actionTypes: { type: Array, required: true },
  assetTypes: { type: Array, required: true },
  loadingTx: { type: Boolean, default: false },
  status: { type: Object, default: () => ({ ok: '', error: '' }) },
})

defineEmits(['create', 'asset-filter-change', 'refresh', 'download-excel', 'download-pdf'])
</script>

<style scoped>
h3 { margin: 0 0 .55rem; font-size: .95rem; }
.subtle { color: var(--color-muted); }
</style>