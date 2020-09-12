<template>
  <div class="home">
    <img alt="RRRssReader" src="../assets/rss.png"/>
    <HelloWorld msg="Welcome to RssReader"/>
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from "@/components/Panel.vue";
import router from "@/router";

export default {
  name: "Home",
  components: {
    HelloWorld
  },
  beforeMount() {
    if (this.$store.state.isLogin) {
      router.push("/overview")
    }
  },
};
</script>
