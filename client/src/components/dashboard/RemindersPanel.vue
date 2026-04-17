<template>
  <AppCard title="Hatırlatıcılar" subtitle="Yaklaşan ve gecikmiş tarihleri takip edin">
    <StatusBanner type="error" :message="remindersDomain.status.error" />
    <StatusBanner type="ok" :message="remindersDomain.status.ok" />

    <div class="grid two-col">
      <form class="inline-form" @submit.prevent="remindersDomain.createReminder">
        <h3>Yeni Hatırlatıcı</h3>
        <input v-model="remindersDomain.form.title" type="text" placeholder="Başlık" required />
        <input v-model="remindersDomain.form.date" type="datetime-local" required />
        <button :disabled="remindersDomain.loading.value">{{ remindersDomain.loading.value ? 'Ekleniyor...' : 'Ekle' }}</button>
      </form>

      <div>
        <h3>Liste</h3>
        <ul class="list" v-if="remindersDomain.reminders.value.length">
          <li v-for="r in remindersDomain.reminders.value" :key="r.id">
            <div>
              <strong>{{ r.title }}</strong>
              <p>{{ formatDate(r.date) }}</p>
              <span class="tag" :class="getReminderMeta(r.date).tone">{{ getReminderMeta(r.date).label }}</span>
            </div>
            <button class="danger ghost" @click="remindersDomain.deleteReminder(r.id)">Sil</button>
          </li>
        </ul>
        <p v-else class="subtle">Hatırlatıcı bulunamadı.</p>
      </div>
    </div>
  </AppCard>
</template>

<script setup>
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'
import { formatDate } from '../../utils/format'

defineProps({
  remindersDomain: { type: Object, required: true },
})

function getReminderMeta(dateValue) {
  const now = Date.now()
  const target = new Date(dateValue).getTime()
  const diffDays = Math.ceil((target - now) / (1000 * 60 * 60 * 24))

  if (Number.isNaN(target)) return { label: 'Tarih belirsiz', tone: 'neutral' }
  if (diffDays < 0) return { label: 'Süresi geçti', tone: 'danger' }
  if (diffDays <= 7) return { label: `${diffDays} gün kaldı`, tone: 'warning' }
  if (diffDays <= 30) return { label: `${Math.ceil(diffDays / 7)} hafta kaldı`, tone: 'info' }
  return { label: 'Planlandı', tone: 'neutral' }
}
</script>

<style scoped>
h3 { margin: 0 0 .55rem; font-size: .95rem; }
.subtle { color: var(--color-muted); }
.tag {
  display: inline-flex;
  margin-top: .45rem;
  padding: .2rem .48rem;
  border-radius: 999px;
  font-size: .72rem;
  border: 1px solid transparent;
}
.tag.neutral { background: rgba(148, 163, 184, 0.15); border-color: rgba(148, 163, 184, 0.4); }
.tag.info { background: rgba(59, 130, 246, 0.15); border-color: rgba(59, 130, 246, 0.45); color: #bfdbfe; }
.tag.warning { background: rgba(245, 158, 11, 0.14); border-color: rgba(245, 158, 11, 0.4); color: #fde68a; }
.tag.danger { background: rgba(239, 68, 68, 0.14); border-color: rgba(239, 68, 68, 0.4); color: #fecaca; }
</style>