<template>
  <transition name="fade">
    <div v-if="message" class="status" :class="typeClass">
      {{ message }}
    </div>
  </transition>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  message: { type: String, default: '' },
  type: { type: String, default: 'ok' },
})

const typeClass = computed(() => `status--${props.type}`)
</script>

<style scoped>
.status {
  margin: 1rem 0;
  border-radius: 12px;
  padding: .8rem 1rem;
  font-weight: 600;
  border: 1px solid transparent;
}
.status--ok {
  background: rgba(34, 197, 94, 0.12);
  border-color: rgba(34, 197, 94, 0.45);
  color: #bbf7d0;
}
.status--error {
  background: rgba(239, 68, 68, 0.12);
  border-color: rgba(239, 68, 68, 0.45);
  color: #fecaca;
}
.fade-enter-active, .fade-leave-active { transition: opacity .24s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>