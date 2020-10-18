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
    <p style="text-align: center;margin: 25px;"><a href="https://github.com/johnpoint/RssReader">RssReader</a> -
      {{ $t("footer.text") }}</p>
  </div>
</template>

<script>
// @ is an alias to /src
import axios from "axios";
import config from "@/config";

export default {
  name: "About",
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
  data() {
    return {
      Sysposts: "",
    };
  },
};
</script>
