<template>
  <section class="app-card" :class="variantClass">
    <header v-if="title || subtitle" class="app-card__head">
      <div>
        <h2>{{ title }}</h2>
        <p v-if="subtitle">{{ subtitle }}</p>
      </div>
      <slot name="actions" />
    </header>
    <slot />
  </section>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
  variant: { type: String, default: 'default' },
})

const variantClass = computed(() => `app-card--${props.variant}`)
</script>

<style scoped>
.app-card {
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--color-surface-light);
  border-radius: 16px;
  padding: 1rem;
  box-shadow: 0 14px 30px rgba(2, 6, 23, 0.3);
}
.app-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: .9rem;
}
.app-card__head h2 { margin: 0; font-size: 1.08rem; }
.app-card__head p { margin: .25rem 0 0; color: var(--color-muted); }
.app-card--highlight {
  border-color: rgba(59, 130, 246, 0.6);
  background: linear-gradient(160deg, rgba(59, 130, 246, 0.22), rgba(30, 41, 59, 0.9));
}
</style>