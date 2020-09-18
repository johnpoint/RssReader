<template>
  <div id="topnav">
    <router-link style="margin: 5px" v-if="this.$store.state.isLogin" to="/posts">{{ $t("nav.posts") }}
    </router-link
    >
    <span style="margin: 5px" v-if="this.$store.state.isLogin"> | </span>
    <router-link style="margin: 5px" v-if="!this.$store.state.isLogin" to="/">{{ $t("nav.home") }}
    </router-link
    >
    <span style="margin: 5px" v-if="!this.$store.state.isLogin"> | </span>
    <router-link style="margin: 5px" v-if="!this.$store.state.isLogin" to="/login"
    >{{ $t("nav.login") }}
    </router-link
    >
    <span style="margin: 5px" v-if="!this.$store.state.isLogin"> | </span>
    <router-link style="margin: 5px" v-if="!this.$store.state.isLogin" to="/register"
    >{{ $t("nav.register") }}
    </router-link
    >
    <span style="margin: 5px" v-if="!this.$store.state.isLogin"> | </span>
    <router-link style="margin: 5px" v-if="this.$store.state.isLogin" to="/feeds"
    >{{ $t("nav.feeds") }}
    </router-link
    >
    <span style="margin: 5px" v-if="this.$store.state.isLogin"> | </span>
    <router-link style="margin: 5px" v-if="this.$store.state.isLogin" to="/settings"
    >{{ $t("nav.settings") }}
    </router-link
    >
    <span style="margin: 5px" v-if="this.$store.state.isLogin"> | </span>
    <router-link style="margin: 5px" to="/about">{{ $t("nav.about") }}</router-link>


    <span style="margin: 5px"> | </span>
    <a @click="changeLg()" style="margin: 5px"
    >{{ $i18n.locale == "zh" ? "EN" : "中文" }}</a
    >
    <span v-if="this.$store.state.isLogin" style="margin: 5px"> | </span>
    <a v-if="this.$store.state.isLogin" @click="logout" style="margin: 5px"
    >{{ $t("nav.exit") }}</a
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
  beforeMount() {
    if (window.localStorage.getItem("i18n") === undefined) {
      window.localStorage.setItem("i18n", "zh");
    } else {
      this.$i18n.locale = window.localStorage.getItem("i18n")
    }
  },
  methods: {
    logout: function () {
      this.$store.commit("setStatus", false);
      window.localStorage.removeItem("jwt")
      window.localStorage.setItem("login", false)
      window.localStorage.removeItem("posts")
      window.localStorage.removeItem("feeds")
      router.push("/");
    },
    changeLg: function () {
      this.$i18n.locale = (this.$i18n.locale === "zh" ? "en" : "zh")
      window.localStorage.setItem("i18n", this.$i18n.locale)
    }
  }
};
</script>

<style scoped></style>
