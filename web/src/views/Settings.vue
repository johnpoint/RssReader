<template>
  <div class="home">
    <b-card
        title="本地缓存"
        style="text-align: left"
    >
      <b-card-text>
        文章: {{ localpostnum }}
      </b-card-text>
      <a class="card-link" @click="clearCache">清除缓存</a>
    </b-card>
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
      localpost: []
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
    getCache: function () {
      this.localpost = []
      Object.keys(window.localStorage).forEach(i => {
        if (i.indexOf("post") != -1 && i.indexOf("posts") == -1) {
          this.localpost.push(i)
          this.localpostnum++
        }
      })
    },
    clearCache: function () {
      this.localpost.forEach(i => {
        window.localStorage.removeItem(i)
      }, this.getCache())
    }
  }
};
</script>
