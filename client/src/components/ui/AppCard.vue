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
  background: linear-gradient(180deg, #ffffff, var(--color-surface-alt));
  border: 1px solid var(--color-border);
  border-radius: 16px;
  padding: 1.05rem;
  box-shadow: var(--shadow-card);
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
  border-color: rgba(37, 99, 235, 0.42);
  background: linear-gradient(160deg, rgba(219, 234, 254, 0.95), #ffffff);
}

@media (max-width: 760px) {
  .app-card {
    padding: .85rem;
  }
  .app-card__head {
    flex-direction: column;
  }
}
</style>