<template>
  <AppCard title="Profil Yönetimi" subtitle="Bilgilerinizi ve güvenlik tercihlerinizi güncelleyin" variant="highlight">
    <div class="profile-head">
      <strong>{{ fullName }}</strong>
      <span>{{ email }}</span>
    </div>

    <StatusBanner type="error" :message="profileDomain.status.error" />
    <StatusBanner type="ok" :message="profileDomain.status.ok" />

    <form class="inline-form section-card" @submit.prevent="profileDomain.updateProfile">
      <h3>Profil Bilgileri</h3>
      <div class="form-grid">
        <input v-model="profileDomain.form.name" placeholder="Yeni ad" />
        <input v-model="profileDomain.form.surname" placeholder="Yeni soyad" />
      </div>
      <input v-model="profileDomain.form.email" type="email" placeholder="Yeni e-posta" />
      <button type="submit" :disabled="profileDomain.loading.value">{{ profileDomain.loading.value ? 'Kaydediliyor...' : 'Profili Güncelle' }}</button>
    </form>

    <form class="inline-form section-card security-card" @submit.prevent="profileDomain.updateProfile">
      <h3>Güvenlik</h3>
      <p>Hesap güvenliğinizi artırmak için şifrenizi düzenli aralıklarla güncelleyin.</p>
      <input v-model="profileDomain.form.password" type="password" placeholder="Yeni şifre" />
      <button type="submit" class="secondary" :disabled="profileDomain.loading.value">{{ profileDomain.loading.value ? 'Kaydediliyor...' : 'Şifre Güncelle' }}</button>
    </form>

    <section class="danger-zone">
      <h3>Danger Zone</h3>
      <p>Bu işlem geri alınamaz. Hesabınızı silerseniz profiliniz, işlem geçmişiniz ve erişim bilgileriniz kalıcı olarak silinir.</p>
      <button class="danger" :disabled="profileDomain.loading.value" @click="openDeleteModal">Hesabı Sil</button>
    </section>

    <div v-if="deleteConfirmOpen" class="overlay" @click.self="closeDeleteModal">
      <section class="modal-card delete-modal">
        <div class="modal-head">
          <h3>Hesabı Kalıcı Olarak Sil</h3>
          <button class="secondary" @click="closeDeleteModal">Kapat</button>
        </div>

        <p>Bu işlem geri alınamaz. Onay için aşağıya <strong>SİL</strong> yazın.</p>
        <input v-model="deleteConfirmation" type="text" placeholder="SİL" />

        <div class="modal-actions">
          <button class="secondary" @click="closeDeleteModal">İptal</button>
          <button
              class="danger"
              :disabled="profileDomain.loading.value || deleteConfirmation.trim().toUpperCase() !== 'SİL'"
              @click="confirmDeleteProfile"
          >
            {{ profileDomain.loading.value ? 'Siliniyor...' : 'Hesabı Kalıcı Olarak Sil' }}
          </button>
        </div>
      </section>
    </div>
  </AppCard>
</template>

<script setup>
import { ref } from 'vue'
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'

const props = defineProps({
  fullName: { type: String, default: '-' },
  email: { type: String, default: '' },
  profileDomain: { type: Object, required: true },
})

const deleteConfirmOpen = ref(false)
const deleteConfirmation = ref('')

function openDeleteModal() {
  deleteConfirmOpen.value = true
  deleteConfirmation.value = ''
}

function closeDeleteModal() {
  if (deleteConfirmOpen.value) {
    deleteConfirmOpen.value = false
    deleteConfirmation.value = ''
  }
}

async function confirmDeleteProfile() {
  if (deleteConfirmation.value.trim().toUpperCase() !== 'SİL') return

  await props.profileDomain.deleteProfile()
  if (!props.profileDomain.status.error) {
    closeDeleteModal()
  }
}
</script>

<style scoped>
.profile-head {
  display: grid;
  gap: .2rem;
  margin-bottom: .8rem;
}
.profile-head strong { font-size: 1.1rem; }
.profile-head span { color: var(--color-muted); }
.section-card {
  border: 1px solid #dbe6f8;
  border-radius: 12px;
  background: #ffffff;
  padding: .8rem;
  margin-bottom: .7rem;
}
.section-card h3 {
  font-size: .95rem;
  margin-bottom: .5rem;
}
.section-card p {
  color: #64748b;
  margin-bottom: .5rem;
  font-size: .9rem;
}
.security-card {
  background: linear-gradient(180deg, #ffffff, #f8fbff);
}
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(120px, 1fr));
  gap: .52rem;
}
.danger-zone {
  border: 1px solid #fecaca;
  border-radius: 12px;
  background: linear-gradient(180deg, #fff7f7, #fff2f2);
  padding: .9rem;
}
.danger-zone h3 {
  color: #b91c1c;
  margin-bottom: .35rem;
  font-size: .95rem;
}
.danger-zone p {
  color: #7f1d1d;
  margin-bottom: .65rem;
}
.delete-modal p {
  color: #334155;
  line-height: 1.5;
  margin-bottom: .55rem;
}
.modal-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: .7rem;
}
.modal-actions {
  margin-top: .8rem;
  display: flex;
  justify-content: flex-end;
  gap: .5rem;
}
@media (max-width: 760px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  .modal-actions {
    flex-direction: column-reverse;
  }
  .modal-actions button {
    width: 100%;
  }
}
</style>