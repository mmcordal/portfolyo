<template>
  <nav class="navbar" v-if="userStore.isAuthenticated">
    <div class="brand" @click="$router.push('/dashboard')">Portfolyo</div>
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.4rem;
  border-bottom: 1px solid var(--color-surface-light);
  position: sticky;
  top: 0;
  z-index: 10;
  background: var(--color-bg);
}
.brand {
  font-weight: 700;
  font-size: 1.1rem;
  cursor: pointer;
}
.actions {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}
.name {
  color: var(--color-muted);
}
</style>