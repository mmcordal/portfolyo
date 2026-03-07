<template>
  <AppCard title="Portföy" subtitle="Varlık dağılımınızı takip edin">
    <template #actions>
      <button @click="$emit('download-pdf')">Portföy PDF</button>
    </template>

    <div class="toolbar">
      <label>Hedef Para</label>
      <select :value="currency" @change="$emit('currency-change', $event.target.value)">
        <option v-for="c in currencies" :key="c" :value="c">{{ c.toUpperCase() }}</option>
      </select>
    </div>

    <p class="summary" v-if="assetsAll">Toplam: {{ formatNumber(assetsAll.total_price) }} {{ assetsAll.currency?.toUpperCase() }}</p>

    <table v-if="assetsAll?.assets?.length">
      <thead><tr><th>Varlık</th><th>Miktar</th><th>Birim</th><th>Toplam</th><th></th></tr></thead>
      <tbody>
      <tr v-for="asset in assetsAll.assets" :key="asset.id">
        <td>{{ asset.asset.toUpperCase() }}</td>
        <td>{{ formatNumber(asset.amount, 4) }}</td>
        <td>{{ formatNumber(asset.price) }}</td>
        <td>{{ formatNumber(asset.total_price_by_asset) }}</td>
        <td><button @click="$emit('show-asset', asset.asset)">Detay</button></td>
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
import AppCard from '../ui/AppCard.vue'
import { formatDate, formatNumber } from '../../utils/format'

defineProps({
  assetsAll: { type: Object, default: null },
  singleAsset: { type: Object, default: null },
  currency: { type: String, required: true },
  currencies: { type: Array, required: true },
})

defineEmits(['download-pdf', 'currency-change', 'show-asset'])
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
.list li {
  display: grid;
  grid-template-columns: .8fr .7fr 1fr;
  gap: .5rem;
}
</style>