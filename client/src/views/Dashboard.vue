<template>
  <main class="page">
    <section class="card grid profile-grid">
      <div>
        <h2>Profil</h2>
        <p>{{ userStore.fullName }}</p>
        <p>{{ userStore.profile?.email }}</p>
      </div>
      <form @submit.prevent="updateProfile" class="inline-form">
        <input v-model="profileForm.name" placeholder="Yeni ad" />
        <input v-model="profileForm.surname" placeholder="Yeni soyad" />
        <input v-model="profileForm.email" type="email" placeholder="Yeni e-posta" />
        <input v-model="profileForm.password" type="password" placeholder="Yeni şifre" />
        <button :disabled="loading.profile">Profili Güncelle</button>
      </form>
      <button class="danger" :disabled="loading.profile" @click="deleteProfile">Hesabı Sil</button>
    </section>

    <section class="card grid two-col">
      <div>
        <h2>Portföy</h2>
        <div class="toolbar">
          <label>Hedef Para</label>
          <select v-model="currency" @change="fetchAssets">
            <option v-for="c in CURRENCIES" :key="c" :value="c">{{ c.toUpperCase() }}</option>
          </select>
          <button @click="downloadAssetsPdf">Portföy PDF</button>
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
            <td><button @click="fetchSingleAsset(asset.asset)">Detay</button></td>
          </tr>
          </tbody>
        </table>
      </div>

      <div>
        <h3>Varlık Detayı</h3>
        <p v-if="singleAsset">{{ singleAsset.asset.toUpperCase() }} - {{ formatNumber(singleAsset.total_price_by_asset) }} {{ singleAsset.target_currency?.toUpperCase() }}</p>
        <ul class="list" v-if="singleAsset?.transactions?.length">
          <li v-for="tx in singleAsset.transactions" :key="tx.id">{{ tx.type }} • {{ formatNumber(tx.amount, 4) }} • {{ formatDate(tx.transaction_date) }}</li>
        </ul>
      </div>
    </section>

    <section class="card grid two-col">
      <div>
        <h2>İşlem Ekle</h2>
        <form @submit.prevent="createTransaction" class="inline-form">
          <select v-model="txForm.type" required>
            <option v-for="a in ACTION_TYPES" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>
          <select v-model="txForm.asset" required>
            <option value="">Varlık seçin</option>
            <option v-for="a in ASSET_TYPES" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>
          <input v-model.number="txForm.amount" type="number" min="0.0001" step="0.0001" placeholder="Miktar" required />
          <input v-model="txForm.transaction_date" type="datetime-local" />
          <input v-model="txForm.description" type="text" placeholder="Açıklama" />
          <button :disabled="loading.tx">Kaydet</button>
        </form>
      </div>
      <div>
        <h3>İşlem Filtreleri</h3>
        <div class="toolbar">
          <select v-model="txAssetFilter">
            <option value="">Tüm Varlıklar</option>
            <option v-for="a in ASSET_TYPES" :key="a.value" :value="a.value">{{ a.label }}</option>
          </select>
          <button @click="fetchTransactions">İşlemleri Yenile</button>
          <button @click="downloadTxExcel">Excel İndir</button>
        </div>
      </div>
    </section>

    <section class="card">
      <h2>İşlemler</h2>
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
          <td><button @click="downloadTxPdf(tx.id)">PDF</button></td>
        </tr>
        </tbody>
      </table>
      <p v-else>İşlem bulunamadı.</p>
    </section>

    <section class="card grid two-col">
      <div>
        <h2>Hatırlatıcı Ekle</h2>
        <form @submit.prevent="createReminder" class="inline-form">
          <input v-model="reminderForm.title" type="text" placeholder="Başlık" required />
          <input v-model="reminderForm.date" type="datetime-local" />
          <button :disabled="loading.reminder">Ekle</button>
        </form>
      </div>
      <div>
        <h3>Hatırlatıcılar</h3>
        <ul class="list" v-if="reminders.length">
          <li v-for="r in reminders" :key="r.id">
            <div>
              <strong>{{ r.title }}</strong>
              <p>{{ formatDate(r.date) }}</p>
            </div>
            <button class="danger" @click="deleteReminder(r.id)">Sil</button>
          </li>
        </ul>
      </div>
    </section>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="ok" class="ok">{{ ok }}</p>
  </main>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { extractErrorMessage } from '../api'
import { ACTION_TYPES, ASSET_TYPES, CURRENCIES } from '../constants/assets'
import { authService, assetService, downloadBlob, reminderService, transactionService } from '../services/portfolio'
import { useUserStore } from '../stores/user'
import { formatDate, formatNumber, toISODateTimeLocal } from '../utils/format'

const userStore = useUserStore()

