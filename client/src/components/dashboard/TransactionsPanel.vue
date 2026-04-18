<template>
  <AppCard title="İşlemler" subtitle="Ekleme, filtreleme ve raporlama">
    <StatusBanner type="error" :message="transactionsDomain.status.error" />
    <StatusBanner type="ok" :message="transactionsDomain.status.ok" />

    <div class="grid two-col panels">
      <form class="inline-form transaction-form" @submit.prevent="$emit('create')">
        <h3>Yeni İşlem</h3>
        <div class="form-grid">
          <select v-model="transactionsDomain.txForm.type" required>
            <option v-for="a in actionTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>
          <select v-model="transactionsDomain.txForm.asset" required>
            <option value="">Varlık seçin</option>
            <option v-for="a in assetTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>
          <input v-model.number="transactionsDomain.txForm.amount" type="number" min="0.0001" step="0.0001" placeholder="Miktar" required />
          <input v-model="transactionsDomain.txForm.transaction_date" type="datetime-local" />
          <input class="full" v-model="transactionsDomain.txForm.description" type="text" placeholder="Açıklama" />
        </div>
        <div class="form-actions">
          <button :disabled="transactionsDomain.loading.value">{{ transactionsDomain.loading.value ? 'Kaydediliyor...' : 'Kaydet' }}</button>
        </div>
      </form>

      <section class="filters-panel">
        <h3>Filtreler & Dışa Aktar</h3>
        <div class="filters-toolbar">
          <select v-model="transactionsDomain.filters.asset">
            <option value="">Tüm Varlıklar</option>
            <option v-for="a in assetTypes" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>

          <select v-model="transactionsDomain.filters.type">
            <option value="">Tüm Tipler</option>
            <option value="add">Ekleme</option>
            <option value="subtract">Çıkarma</option>
          </select>

          <input v-model="transactionsDomain.filters.dateFrom" type="date" />
          <input v-model="transactionsDomain.filters.dateTo" type="date" />

          <input
              class="filter-search"
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
        </div>

        <div class="filter-actions">
          <button class="secondary" @click="transactionsDomain.fetchTransactions">İşlemleri Yenile</button>
          <button class="secondary" @click="transactionsDomain.downloadTxExcel">Excel İndir</button>
        </div>
      </section>
    </div>

    <p class="subtle" v-if="transactionsDomain.pagedTransactions.value.length">
      Toplam sonuç: {{ transactionsDomain.filteredTransactions.value.length }}
    </p>

    <table v-if="transactionsDomain.pagedTransactions.value.length">
      <thead>
      <tr><th>İşlem</th><th>Varlık</th><th>Miktar</th><th>Oluşturulma</th><th>İşlem Tarihi</th><th>Kur</th><th>Toplam</th><th></th></tr>
      </thead>
      <tbody>
      <tr v-for="tx in transactionsDomain.pagedTransactions.value" :key="tx.id">
        <td>
          <div class="tx-type-wrap">
            <span class="type-badge" :class="tx.type">{{ getTransactionTypeLabel(tx.type) }}</span>
            <small>#{{ tx.id }}</small>
          </div>
        </td>
        <td>{{ tx.user_asset?.asset?.toUpperCase() || '-' }}</td>
        <td>{{ formatNumber(tx.amount, 4) }}</td>
        <td>{{ formatDate(tx.created_at, { includeSeconds: false }) }}</td>
        <td>{{ formatDate(tx.transaction_date, { includeSeconds: false }) }}</td>
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

function getTransactionTypeLabel(type) {
  if (type === 'add') return 'Ekleme'
  if (type === 'subtract') return 'Çıkarma'
  return type || '-'
}
</script>

<style scoped>
h3 { margin: 0 0 .55rem; font-size: .95rem; }
.subtle { color: var(--color-muted); margin: .4rem 0; }
.panels { align-items: start; }
.transaction-form,
.filters-panel {
  border: 1px solid #d9e5f6;
  border-radius: 12px;
  padding: .78rem;
  background: linear-gradient(180deg, #ffffff, #f8fbff);
}
.form-grid {
  display: grid;
  gap: .55rem;
  grid-template-columns: repeat(2, minmax(140px, 1fr));
}
.form-grid .full { grid-column: 1 / -1; }
.form-actions {
  margin-top: .12rem;
  display: flex;
  justify-content: flex-end;
}
.filters-toolbar {
  display: grid;
  gap: .55rem;
  grid-template-columns: repeat(2, minmax(120px, 1fr));
  margin-bottom: .55rem;
}
.filter-search { grid-column: 1 / -1; }
.filter-actions {
  display: flex;
  gap: .5rem;
  justify-content: flex-end;
  flex-wrap: wrap;
}
.tx-type-wrap {
  display: grid;
  gap: .18rem;
}
.tx-type-wrap small {
  color: #94a3b8;
  font-size: .7rem;
  font-weight: 600;
}
.type-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: fit-content;
  border-radius: 999px;
  padding: .2rem .5rem;
  border: 1px solid transparent;
  font-size: .73rem;
  font-weight: 700;
}
.type-badge.add {
  background: #e9f9ef;
  border-color: #bbf7d0;
  color: #15803d;
}
.type-badge.subtract {
  background: #fff3f0;
  border-color: #fdba74;
  color: #c2410c;
}
.pager {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: .5rem;
  margin-top: .6rem;
}
@media (max-width: 760px) {
  .form-grid,
  .filters-toolbar {
    grid-template-columns: 1fr;
  }
  .form-actions,
  .filter-actions {
    justify-content: stretch;
  }
  .form-actions button,
  .filter-actions button {
    width: 100%;
  }
}
</style>