import { reactive, computed } from 'vue'
import router from '../router'

const state = reactive({
    profile: null,
    jwt: localStorage.getItem('jwt') || null,
})

export function useUserStore() {
    const isAuthenticated = computed(() => Boolean(state.jwt))
    const fullName = computed(() => state.profile?.full_name || '-')

    function setProfile(profile) {
        state.profile = profile
    }

    function setToken(token) {
        state.jwt = token
        localStorage.setItem('jwt', token)
    }

    function clearSession() {
        state.profile = null
        state.jwt = null
        localStorage.removeItem('jwt')
    }

    function logout() {
        clearSession()
        router.push('/login')
    }

    return {
        get profile() {
            return state.profile
        },
        get jwt() {
            return state.jwt
        },
        isAuthenticated,
        fullName,
        setProfile,
        setToken,
        clearSession,
        logout,
    }
}