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
    <span style="margin: 5px" v-if="Sysposts!==''"> | </span>
    <router-link style="margin: 5px" to="/news">{{ $t("nav.news") }}</router-link>
  </div>
</template
>

<script>

import axios from "axios";
import config from "@/config";

export default {
  name: "topnav",
  data() {
    return {
      Sysposts: "",
    };
  },
  beforeMount() {
    if (window.localStorage.getItem("i18n") === undefined) {
      window.localStorage.setItem("i18n", "zh");
    } else {
      this.$i18n.locale = window.localStorage.getItem("i18n")
    }
    this.getSysPost()
  },
  methods: {
    getSysPost: function () {
      axios.get(config.apiAddress + "/api/syspost").then((response) => {
        if (response.data.code === 200) {
          this.Sysposts = response.data.message
        }
      })
    }
  },
};
</script>

<style scoped></style>
