<template>
  <div style="display: none"></div>
</template>

<script>
import router from "@/router";

export default {
  name: "check",
  methods: {
    check: function() {
      if (!this.$store.state.isLogin) {
        router.push("/home");
      }
    }
  },
  beforeMount() {
    this.check();
  }
};
</script>

<style scoped></style>
