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
// @ is an alias to /src

export default {
  name: "About",
  components: {
  },
  beforeMount() {
    if (window.localStorage.getItem("login") === "true") {
      this.$store.commit("setStatus", true)
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"))
    }
  }
};
</script>
