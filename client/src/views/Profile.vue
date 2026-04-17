<template>
  <main class="page">
    <div class="hero">
      <h1>Profil Ayarları</h1>
      <p>Hesap bilgilerinizi buradan güncelleyebilir veya hesabınızı silebilirsiniz.</p>
    </div>

    <ProfilePanel
        :full-name="userStore.fullName"
        :email="userStore.profile?.email"
        :profile-domain="profile"
    />
  </main>
</template>

<script setup>
import { onMounted } from 'vue'
import ProfilePanel from '../components/dashboard/ProfilePanel.vue'
import { useDashboardData } from '../composables/useDashboardData'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const { profile, bootstrap } = useDashboardData(userStore)

onMounted(async () => {
  await bootstrap({ includeDashboardData: false })
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