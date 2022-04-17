<template>
  <div id="app">
    <topnav/>
    <div id="main">
      <router-view/>
    </div>
  </div>
</template>

<script>
import topnav from "@/components/topnav";
import axios from "axios";
import config from "@/config";

export default {
  components: {
    topnav
  },
  mounted() {
    axios.get("/api_address").then(
        (response) => {
          this.$store.state.apiAddress = response.data
          config.apiAddress = response.data
          console.log("api_address: ", response.data)
        },
    )
  }
};
</script>

<style lang="less">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#main {
  width: 90%;
  margin: auto;
}

#topnav {
  padding: 30px;

a {
  font-weight: bold;
  color: #2c3e50;

&
.router-link-exact-active {
  color: #42b983;
}

}
}
</style>
