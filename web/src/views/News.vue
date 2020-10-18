<template>
  <div class="home">
    <b-card v-if="Sysposts!==''">
      <p>公告</p>
      <p v-html="Sysposts"></p>
    </b-card>
  </div>
</template>

<script>

import axios from "axios";
import config from "@/config";

export default {
  name: "New",
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
