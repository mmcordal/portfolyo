import { reactive, ref } from 'vue'
import { reminderService } from '../services/portfolio'
import { createStatusState, clearStatus, setError, setOk } from './dashboardShared'

export function useDashboardReminders() {
    const reminders = ref([])
    const loading = ref(false)
    const status = createStatusState()
    const form = reactive({ title: '', date: '' })

    async function fetchReminders() {
        try {
            loading.value = true
            const res = await reminderService.getAll()
            reminders.value = res.data || []
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    async function createReminder() {
        try {
            loading.value = true
            clearStatus(status)
            await reminderService.create({
                title: form.title,
                date: form.date ? new Date(form.date).toISOString() : '',
            })
            form.title = ''
            form.date = ''
            setOk(status, 'Hatırlatıcı eklendi.')
            await fetchReminders()
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    async function deleteReminder(id) {
        try {
            loading.value = true
            clearStatus(status)
            await reminderService.remove(id)
            setOk(status, 'Hatırlatıcı silindi.')
            await fetchReminders()
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    return {
        reminders,
        form,
        loading,
        status,
        fetchReminders,
        createReminder,
        deleteReminder,
    }
}