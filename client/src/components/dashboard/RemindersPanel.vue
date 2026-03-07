<template>
  <AppCard title="Hatırlatıcılar" subtitle="Kişisel notlar ve önemli tarihler">
    <div class="grid two-col">
      <form class="inline-form" @submit.prevent="$emit('create')">
        <h3>Yeni Hatırlatıcı</h3>
        <input v-model="reminderForm.title" type="text" placeholder="Başlık" required />
        <input v-model="reminderForm.date" type="datetime-local" />
        <button :disabled="loadingReminder">{{ loadingReminder ? 'Ekleniyor...' : 'Ekle' }}</button>
      </form>

      <div>
        <h3>Liste</h3>
        <ul class="list" v-if="reminders.length">
          <li v-for="r in reminders" :key="r.id">
            <div>
              <strong>{{ r.title }}</strong>
              <p>{{ formatDate(r.date) }}</p>
            </div>
            <button class="danger" @click="$emit('delete', r.id)">Sil</button>
          </li>
        </ul>
        <p v-else class="subtle">Hatırlatıcı bulunamadı.</p>
      </div>
    </div>
  </AppCard>
</template>

<script setup>
import AppCard from '../ui/AppCard.vue'
import { formatDate } from '../../utils/format'

defineProps({
  reminderForm: { type: Object, required: true },
  reminders: { type: Array, required: true },
  loadingReminder: { type: Boolean, default: false },
})

defineEmits(['create', 'delete'])
</script>

<style scoped>
h3 { margin: 0 0 .55rem; font-size: .95rem; }
.subtle { color: var(--color-muted); }
</style>