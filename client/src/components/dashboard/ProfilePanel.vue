<template>
  <AppCard title="Profil" subtitle="Hesap bilgilerinizi güncelleyin" variant="highlight">
    <template #actions>
      <button class="danger" :disabled="profileDomain.loading.value" @click="profileDomain.deleteProfile">Hesabı Sil</button>
    </template>

    <div class="profile-head">
      <strong>{{ fullName }}</strong>
      <span>{{ email }}</span>
    </div>

    <StatusBanner type="error" :message="profileDomain.status.error" />
    <StatusBanner type="ok" :message="profileDomain.status.ok" />

    <form class="inline-form" @submit.prevent="profileDomain.updateProfile">
      <input v-model="profileDomain.form.name" placeholder="Yeni ad" />
      <input v-model="profileDomain.form.surname" placeholder="Yeni soyad" />
      <input v-model="profileDomain.form.email" type="email" placeholder="Yeni e-posta" />
      <input v-model="profileDomain.form.password" type="password" placeholder="Yeni şifre" />
      <button :disabled="profileDomain.loading.value">{{ profileDomain.loading.value ? 'Kaydediliyor...' : 'Profili Güncelle' }}</button>
    </form>
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
</style>