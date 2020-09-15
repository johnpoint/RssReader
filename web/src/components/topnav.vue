<template>
  <div id="topnav" align="left">
    <router-link v-if="this.$store.state.isLogin" to="/posts">文章
    </router-link
    >
    <span v-if="this.$store.state.isLogin"> | </span>
    <router-link v-if="!this.$store.state.isLogin" to="/">主页
    </router-link
    >
    <span v-if="!this.$store.state.isLogin"> | </span>
    <router-link v-if="!this.$store.state.isLogin" to="/login"
    >登录
    </router-link
    >
    <span v-if="!this.$store.state.isLogin"> | </span>
    <router-link v-if="!this.$store.state.isLogin" to="/register"
    >注册
    </router-link
    >
    <span v-if="!this.$store.state.isLogin"> | </span>
    <router-link v-if="this.$store.state.isLogin" to="/feeds"
    >订阅源
    </router-link
    >
    <span v-if="this.$store.state.isLogin"> | </span>
    <router-link v-if="this.$store.state.isLogin" to="/settings"
    >设置
    </router-link
    >
    <span v-if="this.$store.state.isLogin"> | </span>
    <router-link to="/about">关于</router-link>
    <a v-if="this.$store.state.isLogin" @click="logout" style="float: right"
    >退出</a
    >
  </div>
</template
>

<script>
import router from "@/router";

export default {
  name: "topnav",
  data() {
    return {};
  },
  methods: {
    logout: function () {
      this.$store.commit("setStatus", false);
      window.localStorage.removeItem("jwt")
      window.localStorage.setItem("login", false)
      window.localStorage.removeItem("posts")
      window.localStorage.removeItem("feeds")
      router.push("/");
    }
  }
};
</script>

<style scoped></style>
