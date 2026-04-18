<template>
  <AppCard title="İşlemler" subtitle="Filtreleme, raporlama ve işlem geçmişi">
    <template #actions>
      <button @click="showCreateModal = true">+ Yeni İşlem</button>
    </template>

    <StatusBanner type="error" :message="transactionsDomain.status.error" />
    <StatusBanner type="ok" :message="transactionsDomain.status.ok" />

    <section class="filters-panel">
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

        <button class="secondary" @click="transactionsDomain.fetchTransactions">Yenile</button>
        <button class="secondary" @click="transactionsDomain.downloadTxExcel">Excel</button>
      </div>
    </section>

    <section v-if="transactionsDomain.listLoading.value" class="list-loading" aria-live="polite">
      <strong>İşlemler yükleniyor...</strong>
      <p>Filtrelenen sonuçlar hazırlanıyor.</p>
    </section>

    <template v-else>
      <p class="subtle" v-if="transactionsDomain.pagedTransactions.value.length">
        Toplam sonuç: {{ transactionsDomain.filteredTransactions.value.length }}
      </p>

      <div class="table-shell" v-if="transactionsDomain.pagedTransactions.value.length">
        <table class="transactions-table">
          <thead>
          <tr>
            <th>İşlem</th>
            <th>Varlık</th>
            <th>Miktar</th>
            <th>Oluşturulma</th>
            <th>İşlem Tarihi</th>
            <th>Kur</th>
            <th class="align-right">Toplam</th>
            <th class="align-right">Belge</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="tx in transactionsDomain.pagedTransactions.value" :key="tx.id">
            <td>
              <div class="tx-type-wrap">
                <span class="type-badge" :class="tx.type">{{ getTransactionTypeLabel(tx.type) }}</span>
                <small>#{{ tx.id }}</small>
              </div>
            </td>
            <td><span class="asset-badge">{{ tx.user_asset?.asset?.toUpperCase() || '-' }}</span></td>
            <td>{{ formatNumber(tx.amount, 4) }}</td>
            <td>{{ formatDate(tx.created_at, { includeSeconds: false }) }}</td>
            <td>{{ formatDate(tx.transaction_date, { includeSeconds: false }) }}</td>
            <td>{{ formatNumber(tx.target_currency_price) }}</td>
            <td class="total-cell align-right">{{ formatNumber(tx.target_currency_total_price) }} {{ tx.now_target_currency?.toUpperCase() }}</td>
            <td class="align-right">
              <button class="secondary pdf-action" @click="transactionsDomain.downloadTxPdf(tx.id)">PDF</button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>

      <div class="tx-mobile-cards" v-if="transactionsDomain.pagedTransactions.value.length">
        <article class="tx-card" v-for="tx in transactionsDomain.pagedTransactions.value" :key="`mobile-${tx.id}`">
          <div class="tx-card-top">
            <span class="type-badge" :class="tx.type">{{ getTransactionTypeLabel(tx.type) }}</span>
            <span class="asset-badge">{{ tx.user_asset?.asset?.toUpperCase() || '-' }}</span>
          </div>
          <strong>{{ formatNumber(tx.target_currency_total_price) }} {{ tx.now_target_currency?.toUpperCase() }}</strong>
          <p>Miktar: {{ formatNumber(tx.amount, 4) }}</p>
          <p>İşlem: {{ formatDate(tx.transaction_date, { includeSeconds: false }) }}</p>
          <p>Oluşturulma: {{ formatDate(tx.created_at, { includeSeconds: false }) }}</p>
          <button class="secondary" @click="transactionsDomain.downloadTxPdf(tx.id)">PDF İndir</button>
        </article>
      </div>

      <div v-else class="empty-state">
        <h3>Filtreye uygun işlem bulunamadı</h3>
        <p>Tarih, varlık veya arama filtresini sadeleştirerek tekrar deneyin.</p>
        <button class="secondary" @click="clearFilters">Filtreleri Temizle</button>
      </div>
    </template>

    <div class="pager" v-if="transactionsDomain.totalPages.value > 1">
      <button :disabled="transactionsDomain.filters.page <= 1" @click="transactionsDomain.setPage(transactionsDomain.filters.page - 1)">Önceki</button>
      <span>Sayfa {{ transactionsDomain.filters.page }} / {{ transactionsDomain.totalPages.value }}</span>
      <button :disabled="transactionsDomain.filters.page >= transactionsDomain.totalPages.value" @click="transactionsDomain.setPage(transactionsDomain.filters.page + 1)">Sonraki</button>
    </div>

    <div v-if="showCreateModal" class="overlay" @click.self="closeCreateModal">
      <section class="modal-card">
        <div class="modal-head">
          <h3>Yeni İşlem Ekle</h3>
          <button class="secondary" @click="closeCreateModal">Kapat</button>
        </div>

        <form class="inline-form transaction-form" @submit.prevent="submitCreate">
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
      </section>
    </div>
  </AppCard>
