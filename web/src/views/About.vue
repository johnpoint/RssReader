<template>
  <div class="home">
    <img alt="RRRssReader" src="../assets/rss.png" />
    <HelloWorld msg="Welcome to RRRssReader" />
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from "@/components/HelloWorld.vue";

export default {
  name: "About",
  components: {
    HelloWorld
  }
};
</script>
