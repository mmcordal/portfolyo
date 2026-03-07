<template>
  <nav class="navbar">
    <div class="nav-left">
      <h1 class="logo" @click="goDashboard">Portfolyo</h1>
    </div>

    <div class="nav-right" v-if="userStore.jwt">
      <router-link class="nav-link" to="/dashboard">Dashboard</router-link>

      <router-link class="nav-link" to="/add-transaction">
        Yeni İşlem
      </router-link>

      <span class="user-name" v-if="userStore.user">
        {{ userStore.user.name }} {{ userStore.user.surname }}
      </span>

      <button @click="logout" class="logout-btn">
        Çıkış Yap
      </button>
    </div>
  </nav>
</template>

<script setup>
import { useUserStore } from "../stores/user";
import { useRouter } from "vue-router";

const userStore = useUserStore();
const router = useRouter();

const logout = () => {
  userStore.logout();
  router.push("/login");
};

const goDashboard = () => {
  router.push("/dashboard");
};
</script>

<style scoped>
.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 32px;
  background: white;
  border-bottom: 1px solid #eee;
}

.logo {
  font-size: 22px;
  font-weight: bold;
  cursor: pointer;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-name {
  font-weight: 500;
}

.nav-link {
  text-decoration: none;
  font-weight: 500;
  color: #111;
}

.nav-link:hover {
  color: #4f46e5;
}

.logout-btn {
  background: #ef4444;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
}
</style>