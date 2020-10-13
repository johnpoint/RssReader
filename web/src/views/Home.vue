<template>
  <div class="home">
    <img alt="RRRssReader" src="../assets/logo.png" style="width: 256px;height: 256px"/>
    <h1>Welcome to RssReader</h1>
    <p>
      {{ $t("about.text") }}
    </p>
    <HelloWorld msg="Welcome to RssReader"/>
  </div>
</template>

<script>
// @ is an alias to /src
import router from "@/router";

export default {
  name: "Home",
  components: {},
  beforeMount() {
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true)
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"))
    }
    if (this.$store.state.isLogin) {
      router.push("/posts")
    }
  },
};
</script>