const currency = ref('try')
const txAssetFilter = ref('')
const assetsAll = ref(null)
const singleAsset = ref(null)
const transactions = ref([])
const reminders = ref([])

const txForm = reactive({ type: 'add', asset: '', amount: null, transaction_date: '', description: '' })
const reminderForm = reactive({ title: '', date: '' })
const profileForm = reactive({ name: '', surname: '', email: '', password: '' })

const loading = reactive({ profile: false, tx: false, reminder: false })
const error = ref('')
const ok = ref('')

function resetMessage() {
  error.value = ''
  ok.value = ''
}

async function bootstrap() {
  try {
    const meRes = await authService.me()
    userStore.setProfile(meRes.data.data)
    await Promise.all([fetchAssets(), fetchTransactions(), fetchReminders()])
  } catch (err) {
    error.value = extractErrorMessage(err)
    if (err?.response?.status === 401) {
      userStore.logout()
    }
  }
}

async function fetchAssets() {
  try {
    resetMessage()
    const res = await assetService.getAll(currency.value)
    assetsAll.value = res.data.data
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function fetchSingleAsset(asset) {
  try {
    resetMessage()
    const res = await assetService.getByAsset(asset, currency.value)
    singleAsset.value = res.data.data
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function downloadAssetsPdf() {
  try {
    const res = await assetService.downloadPDF(currency.value)
    downloadBlob(res.data, `portfolio-${currency.value}.pdf`)
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function createTransaction() {
  try {
    loading.tx = true
    resetMessage()
    await transactionService.create({
      ...txForm,
      transaction_date: txForm.transaction_date ? new Date(txForm.transaction_date).toISOString() : '',
    })
    txForm.asset = ''
    txForm.amount = null
    txForm.transaction_date = ''
    txForm.description = ''
    ok.value = 'İşlem başarıyla eklendi.'
    await Promise.all([fetchAssets(), fetchTransactions()])
  } catch (err) {
    error.value = extractErrorMessage(err)
  } finally {
    loading.tx = false
  }
}

async function fetchTransactions() {
  try {
    resetMessage()
    const res = txAssetFilter.value
        ? await transactionService.getByAsset(txAssetFilter.value, currency.value)
        : await transactionService.getAll(currency.value)
    transactions.value = res.data.data || []
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function downloadTxExcel() {
  try {
    const res = await transactionService.downloadExcel(currency.value)
    downloadBlob(res.data, `transactions-${currency.value}.xlsx`)
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function downloadTxPdf(txId) {
  try {
    const res = await transactionService.downloadPDF(txId, currency.value)
    downloadBlob(res.data, `transaction-${txId}.pdf`)
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function fetchReminders() {
  try {
    const res = await reminderService.getAll()
    reminders.value = res.data.data || []
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function createReminder() {
  try {
    loading.reminder = true
    resetMessage()
    await reminderService.create({
      title: reminderForm.title,
      date: reminderForm.date ? new Date(reminderForm.date).toISOString() : '',
    })
    reminderForm.title = ''
    reminderForm.date = ''
    ok.value = 'Hatırlatıcı eklendi.'
    await fetchReminders()
  } catch (err) {
    error.value = extractErrorMessage(err)
  } finally {
    loading.reminder = false
  }
}

async function deleteReminder(id) {
  try {
    resetMessage()
    await reminderService.remove(id)
    ok.value = 'Hatırlatıcı silindi.'
    await fetchReminders()
  } catch (err) {
    error.value = extractErrorMessage(err)
  }
}

async function updateProfile() {
  try {
    loading.profile = true
    resetMessage()
    const payload = Object.fromEntries(Object.entries(profileForm).filter(([, v]) => String(v || '').trim() !== ''))
    if (!Object.keys(payload).length) {
      ok.value = 'Güncellenecek alan bulunamadı.'
      return
    }
    await authService.updateMe(payload)
    Object.keys(profileForm).forEach((key) => { profileForm[key] = '' })
    ok.value = 'Profil güncellendi.'
    const meRes = await authService.me()
    userStore.setProfile(meRes.data.data)
  } catch (err) {
    error.value = extractErrorMessage(err)
  } finally {
    loading.profile = false
  }
}

async function deleteProfile() {
  const confirmed = window.confirm('Hesabınızı silmek istediğinize emin misiniz?')
  if (!confirmed) return

  try {
    loading.profile = true
    resetMessage()
    await authService.deleteMe()
    userStore.logout()
  } catch (err) {
    error.value = extractErrorMessage(err)
  } finally {
    loading.profile = false
  }
}

onMounted(async () => {
  await bootstrap()
  profileForm.email = userStore.profile?.email || ''
  profileForm.name = ''
  profileForm.surname = ''
  profileForm.password = ''
  reminderForm.date = toISODateTimeLocal(new Date())
})
</script>