</template>

<script setup>
import { ref, watch } from 'vue'
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'
import { formatDate, formatNumber } from '../../utils/format'

const props = defineProps({
  transactionsDomain: { type: Object, required: true },
  actionTypes: { type: Array, required: true },
  assetTypes: { type: Array, required: true },
  openCreate: { type: Boolean, default: false },
})

const emit = defineEmits(['create', 'close-create'])
const showCreateModal = ref(false)

watch(() => props.openCreate, (next) => {
  showCreateModal.value = next
}, { immediate: true })

function closeCreateModal() {
  showCreateModal.value = false
  emit('close-create')
}

function submitCreate() {
  emit('create', closeCreateModal)
}

function getTransactionTypeLabel(type) {
  if (type === 'add') return 'Ekleme'
  if (type === 'subtract') return 'Çıkarma'
  return type || '-'
}

function clearFilters() {
  props.transactionsDomain.filters.asset = ''
  props.transactionsDomain.filters.type = ''
  props.transactionsDomain.filters.search = ''
  props.transactionsDomain.filters.dateFrom = ''
  props.transactionsDomain.filters.dateTo = ''
  props.transactionsDomain.setPage(1)
  props.transactionsDomain.fetchTransactions()
}
</script>

<style scoped>
h3 { margin: 0 0 .55rem; font-size: .95rem; }
.subtle { color: var(--color-muted); margin: .4rem 0; }
.filters-panel {
  border: 1px solid #d9e5f6;
  border-radius: 12px;
  padding: .78rem;
  background: linear-gradient(180deg, #ffffff, #f8fbff);
}
.filters-toolbar {
  display: grid;
  grid-template-columns: repeat(8, minmax(120px, 1fr));
  gap: .62rem;
  align-items: center;
}
.filter-search {
  grid-column: span 2;
}
.list-loading {
  border: 1px solid #d9e5f6;
  border-radius: 12px;
  padding: .78rem;
  background: #f8fbff;
}
.list-loading p {
  margin-top: .2rem;
  color: var(--color-muted);
}
.table-shell { overflow-x: auto; }
.transactions-table { min-width: 940px; }
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
.asset-badge {
  display: inline-flex;
  padding: .2rem .45rem;
  border-radius: 7px;
  background: #ecf3ff;
  border: 1px solid #cddcf8;
  color: #1e40af;
  font-weight: 700;
  font-size: .74rem;
}
.align-right {
  text-align: right;
}
.total-cell {
  color: #0f2f73;
  font-weight: 800;
  letter-spacing: .01em;
}
.pdf-action {
  padding: .38rem .58rem;
  font-size: .74rem;
}
.empty-state {
  border: 1px dashed #c9d9f4;
  border-radius: 12px;
  padding: .95rem;
  background: #f9fbff;
  display: grid;
  gap: .42rem;
}
.empty-state h3 {
  margin: 0;
  font-size: .98rem;
}
.empty-state p {
  color: #64748b;
}
.empty-state button {
  width: fit-content;
}
.pager {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: .5rem;
  margin-top: .6rem;
}
.tx-mobile-cards {
  display: none;
  gap: .6rem;
}
.tx-card {
  border: 1px solid #d9e5f6;
  border-radius: 12px;
  background: #ffffff;
  padding: .72rem;
  display: grid;
  gap: .3rem;
}
.tx-card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: .45rem;
}
.modal-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: .6rem;
  margin-bottom: .7rem;
}
.transaction-form {
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
@media (max-width: 1120px) {
  .filters-toolbar {
    grid-template-columns: repeat(4, minmax(130px, 1fr));
  }
  .filter-search {
    grid-column: span 2;
  }
}
@media (max-width: 760px) {
  .filters-toolbar,
  .form-grid {
    grid-template-columns: 1fr;
  }
  .filter-search {
    grid-column: auto;
  }
  .table-shell { display: none; }
  .tx-mobile-cards {
    display: grid;
  }
  .form-actions {
    justify-content: stretch;
  }
  .form-actions button {
    width: 100%;
  }
}
</style>