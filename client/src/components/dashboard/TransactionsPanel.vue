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
        <div class="toolbar two-col">
          <select :value="txAssetFilter" @change="$emit('asset-filter-change', $event.target.value)">
            <option value="">Tüm Varlıklar</option>
            <option v-for="a in assetTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>

          <select :value="txTypeFilter" @change="$emit('type-filter-change', $event.target.value)">
            <option value="">Tüm Tipler</option>
            <option value="add">Ekle</option>
            <option value="subtract">Çıkar</option>
          </select>

          <input :value="txDateFrom" type="date" @input="$emit('date-from-change', $event.target.value)" />
          <input :value="txDateTo" type="date" @input="$emit('date-to-change', $event.target.value)" />

          <input
              :value="txSearch"
              type="text"
              placeholder="ID / açıklama / varlık ara"
              @input="$emit('search-change', $event.target.value)"
          />

          <select :value="txPerPage" @change="$emit('per-page-change', Number($event.target.value))">
            <option :value="3">3 / sayfa</option>
            <option :value="6">6 / sayfa</option>
            <option :value="9">9 / sayfa</option>
            <option :value="12">12 / sayfa</option>
          </select>

          <button class="secondary" @click="$emit('refresh')">İşlemleri Yenile</button>
          <button class="secondary" @click="$emit('download-excel')">Excel İndir</button>
        </div>
      </div>
    </div>

    <p class="subtle" v-if="transactions.length">Toplam sonuç: {{ totalFilteredCount }}</p>

    <table v-if="transactions.length">
      <thead>
      <tr><th>ID</th><th>Tip</th><th>Miktar</th><th>Oluşturulma</th><th>İşlem Tarihi</th><th>Kur</th><th>Toplam</th><th></th></tr>
      </thead>
      <tbody>
      <tr v-for="tx in transactions" :key="tx.id">
        <td>{{ tx.id }}</td>
        <td>{{ tx.type }}</td>
        <td>{{ formatNumber(tx.amount, 4) }}</td>
        <td>{{ formatDate(tx.created_at) }}</td>
        <td>{{ formatDate(tx.transaction_date) }}</td>
        <td>{{ formatNumber(tx.target_currency_price) }}</td>
        <td>{{ formatNumber(tx.target_currency_total_price) }} {{ tx.now_target_currency?.toUpperCase() }}</td>
        <td><button class="secondary" @click="$emit('download-pdf', tx.id)">PDF</button></td>
      </tr>
      </tbody>
    </table>
    <p v-else class="subtle">İşlem bulunamadı.</p>

    <div class="pager" v-if="totalTxPages > 1">
      <button :disabled="txPage <= 1" @click="$emit('page-change', txPage - 1)">Önceki</button>
      <span>Sayfa {{ txPage }} / {{ totalTxPages }}</span>
      <button :disabled="txPage >= totalTxPages" @click="$emit('page-change', txPage + 1)">Sonraki</button>
    </div>
  </AppCard>
</template>

<script setup>
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'
import { formatDate, formatNumber } from '../../utils/format'

defineProps({
  txForm: { type: Object, required: true },
  txAssetFilter: { type: String, required: true },
  txTypeFilter: { type: String, required: true },
  txSearch: { type: String, required: true },
  txDateFrom: { type: String, required: true },
  txDateTo: { type: String, required: true },
  txPage: { type: Number, required: true },
  txPerPage: { type: Number, required: true },
  totalTxPages: { type: Number, required: true },
  totalFilteredCount: { type: Number, required: true },
  transactions: { type: Array, required: true },
  actionTypes: { type: Array, required: true },
  assetTypes: { type: Array, required: true },
  loadingTx: { type: Boolean, default: false },
  status: { type: Object, default: () => ({ ok: '', error: '' }) },
})

defineEmits([
  'create',
  'asset-filter-change',
  'type-filter-change',
  'search-change',
  'date-from-change',
  'date-to-change',
  'per-page-change',
  'page-change',
  'refresh',
  'download-excel',
  'download-pdf',
])
</script>

<style scoped>
h3 { margin: 0 0 .55rem; font-size: .95rem; }
.subtle { color: var(--color-muted); margin: .4rem 0; }
.pager {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: .5rem;
  margin-top: .6rem;
}
</style>