<template>
  <main class="page">
    <div class="hero">
      <h1>Profil Ayarları</h1>
      <p>Hesap bilgilerinizi buradan güncelleyebilir veya hesabınızı silebilirsiniz.</p>
    </div>

    <ProfilePanel
        :full-name="userStore.fullName"
        :email="userStore.profile?.email"
        :form="profileForm"
        :loading="loading.profile"
        :status="status.profile"
        @update="updateProfile"
        @delete="deleteProfile"
    />
  </main>
</template>

<script setup>
import { onMounted } from 'vue'
import ProfilePanel from '../components/dashboard/ProfilePanel.vue'
import { useDashboardData } from '../composables/useDashboardData'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const { profileForm, loading, status, bootstrap, updateProfile, deleteProfile } = useDashboardData(userStore)

onMounted(async () => {
  await bootstrap()
})
</script>

<style scoped>
.hero {
  margin: .25rem 0 .2rem;
}
.hero h1 {
  margin: 0;
  font-size: 1.4rem;
}
.hero p {
  margin: .2rem 0 0;
  color: var(--color-muted);
}
</style>