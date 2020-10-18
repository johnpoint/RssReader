<template>
  <div class="home">
    <img alt="RRRssReader" src="../assets/logo.png" style="width: 256px;height: 256px"/>
    <h1>Welcome to RssReader</h1>
    <p>
      {{ $t("about.text") }}
    </p>
    <b-card v-if="Sysposts!==''">
      <p>公告</p>
      <p v-html="Sysposts"></p>
    </b-card>
  </div>
</template>

<script>
// @ is an alias to /src
import router from "@/router";
import axios from "axios";
import config from "@/config";

export default {
  name: "Home",
  components: {},
  beforeMount() {
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true);
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"));
      if (window.localStorage.getItem("config") !== null) {
        this.$store.state.config = JSON.parse(window.localStorage.getItem("config"))
      } else {
        window.localStorage.setItem("config", JSON.stringify({"postnum": 50}))
      }
    }
    if (this.$store.state.isLogin) {
      router.push("/posts")
    }
    this.getSysPost()
  },
  methods:{
    getSysPost: function () {
      axios.get(config.apiAddress + "/api/syspost").then((response) => {
        if (response.data.code === 200) {
          this.Sysposts = response.data.message
        }
      })
    }
  },
  data() {
    return {
      Sysposts: "",
    };
  },
};
</script>
