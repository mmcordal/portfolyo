<template>
  <div>
    <h2>Kayıt Ol</h2>
    <form @submit.prevent="register">
      <input v-model="name" placeholder="İsim" required />
      <input v-model="surname" placeholder="Soyisim" required />
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Şifre" required />
      <button type="submit">Kayıt Ol</button>
    </form>
    <p v-if="error" style="color:red">{{ error }}</p>
    <router-link to="/login">Login</router-link>
  </div>
</template>

<script setup>
import { ref } from "vue";
import api from "../api";
import { useRouter } from "vue-router";

const name = ref("");
const surname = ref("");
const email = ref("");
const password = ref("");
const error = ref(null);
const router = useRouter();

const register = async () => {
  try {
    await api.post("/auth/register", { name: name.value, surname: surname.value, email: email.value, password: password.value });
    error.value = null;
    router.push("/login");
  } catch (err) {
    error.value = err.response?.data?.message || err.message;
  }
};
</script>