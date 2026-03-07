<template>
  <section class="auth-container">
    <div class="card auth-card">
      <h2>Giriş Yap</h2>
      <form @submit.prevent="submit">
        <input v-model="form.email" type="email" placeholder="E-posta" required />
        <input v-model="form.password" type="password" placeholder="Şifre" required />
        <button type="submit" :disabled="loading">{{ loading ? 'Giriş...' : 'Giriş Yap' }}</button>
      </form>
      <p class="error" v-if="error">{{ error }}</p>
      <router-link to="/register">Hesabın yok mu? Kayıt ol.</router-link>
    </div>
  </section>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { extractErrorMessage } from '../api'
import { authService } from '../services/portfolio'
import { useUserStore } from '../stores/user'

const form = reactive({ email: '', password: '' })
const loading = ref(false)
const error = ref('')
const router = useRouter()
const userStore = useUserStore()

async function submit() {
  try {
    loading.value = true
    error.value = ''

    const loginRes = await authService.login(form)
    userStore.setToken(loginRes.data.data.token)

    const meRes = await authService.me()
    userStore.setProfile(meRes.data.data)

    router.push('/dashboard')
  } catch (err) {
    error.value = extractErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>