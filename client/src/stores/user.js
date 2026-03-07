import { defineStore } from "pinia";
import router from "../router";

export const useUserStore = defineStore("user", {
    state: () => ({
        user: null,
        jwt: localStorage.getItem("jwt") || null,
    }),
    actions: {
        setUser(user) {
            this.user = user;
        },
        setToken(token) {
            this.jwt = token;
            localStorage.setItem("jwt", token);
        },
        logout() {
            this.user = null;
            this.jwt = null;
            localStorage.removeItem("jwt");
            router.push("/login"); // logout sonrası login sayfasına yönlendir
        },
    },
});