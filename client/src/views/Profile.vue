<template>
  <main class="page profile-page">
    <div class="hero">
      <h1>Profil Ayarları</h1>
      <p>Hesap bilgilerinizi güvenli şekilde yönetin, güncelleyin ve kritik işlemleri kontrollü biçimde yapın.</p>
    </div>

    <section class="profile-layout">
      <AppCard title="Kullanıcı Özeti" subtitle="Hesap durumu ve kimlik bilgileri">
        <div class="user-summary">
          <div class="avatar">{{ initials }}</div>
          <div>
            <strong>{{ userStore.fullName }}</strong>
            <p>{{ userStore.profile?.email }}</p>
          </div>
        </div>
        <ul class="summary-list">
          <li>
            <span>E-posta</span>
            <strong>{{ userStore.profile?.email || '-' }}</strong>
          </li>
          <li>
            <span>Hesap Durumu</span>
            <strong>Aktif</strong>
          </li>
          <li>
            <span>Panel Rolü</span>
            <strong>Standart Kullanıcı</strong>
          </li>
        </ul>
      </AppCard>

      <ProfilePanel
          :full-name="userStore.fullName"
          :email="userStore.profile?.email"
          :profile-domain="profile"
      />
    </section>
  </main>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import AppCard from '../components/ui/AppCard.vue'
import ProfilePanel from '../components/dashboard/ProfilePanel.vue'
import { useDashboardData } from '../composables/useDashboardData'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const { profile, bootstrap } = useDashboardData(userStore)

const initials = computed(() => {
  const fallback = 'U'
  const fullName = userStore.fullName || ''
  const tokens = fullName.split(' ').filter(Boolean)
  if (!tokens.length) return fallback
  return tokens.slice(0, 2).map((part) => part[0]?.toUpperCase() || '').join('')
})

onMounted(async () => {
  await bootstrap({ includeDashboardData: false })
})
</script>

<style scoped>
.profile-page {
  gap: .9rem;
}
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
.profile-layout {
  display: grid;
  grid-template-columns: minmax(220px, .85fr) minmax(0, 1.7fr);
  gap: .95rem;
}
.user-summary {
  display: flex;
  gap: .65rem;
  align-items: center;
}
.avatar {
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: linear-gradient(160deg, #1d4ed8, #60a5fa);
  color: #ffffff;
  display: grid;
  place-items: center;
  font-weight: 700;
}
.user-summary p {
  color: #64748b;
  margin-top: .2rem;
}
.summary-list {
  margin: .8rem 0 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: .55rem;
}
.summary-list li {
  display: flex;
  justify-content: space-between;
  gap: .5rem;
  border: 1px solid #dbe6f7;
  border-radius: 10px;
  padding: .6rem;
  background: #fbfdff;
}
.summary-list span {
  color: #64748b;
  font-size: .82rem;
}
.summary-list strong {
  font-size: .86rem;
  color: #1e3a8a;
}
@media (max-width: 900px) {
  .profile-layout {
    grid-template-columns: 1fr;
  }
}
</style>