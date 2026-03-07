import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', component: () => import('../views/Login.vue') },
    { path: '/register', component: () => import('../views/Register.vue') },
    {
        path: '/profile',
        component: () => import('../views/Profile.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true },
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

router.beforeEach((to) => {
    const userStore = useUserStore()
    const token = userStore.jwt || localStorage.getItem('jwt')

    if (to.meta.requiresAuth && !token) {
        return '/login'
    }
    if ((to.path === '/login' || to.path === '/register') && token) {
        return '/dashboard'
    }
    return true
})

export default router