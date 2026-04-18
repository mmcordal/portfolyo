<template>
  <nav class="navbar" v-if="userStore.isAuthenticated">
    <button class="brand" @click="$router.push('/dashboard')">Portfolyo</button>
    <div class="actions">
      <button class="name-link" type="button" @click="$router.push('/profile')">{{ userStore.fullName }}</button>
      <button type="button" class="danger" @click="userStore.logout()">Çıkış</button>
    </div>
  </nav>
</template>

<script setup>
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
</script>

<style scoped>
.navbar {
  position: sticky;
  top: 0;
  z-index: 20;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: .9rem 1.15rem;
  backdrop-filter: blur(10px);
  background: rgba(248, 251, 255, 0.88);
  border-bottom: 1px solid #d8e4f8;
  box-shadow: 0 6px 24px rgba(15, 23, 42, 0.06);
}
.brand {
  background: #ffffff;
  border: 1px solid #bfd4f8;
  color: #1e3a8a;
  font-weight: 700;
  padding: .55rem .9rem;
}
.actions { display: flex; gap: .75rem; align-items: center; }
.name-link {
  background: #ffffff;
  border: 1px solid #d1ddf1;
  color: var(--color-muted);
  font-weight: 600;
}
</style>