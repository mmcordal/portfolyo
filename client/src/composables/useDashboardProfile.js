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
                Object.entries({
                    name: form.name,
                    surname: form.surname,
                    email: form.email,
                }).filter(([, value]) => String(value || '').trim() !== ''),
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
            })

            await fetchProfile()
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    async function updatePassword() {
        try {
            loading.value = true
            clearStatus(status)
            const password = String(form.password || '').trim()

            if (!password) {
                setOk(status, 'Lütfen yeni şifreyi girin.')
                return
            }

            await authService.updateMe({ password })
            setOk(status, 'Şifreniz güncellendi.')
            form.password = ''
        } catch (err) {
            setError(status, err)
        } finally {
            loading.value = false
        }
    }

    async function deleteProfile() {
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
        updatePassword,
        deleteProfile,
    }
}