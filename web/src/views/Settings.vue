<template>
  <div class="home">
    <div class="setting">
      <label>
        {{ $t("setting.CachedArticle") }}: {{ localpostnum }}
      </label>
      <b-button size="sm" style="float: right" variant="outline-primary" @click="clearCache">{{ $t("setting.clearCache") }}</b-button>
    </div>
    <hr>
    <div class="setting">
      <span>语言 / Language</span>
      <b-button size="sm" style="float: right" @click="changeLg()" variant="outline-primary">{{
          $i18n.locale == "zh" ? "EN" : "中文"
        }}
      </b-button>
    </div>
    <hr>
    <div class="setting">
        <b-button size="sm" style="float: right" variant="outline-danger" @click="logout">{{ $t("setting.exit") }}</b-button>
    </div>
    <div v-if="showLoading">
      <b-spinner class="loading" variant="success" label="Spinning"></b-spinner>
    </div>
  </div>
</template>

<script>

import router from "@/router";

export default {
  name: "Settings",
  components: {},
  data() {
    return {
      localpostnum: 0,
      localpost: [],
      showLoading: false
    };
  },
  beforeMount() {
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true);
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"));
    }
    if (!this.$store.state.isLogin) {
      router.push("/login");
    }
    this.getCache()
  },
  methods: {
    changeLg: function () {
      this.$i18n.locale = (this.$i18n.locale === "zh" ? "en" : "zh")
      window.localStorage.setItem("i18n", this.$i18n.locale)
    },
    logout: function () {
      this.$store.commit("setStatus", false);
      window.localStorage.removeItem("jwt")
      window.localStorage.setItem("login", false)
      window.localStorage.removeItem("posts")
      window.localStorage.removeItem("feeds")
      router.push("/");
    },
    getCache: function () {
      this.localpost = []
      this.localpostnum = 0
      Object.keys(window.localStorage).forEach(i => {
        if (i.indexOf("post") != -1 && i.indexOf("posts") == -1) {
          this.localpost.push(i)
          this.localpostnum++
        }
      })
    },
    clearCache: function () {
      for (let i of this.localpost) {
        window.localStorage.removeItem(i)
      }
      this.getCache()
    }
  }
};
</script>
