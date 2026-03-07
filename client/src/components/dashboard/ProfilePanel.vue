<template>
  <AppCard title="Profil" subtitle="Hesap bilgilerinizi güncelleyin" variant="highlight">
    <template #actions>
      <button class="danger" :disabled="loading" @click="$emit('delete')">Hesabı Sil</button>
    </template>

    <div class="profile-head">
      <strong>{{ fullName }}</strong>
      <span>{{ email }}</span>
    </div>

    <form class="inline-form" @submit.prevent="$emit('update')">
      <input v-model="form.name" placeholder="Yeni ad" />
      <input v-model="form.surname" placeholder="Yeni soyad" />
      <input v-model="form.email" type="email" placeholder="Yeni e-posta" />
      <input v-model="form.password" type="password" placeholder="Yeni şifre" />
      <button :disabled="loading">{{ loading ? 'Kaydediliyor...' : 'Profili Güncelle' }}</button>
    </form>
  </AppCard>
</template>

<script setup>
import AppCard from '../ui/AppCard.vue'

defineProps({
  fullName: { type: String, default: '-' },
  email: { type: String, default: '' },
  form: { type: Object, required: true },
  loading: { type: Boolean, default: false },
})

defineEmits(['update', 'delete'])
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