<template>
  <section class="auth-container">
    <div class="auth-card">
      <h2>Kayıt Ol</h2>
      <form @submit.prevent="submit" class="inline-form auth-form">
        <input v-model="form.name" type="text" placeholder="Ad" required minlength="2" />
        <input v-model="form.surname" type="text" placeholder="Soyad" required minlength="2" />
        <input v-model="form.email" type="email" placeholder="E-posta" required />
        <input v-model="form.password" type="password" placeholder="Şifre" required minlength="8" />
        <button type="submit" :disabled="loading">{{ loading ? 'Kaydediliyor...' : 'Kayıt Ol' }}</button>
      </form>
      <p class="ok auth-feedback" v-if="success">{{ success }}</p>
      <p class="error auth-feedback" v-if="error">{{ error }}</p>
      <router-link to="/login" class="auth-link">Hesabın var mı? Giriş yap.</router-link>
    </div>
  </section>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { extractErrorMessage } from '../api'
import { authService } from '../services/portfolio'

const form = reactive({ name: '', surname: '', email: '', password: '' })
const loading = ref(false)
const success = ref('')
const error = ref('')
const router = useRouter()

async function submit() {
  try {
    loading.value = true
    success.value = ''
    error.value = ''

    await authService.register(form)
    success.value = 'Kayıt başarılı. Giriş sayfasına yönlendiriliyorsunuz...'
    setTimeout(() => router.push('/login'), 900)
  } catch (err) {
    error.value = extractErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>