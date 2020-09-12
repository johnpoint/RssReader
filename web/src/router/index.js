import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";

Vue.use(VueRouter);

const routes = [
    {
        path: "/",
        name: "Home",
        component: Home
    },
    {
        path: "/overview",
        name: "Overview",
        component: () => import("../views/Overview.vue")
    },
    {
        path: "/manager",
        name: "Manager",
        component: () => import("../views/Manager.vue")
    },
    {
        path: "/about",
        name: "About",
        component: () => import("../views/About.vue")
    },
    {
        path: "/login",
        name: "Login",
        component: () => import("../views/Login.vue")
    },
    {
        path: "/register",
        name: "Register",
        component: () => import("../views/Register.vue")
    },
    {
        path: "*",
        name: "404",
        component: () => import("../views/404.vue")
    }
];

const router = new VueRouter({
    mode: "history",
    base: process.env.BASE_URL,
    routes
});

export default router;
