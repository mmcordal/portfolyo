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
  background: #eefbf2;
  border-color: rgba(22, 163, 74, 0.3);
  color: #166534;
}
.status--error {
  background: #fff2f2;
  border-color: rgba(220, 38, 38, 0.3);
  color: #b91c1c;
}
.fade-enter-active, .fade-leave-active { transition: opacity .24s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>