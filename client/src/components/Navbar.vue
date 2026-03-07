<template>
  <nav class="navbar" v-if="userStore.isAuthenticated">
    <button class="brand" @click="$router.push('/dashboard')">Portfolyo</button>
    <div class="actions">
      <span class="name">{{ userStore.fullName }}</span>
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
  background: rgba(2, 6, 23, 0.72);
  border-bottom: 1px solid rgba(37, 50, 82, .8);
}
.brand {
  background: transparent;
  border: 1px solid rgba(59,130,246,.4);
  color: #dbeafe;
  font-weight: 700;
  padding: .55rem .9rem;
}
.actions { display: flex; gap: .75rem; align-items: center; }
.name { color: var(--color-muted); font-weight: 600; }
</style>