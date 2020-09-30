<template>
  <div>
    <h3>404 | Not Found</h3>
    <h5>3 秒后跳转回主页</h5>
  </div>
</template>

<script>
import router from "@/router";

export default {
  name: "notFound",
  components: {},
  mounted() {
    this.gotime()
  },
  methods: {
    gotime: function () {
      setTimeout(function() { router.push("/"); }, 3000);
    }
  }
};
</script>
