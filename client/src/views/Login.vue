<template>
  <div>
    <h2>Login</h2>
    <form @submit.prevent="login">
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="password" type="password" placeholder="Şifre" required />
      <button type="submit">Login</button>
    </form>
    <p v-if="error" style="color:red">{{ error }}</p>
    <router-link to="/register">Kayıt Ol</router-link>
  </div>
</template>

<script setup>
import { ref } from "vue";
import api from "../api";
import { useUserStore } from "../stores/user";
import { useRouter } from "vue-router";

const email = ref("");
const password = ref("");
const error = ref(null);

const userStore = useUserStore();
const router = useRouter();

const login = async () => {
  try {
    const res = await api.post("/auth/login", {
      email: email.value,
      password: password.value
    });

    // TOKEN
    const token = res.data.data.token;
    userStore.setToken(token);

    // USER INFO
    const userRes = await api.get("/users/me");
    userStore.setUser(userRes.data);

    error.value = null;

    router.push("/dashboard");

  } catch (err) {
    error.value = err.response?.data?.message || err.message;
  }
};
</script>