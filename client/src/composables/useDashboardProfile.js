import { reactive, ref } from 'vue'
import { authService } from '../services/portfolio'
import { createStatusState, clearStatus, setError, setOk } from './dashboardShared'

export function useDashboardProfile(userStore) {
    const loading = ref(false)
    const status = createStatusState()
    const form = reactive({ name: '', surname: '', email: '', password: '' })

    async function fetchProfile() {
        const meRes = await authService.me()
        userStore.setProfile(meRes.data)
        form.email = meRes.data?.email || ''
        return meRes
    }

    async function updateProfile() {
        try {
            loading.value = true
            clearStatus(status)
            const payload = Object.fromEntries(
                Object.entries(form).filter(([, value]) => String(value || '').trim() !== ''),
            )

            if (!Object.keys(payload).length) {
                setOk(status, 'Güncellenecek alan bulunamadı.')
                return
            }

            await authService.updateMe(payload)
            setOk(status, 'Profil güncellendi.')
            Object.assign(form, {
                name: '',
                surname: '',
                email: payload.email || form.email,
                password: '',
            })

            await fetchProfile()
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    async function deleteProfile() {
        if (!window.confirm('Hesabınızı silmek istediğinize emin misiniz?')) return

        try {
            loading.value = true
            clearStatus(status)
            await authService.deleteMe()
            userStore.logout()
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    return {
        form,
        loading,
        status,
        fetchProfile,
        updateProfile,
        deleteProfile,
    }
}