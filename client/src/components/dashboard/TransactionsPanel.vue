<template>
  <AppCard title="İşlemler" subtitle="Ekleme, filtreleme ve raporlama">
    <StatusBanner type="error" :message="transactionsDomain.status.error" />
    <StatusBanner type="ok" :message="transactionsDomain.status.ok" />

    <div class="grid two-col">
      <form class="inline-form" @submit.prevent="$emit('create')">
        <h3>Yeni İşlem</h3>
        <select v-model="transactionsDomain.txForm.type" required>
          <option v-for="a in actionTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
        </select>
        <select v-model="transactionsDomain.txForm.asset" required>
          <option value="">Varlık seçin</option>
          <option v-for="a in assetTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
        </select>
        <input v-model.number="transactionsDomain.txForm.amount" type="number" min="0.0001" step="0.0001" placeholder="Miktar" required />
        <input v-model="transactionsDomain.txForm.transaction_date" type="datetime-local" />
        <input v-model="transactionsDomain.txForm.description" type="text" placeholder="Açıklama" />
        <button :disabled="transactionsDomain.loading.value">{{ transactionsDomain.loading.value ? 'Kaydediliyor...' : 'Kaydet' }}</button>
      </form>

      <div>
        <h3>Filtreler & Dışa Aktar</h3>
        <div class="toolbar two-col">
          <select v-model="transactionsDomain.filters.asset">
            <option value="">Tüm Varlıklar</option>
            <option v-for="a in assetTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>

          <select v-model="transactionsDomain.filters.type">
            <option value="">Tüm Tipler</option>
            <option value="add">Ekle</option>
            <option value="subtract">Çıkar</option>
          </select>

          <input v-model="transactionsDomain.filters.dateFrom" type="date" />
          <input v-model="transactionsDomain.filters.dateTo" type="date" />

          <input
              v-model="transactionsDomain.filters.search"
              type="text"
              placeholder="ID / açıklama / varlık ara"
          />

          <select :value="transactionsDomain.filters.perPage" @change="transactionsDomain.setPerPage(Number($event.target.value))">
            <option :value="3">3 / sayfa</option>
            <option :value="6">6 / sayfa</option>
            <option :value="9">9 / sayfa</option>
            <option :value="12">12 / sayfa</option>
          </select>

          <button class="secondary" @click="transactionsDomain.fetchTransactions">İşlemleri Yenile</button>
          <button class="secondary" @click="transactionsDomain.downloadTxExcel">Excel İndir</button>
        </div>
      </div>
    </div>

    <p class="subtle" v-if="transactionsDomain.pagedTransactions.value.length">
      Toplam sonuç: {{ transactionsDomain.filteredTransactions.value.length }}
    </p>

    <table v-if="transactionsDomain.pagedTransactions.value.length">
      <thead>
      <tr><th>ID</th><th>Tip</th><th>Miktar</th><th>Oluşturulma</th><th>İşlem Tarihi</th><th>Kur</th><th>Toplam</th><th></th></tr>
      </thead>
      <tbody>
      <tr v-for="tx in transactionsDomain.pagedTransactions.value" :key="tx.id">
        <td>{{ tx.id }}</td>
        <td>{{ tx.type }}</td>
        <td>{{ formatNumber(tx.amount, 4) }}</td>
        <td>{{ formatDate(tx.created_at) }}</td>
        <td>{{ formatDate(tx.transaction_date) }}</td>
        <td>{{ formatNumber(tx.target_currency_price) }}</td>
        <td>{{ formatNumber(tx.target_currency_total_price) }} {{ tx.now_target_currency?.toUpperCase() }}</td>
        <td><button class="secondary" @click="transactionsDomain.downloadTxPdf(tx.id)">PDF</button></td>
      </tr>
      </tbody>
    </table>
    <p v-else class="subtle">İşlem bulunamadı.</p>

    <div class="pager" v-if="transactionsDomain.totalPages.value > 1">
      <button :disabled="transactionsDomain.filters.page <= 1" @click="transactionsDomain.setPage(transactionsDomain.filters.page - 1)">Önceki</button>
      <span>Sayfa {{ transactionsDomain.filters.page }} / {{ transactionsDomain.totalPages.value }}</span>
      <button :disabled="transactionsDomain.filters.page >= transactionsDomain.totalPages.value" @click="transactionsDomain.setPage(transactionsDomain.filters.page + 1)">Sonraki</button>
    </div>
  </AppCard>
</template>

<script setup>
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'
import { formatDate, formatNumber } from '../../utils/format'

defineProps({
  transactionsDomain: { type: Object, required: true },
  actionTypes: { type: Array, required: true },
  assetTypes: { type: Array, required: true },
})

defineEmits(['create'])
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