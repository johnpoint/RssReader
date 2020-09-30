<template>
  <div class="home">
    <hr>
    <div class="setting">
      <label>
        {{ $t("setting.CachedArticle") }}: {{ localpostnum }}
      </label>
      <b-button size="sm" style="float: right;margin: 5px" variant="outline-primary" @click="clearCache">
        {{ $t("setting.clearCache") }}
      </b-button>
      <b-button size="sm" style="float: right;margin: 5px" variant="outline-info"
                @click="showAnalysis();Analysis=!Analysis">
        {{ Analysis ? $t("setting.hide") : $t("setting.info") }}
      </b-button>
      <div v-if="Analysis">
        <div v-for="i in cacheAnalysis" :key="i">
          {{ i.source }}: {{ i.num }}
        </div>
      </div>
    </div>
    <hr>
    <div class="setting">
      <span>语言 / Language</span>
      <b-button size="sm" style="float: right;margin: 5px" @click="changeLg()" variant="outline-primary">{{
          $i18n.locale == "zh" ? "EN" : "中文"
        }}
      </b-button>
    </div>
    <hr>
    <div class="setting">
      <b-button size="sm" style="float: right;margin: 5px" variant="outline-danger" @click="logout">
        {{ $t("setting.exit") }}
      </b-button>
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
      showLoading: false,
      cache: [],
      cacheAnalysis: [],
      Analysis: false
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
    showAnalysis: function () {
      this.cache = []
      let keys = Object.keys(this.cache)
      for (let i of this.localpost) {
        let csource = JSON.parse(window.localStorage.getItem(i)).source
        if (keys.indexOf(csource) == -1) {
          this.cache[csource] = 1
        } else {
          this.cache[csource] += 1
        }
      }
      this.cacheAnalysis = []
      keys = Object.keys(this.cache)
      for (let i of keys) {
        this.cacheAnalysis.push(
            {
              "source": i,
              "num": this.cache[i]
            }
        )
      }
    },
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
        if (i.indexOf("post") !== -1 && i.indexOf("posts") === -1) {
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
