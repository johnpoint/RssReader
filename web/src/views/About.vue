<template>
  <div class="home">
    <img alt="RssReader" src="../assets/rss.png" />
    <HelloWorld msg="Welcome to RssReader" />
    <a href="https://github.com/johnpoint/RssReader">Github</a>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from "@/components/Panel.vue";

export default {
  name: "About",
  components: {
    HelloWorld
  },
  beforeMount() {
    if (window.localStorage.getItem("login")==="true") {
      this.$store.commit("setStatus", true)
      this.$store.commit("setjwt", window.localStorage.getItem("jwt"))
    }
  }
};
</script>
