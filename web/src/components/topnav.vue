<template>
  <div id="topnav" align="left">
    <router-link v-if="this.$store.state.isLogin" to="/">Home</router-link
    ><span v-if="this.$store.state.isLogin"> | </span>
    <router-link v-if="!this.$store.state.isLogin" to="/home">Home</router-link
    ><span v-if="!this.$store.state.isLogin"> | </span>
    <router-link v-if="!this.$store.state.isLogin" to="/login"
      >Login</router-link
    ><span v-if="!this.$store.state.isLogin"> | </span>
    <router-link v-if="!this.$store.state.isLogin" to="/register"
      >Register</router-link
    ><span v-if="!this.$store.state.isLogin"> | </span>
    <router-link v-if="this.$store.state.isLogin" to="/manager"
      >Manager</router-link
    ><span v-if="this.$store.state.isLogin"> | </span>
    <router-link to="/about">About</router-link>
    <a v-if="this.$store.state.isLogin" @click="logout" style="float: right"
      >Exit</a
    >
  </div></template
>

<script>
import router from "@/router";

export default {
  name: "topnav",
  data() {
    return {};
  },
  methods: {
    logout: function() {
      this.$store.commit("setStatus", false);
      router.push("/home");
    }
  }
};
</script>

<style scoped></style>
