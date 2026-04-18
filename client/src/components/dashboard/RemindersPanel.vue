<template>
  <AppCard title="Hatırlatıcılar" subtitle="Yaklaşan planları düzenli ve sade bir görünümde takip edin">
    <template #actions>
      <button class="secondary" @click="showCreateModal = true">+ Hatırlatıcı</button>
    </template>

    <StatusBanner type="error" :message="remindersDomain.status.error" />
    <StatusBanner type="ok" :message="remindersDomain.status.ok" />

    <ul class="list" v-if="remindersDomain.reminders.value.length">
      <li v-for="r in remindersDomain.reminders.value" :key="r.id" class="reminder-item">
        <div class="reminder-content">
          <strong>{{ r.title }}</strong>
          <p>{{ formatDate(r.date, { includeSeconds: false }) }}</p>
        </div>
        <div class="reminder-actions">
          <span class="tag" :class="getReminderMeta(r.date).tone">{{ getReminderMeta(r.date).label }}</span>
          <button class="danger ghost" @click="remindersDomain.deleteReminder(r.id)">Sil</button>
        </div>
      </li>
    </ul>
    <p v-else class="subtle">Hatırlatıcı bulunamadı.</p>

    <div v-if="showCreateModal" class="overlay" @click.self="closeCreateModal">
      <section class="modal-card reminder-modal">
        <div class="modal-head">
          <h3>Yeni Hatırlatıcı</h3>
          <button class="secondary" @click="closeCreateModal">Kapat</button>
        </div>

        <form class="inline-form compact-form" @submit.prevent="submitReminder">
          <input v-model="remindersDomain.form.title" type="text" placeholder="Başlık" required />
          <input v-model="remindersDomain.form.date" type="datetime-local" required />
          <button :disabled="remindersDomain.loading.value">{{ remindersDomain.loading.value ? 'Ekleniyor...' : 'Ekle' }}</button>
        </form>
      </section>
    </div>
  </AppCard>
</template>

<script setup>
import { ref, watch } from 'vue'
import AppCard from '../ui/AppCard.vue'
import StatusBanner from '../ui/StatusBanner.vue'
import { formatDate } from '../../utils/format'

const props = defineProps({
  remindersDomain: { type: Object, required: true },
  openCreate: { type: Boolean, default: false },
})

const emit = defineEmits(['close-create'])
const showCreateModal = ref(false)

watch(() => props.openCreate, (next) => {
  showCreateModal.value = next
}, { immediate: true })

function closeCreateModal() {
  showCreateModal.value = false
  emit('close-create')
}

async function submitReminder() {
  await props.remindersDomain.createReminder()
  if (!props.remindersDomain.status.error) {
    closeCreateModal()
  }
}

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
.reminder-item {
  background: linear-gradient(180deg, #ffffff, #f8fbff);
}
.reminder-content p {
  margin-top: .2rem;
  color: #475569;
  font-size: .88rem;
}
.reminder-actions {
  display: inline-flex;
  align-items: center;
  gap: .5rem;
}
.tag {
  display: inline-flex;
  padding: .2rem .48rem;
  border-radius: 999px;
  font-size: .72rem;
  border: 1px solid transparent;
  font-weight: 600;
}
.tag.neutral { background: #eff4fc; border-color: #c8d6ed; color: #475569; }
.tag.info { background: #e8f1ff; border-color: #bcd4ff; color: #1d4ed8; }
.tag.warning { background: #fff7e8; border-color: #f9d9a0; color: #a16207; }
.tag.danger { background: #fff1f1; border-color: #fecaca; color: #b91c1c; }
.compact-form {
  border: 1px solid #dbe6f7;
  border-radius: 12px;
  background: #ffffff;
  padding: .72rem;
}
.modal-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: .7rem;
}
@media (max-width: 760px) {
  .reminder-actions {
    width: 100%;
    justify-content: space-between;
  }
}
</style>