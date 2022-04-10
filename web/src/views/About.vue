<template>
  <div class="home">
    <img alt="RRRssReader" src="../assets/logo.png" style="width: 256px;height: 256px"/>
    <h1>Welcome to RssReader</h1>
    <p>
      {{ $t("about.text") }}
    </p>
    <p style="text-align: center;margin: 25px;"><a href="https://github.com/johnpoint/RssReader">RssReader</a> -
      {{ $t("footer.text") }}</p>
  </div>
</template>

<script>

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
  },
  data() {
    return {
    };
  },
};
</script>
