import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/Login.vue";
import Register from "../views/Register.vue";
import Dashboard from "../views/Dashboard.vue";
import AddTransaction from "../pages/AddTransaction.vue";
import { useUserStore } from "../stores/user";

const routes = [
    { path: "/", redirect: "/login" },

    { path: "/login", component: Login },
    { path: "/register", component: Register },

    {
        path: "/dashboard",
        component: Dashboard,
        meta: { requiresAuth: true }
    },

    {
        path: "/add-transaction",
        component: AddTransaction,
        meta: { requiresAuth: true }
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, from, next) => {
    const userStore = useUserStore();
    const token = userStore.jwt || localStorage.getItem("jwt");

    if (to.meta.requiresAuth && !token) {
        next("/login");
    }
    else if ((to.path === "/login" || to.path === "/register") && token) {
        next("/dashboard");
    }
    else {
        next();
    }
});

export default router;