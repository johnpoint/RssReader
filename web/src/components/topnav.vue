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
  </div>
</template
>

<script>

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
  }
};
</script>

<style scoped></style>
