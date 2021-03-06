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
        path: "/posts",
        name: "Posts",
        component: () => import("../views/Posts.vue")
    },
    {
        path: "/news",
        name: "News",
        component: () => import("../views/News.vue")
    },
    {
        path: "/post/:id",
        name: "Post",
        component: () => import("../views/Post.vue")
    },
    {
        path: "/settings",
        name: "Settings",
        component: () => import("../views/Settings.vue")
    },
    {
        path: "/account",
        name: "Account",
        component: () => import("../views/Account.vue")
    },
    {
        path: "/feeds",
        name: "Feeds",
        component: () => import("../views/Feeds.vue")
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
