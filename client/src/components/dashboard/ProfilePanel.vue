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
      <button :disabled="profileDomain.loading.value">{{ profileDomain.loading.value ? 'Kaydediliyor...' : 'Profili Güncelle' }}</button>
    </form>

    <form class="inline-form section-card" @submit.prevent="profileDomain.updateProfile">
      <h3>Güvenlik</h3>
      <p>Hesap güvenliğinizi artırmak için şifrenizi düzenli aralıklarla güncelleyin.</p>
      <input v-model="profileDomain.form.password" type="password" placeholder="Yeni şifre" />
      <button class="secondary" :disabled="profileDomain.loading.value">{{ profileDomain.loading.value ? 'Kaydediliyor...' : 'Şifreyi Güncelle' }}</button>
    </form>

    <section class="danger-zone">
      <h3>Danger Zone</h3>
      <p>Bu işlem geri alınamaz. Hesabınızı silerseniz tüm profil erişiminiz sonlandırılır.</p>
      <button class="danger" :disabled="profileDomain.loading.value" @click="profileDomain.deleteProfile">Hesabı Sil</button>
    </section>
  </AppCard>
</template>

<script setup>
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'

defineProps({
  fullName: { type: String, default: '-' },
  email: { type: String, default: '' },
  profileDomain: { type: Object, required: true },
})
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
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(120px, 1fr));
  gap: .52rem;
}
.danger-zone {
  border: 1px solid #fecaca;
  border-radius: 12px;
  background: #fff6f6;
  padding: .8rem;
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
@media (max-width: 760px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>