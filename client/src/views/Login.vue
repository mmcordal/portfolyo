<template>
  <section class="auth-container">
    <div class="auth-card">
      <p class="auth-eyebrow">Portfolyo</p>
      <h1>Giriş Yap</h1>
      <p class="auth-subtitle">Hesabınıza erişmek için bilgilerinizi girin.</p>
      <form @submit.prevent="submit" class="inline-form auth-form">
        <input v-model="form.email" type="email" placeholder="E-posta adresiniz" required />
        <input v-model="form.password" type="password" placeholder="Şifreniz" required />
        <button type="submit" :disabled="loading">{{ loading ? 'Giriş...' : 'Giriş Yap' }}</button>
      </form>
      <p class="error auth-feedback" v-if="error">{{ error }}</p>
      <p class="auth-switch">
        Hesabın yok mu?
        <router-link to="/register" class="auth-link">Kayıt ol</router-link>
      </p>
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

    const loginData = await authService.login(form)
    userStore.setToken(loginData.data.token)

    const meData = await authService.me()
    userStore.setProfile(meData.data)

    router.push('/dashboard')
  } catch (err) {
    error.value = extractErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>